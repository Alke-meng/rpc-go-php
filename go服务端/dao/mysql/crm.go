package mysql

import (
	"bytes"
	"ccgo/logger"
	"ccgo/tool"
	"fmt"
	"io"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/spf13/viper"
	"github.com/xuri/excelize/v2"
)

const (
	ImportMapLength = 10000
	GoHandleData    = 50000
	ReportSheet1    = "报告"
	ReportSheet2    = "重复"
	ReportSheet3    = "失败"
	DataTitle       = 1
	DataFail        = 2
	DataRepeatFile  = 3
	DataRepeatDb    = 4
	PageSize        = 2000
)

type ImportCrmResult struct {
	Total      int                 `json:"total"`
	Success    int                 `json:"success"`
	Fail       int                 `json:"fail"`
	Repeat     int                 `json:"repeat"`
	ReportFile string              `json:"repeatFile"`
	PhoneData  map[string]struct{} `json:"-"`
}

type outFileResult struct {
	Total   int            `json:"total"`
	GoNum   int            `json:"-"`
	OutFile map[int]string `json:"-"`
}

func ImportCrm(data map[string]any, traceID string, startTime time.Time) (result ImportCrmResult, err error) {
	// 数据准备
	dataTmp := data["data"].(map[string]any)
	file := data["file"].(string)
	defaultData := dataTmp["default"].(map[string]any)
	table := dataTmp["table"].(string) + defaultData["customer_id"].(string)
	field := dataTmp["field"].(string)
	filterTmp := data["filter"].(map[string]any)
	// 数据条件
	rowFilter := map[string]any{
		"number":          0,
		"numberCheck":     filterTmp["numberCheck"],
		"repeatCheckFile": filterTmp["repeatCheckFile"],
	}
	fieldArray := strings.Split(dataTmp["field"].(string), ",")
	for key, val := range fieldArray {
		if val == "number1" {
			rowFilter["number"] = key
			break
		}
	}

	// 逻辑执行
	f, err := excelize.OpenFile(file)
	if err != nil {
		logger.CCgoActionLogger(traceID, fmt.Sprintf("excelize OpenFile fail:%v", err))
		return
	}
	defer func() {
		//关闭工作簿
		if err := f.Close(); err != nil {
			logger.CCgoActionLogger(traceID, fmt.Sprintf("excelize Close fail:%v", err))
		}
	}()

	sheetName := f.GetSheetName(f.GetActiveSheetIndex())
	rows, err := f.Rows(sheetName)
	if err != nil {
		logger.CCgoActionLogger(traceID, fmt.Sprintf("excelize read rows fail:%v", err))
		return
	}

	bytesTmp := &bytes.Buffer{}

	result.PhoneData = make(map[string]struct{}, ImportMapLength)
	//导入报告
	reportFile := excelize.NewFile()
	// 创建报告的sheet页
	reportFile.NewSheet(ReportSheet1)
	reportFile.NewSheet(ReportSheet2)
	reportFile.NewSheet(ReportSheet3)

	//迭代
	//记录循环获取
	rowIndex := 0
	for rows.Next() {
		// 总数
		result.Total++
		row, err := rows.Columns()
		if err != nil {
			// 失败
			result.Fail++
			reportFile.SetSheetRow(ReportSheet3, fmt.Sprintf("A%d", result.Fail+1), row)
			continue
		}

		//数据处理
		flag := HandleImportData(row, &result, reportFile, rowFilter)

		row = append(
			row,
			defaultData["customer_id"].(string),
		)

		if flag > 0 {
			continue
		}

		// 总数
		result.Success++
		line := strings.Join(row, "\x1F")
		bytesTmp.WriteString(line)
		bytesTmp.WriteString("\x1E")
		rowIndex++
		if rowIndex%GoHandleData == 0 {
			bytesTmp.WriteString("\x1A")
		}
	}

	if err = rows.Close(); err != nil {
		logger.CCgoActionLogger(traceID, fmt.Sprintf("rows.Close fail:%v", err))
	}

	logger.CCgoActionLogger(traceID, fmt.Sprintf("import data init,cost: %v", time.Now().Sub((startTime))))

	var wg sync.WaitGroup
	bytesSlice := strings.Split(bytesTmp.String(), "\x1A")
	for key, readerData := range bytesSlice {
		wg.Add(1)
		go func(i int, readerData string) {
			defer func() {
				wg.Done()
			}()

			reader := strings.NewReader(readerData)
			registerData := "data" + strconv.Itoa(i)
			mysql.RegisterReaderHandler(registerData, func() io.Reader {
				return io.Reader(reader)
			})
			registerReader := "Reader::" + registerData
			strSql := "LOAD DATA LOCAL INFILE '" + registerReader + "' INTO TABLE " + table + "\n\t\t\tCHARACTER SET UTF8\n\t\t    FIELDS TERMINATED BY X'1F'\n\t\t    LINES TERMINATED BY X'1E'\n\t\t"
			strSql += "(" + field + ");"

			_, err = db.Exec(strSql)
			if err != nil {
				result.Success = 0
				logger.CCgoActionLogger(traceID, fmt.Sprintf("mysql db Exec infile fail:%v", err))
				return
			}
		}(key, readerData)
	}
	wg.Wait()

	logger.CCgoActionLogger(traceID, fmt.Sprintf("import data end,cost: %v", time.Now().Sub((startTime))))

	result.Total--
	// 写报告
	reportFile.SetSheetRow(ReportSheet1, "A1", &[]interface{}{"总数", "成功", "失败", "重复"})
	reportFile.SetSheetRow(ReportSheet1, "A2", &[]interface{}{result.Total, result.Success, result.Fail, result.Repeat})
	reportFile.DeleteSheet("Sheet1")
	result.ReportFile = viper.GetString("report_file_path") + "/report-" + fmt.Sprintf("%v", time.Now().Unix()) + ".xlsx"
	reportFile.SaveAs(result.ReportFile)
	return
}

