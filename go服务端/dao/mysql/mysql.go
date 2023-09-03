package mysql

import (
	"ccgo/settings"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var db *sqlx.DB

func Init(cfg *settings.MySQLConfig) (err error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DbName,
	)

	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Println("connect db failed, err: ", zap.Error(err))
		return
	}

	db.SetMaxOpenConns(viper.GetInt("mysql.max_open_conns"))
	db.SetMaxIdleConns(viper.GetInt("mysql.max_idle_conns"))
	return
}

func GetResultRowsForArray(rows *sql.Rows) (dataMaps [][]interface{}, err error) {
	// 1. 查询到的数据列名、返回值
	columns, _ := rows.Columns() //列名
	count := len(columns)
	values, valuesPoints := make([]interface{}, count), make([]interface{}, count)

	// 2. 遍历Rows读取每一行
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuesPoints[i] = &values[i]
		}

		// 2.1 数据库中读取出每一行数据
		rows.Scan(valuesPoints...) //将所有内容读取进values

		// 2.2 相当于准备接收数据的结构体Product
		var row []interface{}

		// 2.3 将读取到的数据填充到product
		for _, val := range values { // val是每个列对应的值

			// 判断val的值的类型
			var v interface{}
			b, ok := val.([]byte) //判断是否为[]byte
			if ok {
				v = string(b)
			} else {
				v = val
			}

			// 列名与值对应
			row = append(row, v)
		}

		// 将product归到集合中
		dataMaps = append(dataMaps, row)
	}
	return dataMaps, nil
}

func GetResultRowsForMap(rows *sql.Rows) (dataMaps []map[string]interface{}, err error) {
	// 1. 查询到的数据列名、返回值
	columns, _ := rows.Columns() //列名
	count := len(columns)
	values, valuesPoints := make([]interface{}, count), make([]interface{}, count)

	// 2. 遍历Rows读取每一行
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuesPoints[i] = &values[i]
		}

		// 2.1 数据库中读取出每一行数据
		rows.Scan(valuesPoints...) //将所有内容读取进values

		// 2.2 相当于准备接收数据的结构体Product
		row := make(map[string]interface{})

		// 2.3 将读取到的数据填充到product
		for i, val := range values { // val是每个列对应的值
			key := columns[i] //列名

			// 判断val的值的类型
			var v interface{}
			b, ok := val.([]byte) //判断是否为[]byte
			if ok {
				v = string(b)
			} else {
				v = val
			}

			// 列名与值对应
			row[key] = v
		}

		// 将product归到集合中
		dataMaps = append(dataMaps, row)
	}
	return dataMaps, nil
}

func Close() {
	_ = db.Close()
}
