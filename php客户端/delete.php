<?php

require("rpc.php");

$rpc = new socketRpc();

$data = array(
    'data' => array(
          'table' => 'test.tbl_crm_', // 表前缀
          'filter' => 'id>3', // 查询条件
          'recycle' => true,  // 是否回收数据到另外一个表
          'default' => array(
                'columns' =>'id,name,sex,number1,number2,remark', // 查询字段内容
                'customer_id' => "9"
          )
    )

);

//  rpc 调用
$res = $rpc->call("CCgo.DeleteCrm",json_encode($data),"C9A10DeleteCrm",true);

var_dump($res);