func HandleImportData(row []string, result *ImportCrmResult, reportFile *excelize.File, rowFilter map[string]any) (flag int) {
	tmp := make([]interface{}, len(row)+1)

	if result.Total == 1 {
		sheet2Title := make([]interface{}, len(row)+1)
		sheet3Title := make([]interface{}, len(row)+1)
		for item, colCell := range row {
			sheet2Title[item] = colCell
			sheet3Title[item] = colCell
		}
		sheet2Title[len(row)] = "原因(1:文件重复)"
		sheet3Title[len(row)] = "原因(1:号码格式问题)"
		// 表头
		reportFile.SetSheetRow(ReportSheet2, "A1", &sheet2Title)
		reportFile.SetSheetRow(ReportSheet3, "A1", &sheet3Title)
		return DataTitle
	}

	for item, colCell := range row {
		tmp[item] = colCell
		if item == rowFilter["number"] {
			// 号码检查
			if rowFilter["numberCheck"] == true {
				res := tool.CheckMobile(colCell)
				if !res {
					flag = DataFail
					continue
				}
			}

			// 文件重复检查
			if rowFilter["repeatCheckFile"] == true {
				if _, ok := result.PhoneData[colCell]; ok {
					flag = DataRepeatFile
				}
			}

			result.PhoneData[colCell] = struct{}{}
		}
	}

	switch flag {
	case DataRepeatFile:
		// 文件重复
		result.Repeat++
		tmp[len(row)] = 1
		reportFile.SetSheetRow(ReportSheet2, fmt.Sprintf("A%d", result.Repeat+1), &tmp)
	case DataFail:
		// 错误
		result.Fail++
		tmp[len(row)] = 1
		reportFile.SetSheetRow(ReportSheet3, fmt.Sprintf("A%d", result.Fail+1), &tmp)
	}

	return
}

func CallerListADD(data map[string]any, traceID string) error {
	bytesTmp := &bytes.Buffer{}
	for item := range data["number"].(map[string]any) {
		row := []string{data["customer_id"].(string), item}
		line := strings.Join(row, "\x1F")
		bytesTmp.WriteString(line)
		bytesTmp.WriteString("\x1E")
	}

	mysql.RegisterReaderHandler("data", func() io.Reader {
		return io.Reader(bytesTmp)
	})

	table := "test.tbl_callee_pool_" + data["customer_id"].(string)
	field := "customer_id,number"

	strSql := "LOAD DATA LOCAL INFILE 'Reader::data' INTO TABLE " + table + "\n\t\t\tCHARACTER SET UTF8\n\t\t    FIELDS TERMINATED BY X'1F'\n\t\t    LINES TERMINATED BY X'1E'\n\t\t"
	strSql += "(" + field + ");"

	_, err := db.Exec(strSql)
	if err != nil {
		logger.CCgoActionLogger(traceID, fmt.Sprintf("mysql db Exec infile fail:%v", err))
		return err
	}

	return nil
}

