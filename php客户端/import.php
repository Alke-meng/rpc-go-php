<?php

require("rpc.php");

$rpc = new socketRpc();

$data = array(
    'file' => "/var/www/test.xlsx", // excel 文件地址
     //主要数据
    'data' => array(
          'table' => 'test.tbl_crm_', // 表名前缀
          'field' => 'name,sex,number1,number2,remark', // 表字段
          'default' => array(
                'customer_id' => "9",
          )
    ),
    //插入时的检查 和 插入后的其他操作
    'filter' => array(
                    'numberCheck' => false,  // 号码格式检查
                    'repeatCheckFile' => false,  // 号码文件重复检查
                    'addCalleeList' => true // 后续的添加号码到另外一个表的逻辑
             )
);

// rpc 调用
$res = $rpc->call("CCgo.ImportCrm",json_encode($data),"C9A10Import",false);

var_dump($res);