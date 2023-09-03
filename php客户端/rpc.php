<?php

class socketRpc {
    private $conn;
    const TCP_PORT = 8081;
    public function __construct() {
        $this->conn = socket_create(AF_INET, SOCK_STREAM, SOL_TCP);

        if (!$this->conn) {
            die("socket create error");
        }
        $check = $this->connect_socket();
        if ($check["code"] !== 0) {
            die("socket connect fail");
        }
    }

    public function connect_socket() {
        socket_set_option($this->conn, SOL_SOCKET, SO_SNDTIMEO, array("sec" => 0, "usec" => 50000));
        $result = @socket_connect($this->conn, '127.0.0.1', self::TCP_PORT);
        if ($result === false) {
            return array("code" => -1, "msg" => "socket connect fail");
        }
        return array("code" => 0, "msg" => "socket connect success");
    }

    public function write_socket($data) {
        $result = @socket_write($this->conn, json_encode($data));
        if (empty($result)) {
             return array("code" => -2, "msg" => "socket write fail");
        }

        return array("code" => 0, "msg" => "socket write success","data" => $data);
    }

    public function read_socket() {
        $rpc = json_decode(@socket_read($this->conn, 1024), true);
         if (empty($rpc)) {
            return array("code" => -3, "msg" => "socket read fail");
         }

        return array("code" => 0, "msg" => "rpc call success","data" => $rpc['result']);
    }

    /**
     * $w = ['id' => time(), 'params' => ["1"], 'method' => 'CCgo.Hi'];
     * $sync 为 true,同步请求等rpc返回
     * $sync 为 false,直接结束
     */
    public function call($method,$params,$id,$sync = true) {
        $rpcData = json_decode($params, true);
        if (!is_string($params) || empty($id) || is_null($rpcData)) {
             return array("code" => -4, "msg" => "rpc call fail");
        }
        $rpcData['traceID'] = $id.'_'.self::GetRandStr();
        $data = array('params' => array(json_encode($rpcData)), 'method' => $method);
        $check = $this->write_socket($data);
        if ($check["code"] !== 0) {
            return $check;
        }

        if ($sync) {
             return $this->read_socket();
        }

        return array("code" => 0, "msg" => "rpc call success");
    }

    public function __destruct() {
        socket_close($this->conn);
    }

    private static function GetRandStr() {
        //字符组合
        $str = 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789';
        $len = strlen($str) - 1;
        $randstr = '';
        for ($i = 0; $i < 6; $i++) {
            $num = mt_rand(0, $len);
            $randstr .= $str[$num];
        }
        return $randstr;
    }
}