func DeleteCrm(data map[string]any, traceID string) (result outFileResult, err error) {
	// 数据准备
	dataTmp := data["data"].(map[string]any)
	defaultData := dataTmp["default"].(map[string]any)
	table := dataTmp["table"].(string) + defaultData["customer_id"].(string)
	filterTmp := dataTmp["filter"]
	columns := defaultData["columns"].(string)
	crmRecycle := dataTmp["recycle"]

	sqlStrTotal := fmt.Sprintf("SELECT count(*) FROM %s WHERE %s", table, filterTmp)
	rowsTotal, err := db.Query(sqlStrTotal)
	tmpTotal, err := GetResultRowsForArray(rowsTotal)
	total, err := strconv.Atoi(tmpTotal[0][0].(string))
	// 协程个数
	goNum := runtime.NumCPU()
	goNumFor := total/(PageSize*goNum) + 1

	// 接收文件
	result.OutFile = make(map[int]string, 10)

	if total <= PageSize {
		outFileName := traceID + "_1.txt"
		if crmRecycle == true {
			sqlStr := fmt.Sprintf("SELECT %s FROM %s WHERE id>=(SELECT id FROM %s WHERE  %s ORDER BY id LIMIT %d,1) AND  %s  ORDER BY id LIMIT %s ",
				columns, table, table, filterTmp, 0, filterTmp, strconv.Itoa(PageSize))
			intoSql := sqlStr + " INTO OUTFILE '" + viper.GetString("tmp_file_path") + "/" + outFileName + "'\n\t\t FIELDS TERMINATED BY X'1F'\n\t\t    LINES TERMINATED BY X'1E'\n\t\t"

			_, err = db.Exec(intoSql)
			if err != nil {
				logger.CCgoActionLogger(traceID, fmt.Sprintf("mysql db Exec into outfile fail:%v", err))
				return
			}
		}

		deleteSql := fmt.Sprintf("Delete FROM %s WHERE id>=(SELECT id FROM %s WHERE  %s ORDER BY id LIMIT %d,1) AND  %s  ORDER BY id LIMIT %s ",
			table, table, filterTmp, 0, filterTmp, strconv.Itoa(PageSize))
		_, err = db.Exec(deleteSql)
		if err != nil {
			logger.CCgoActionLogger(traceID, fmt.Sprintf("mysql db Exec delete fail:%v", err))
			return
		}

		if crmRecycle == true {
			result.OutFile[1] = outFileName
		}
		goNum = 1
	} else {
		var wg sync.WaitGroup
		for i := 1; i <= goNum; i++ {
			wg.Add(1)
			go func(i int) {
				defer func() {
					wg.Done()
				}()
				start := (i - 1) * PageSize * goNumFor
				end := PageSize * goNumFor

				sqlStr := fmt.Sprintf("SELECT MIN(id) as min,MAX(id) as max FROM (SELECT id FROM %s WHERE  %s ORDER BY id LIMIT %d,%d) AS tmp",
					table, filterTmp, start, end)

				rows, err := db.Query(sqlStr)
				if err != nil {
					fmt.Printf("query failed, err:%v\n", err)
					return
				}
				// 非常重要：关闭rows释放持有的数据库链接
				defer rows.Close()
				tmp, err := GetResultRowsForMap(rows)
				if err != nil {
					fmt.Printf("rows failed err:%v\n", err)
					return
				}

				// 区间
				region := "id>=" + tmp[0]["min"].(string) + " and id<=" + tmp[0]["max"].(string)
				sqlStr = fmt.Sprintf("SELECT %s FROM %s WHERE %s AND %s",
					columns, table, region, filterTmp)

				outFileName := traceID + "_" + strconv.Itoa(i) + ".txt"
				if crmRecycle == true {
					intoSql := sqlStr + " INTO OUTFILE '" + viper.GetString("tmp_file_path") + "/" + outFileName + "'\n\t\t FIELDS TERMINATED BY X'1F'\n\t\t    LINES TERMINATED BY X'1E'\n\t\t"

					_, err = db.Exec(intoSql)
					if err != nil {
						logger.CCgoActionLogger(traceID, fmt.Sprintf("mysql db Exec into outfile fail:%v", err))
						return
					}
				}

				// 开始删除
				for j := 1; j <= goNumFor; j++ {
					deleteSql := fmt.Sprintf("Delete FROM %s WHERE %s AND %s ORDER BY id LIMIT %s ",
						table, region, filterTmp, strconv.Itoa(PageSize))

					_, err = db.Exec(deleteSql)
					if err != nil {
						logger.CCgoActionLogger(traceID, fmt.Sprintf("mysql db Exec delete fail:%v", err))
						return
					}
				}
				if crmRecycle == true {
					result.OutFile[i] = outFileName
				}
			}(i)
		}
		wg.Wait()
	}

	result.Total = total
	result.GoNum = goNum

	return
}

func CrmRecycle(data map[string]string, traceID string) error {
	fileName := viper.GetString("tmp_file_path") + "/" + data["data"]
	mysql.RegisterLocalFile(fileName)
	strSql := "LOAD DATA LOCAL INFILE '" + fileName + "' INTO TABLE " + data["table"] + "\n\t\t\tCHARACTER SET UTF8\n\t\t    FIELDS TERMINATED BY X'1F'\n\t\t    LINES TERMINATED BY X'1E'\n\t\t"
	strSql += "(" + data["field"] + ");"

	_, err := db.Exec(strSql)
	if err != nil {
		logger.CCgoActionLogger(traceID, fmt.Sprintf("mysql db Exec infile fail:%v", err))
		return err
	}

	return nil
}
