package main

import (
	"flag"
	"fmt"
	logger "github.com/shengkehua/xlog4go"
	"io"
	"model"
	"net"
	"process"
	"rpc"
	"time"
)

var (
	logFile = flag.String("l", "./conf/log.json", "log config file path")
)

/*
	处理与客户端通信，参数是连接信息
*/
func process(conn net.Conn) {
	defer conn.Close()

	//调用总控
	processor := processes.Processor{
		Conn: conn,
	}

	err := processor.Process2()

	if err != nil {
		if err == io.EOF {
			//fmt.Println("客户端下线")
			return
		}
		logger.Error("client communicate with service err, err=%s", err.Error())
		return
	}
}

//定义一个函数，完成对UserDao的初始化
func initUserDao() {

	/*
		pool是一个全局变量，在redis.go中定义的
		注意，必须先初始化Pool，在执行该函数，否则没有一个pool对象
	*/
	model.MyUserDao = model.NewUserDao(rpc.Pool)
}

func main() {

	// init log
	if err := initLog(*logFile); err != nil {
		panic(err)
	}
	logger.Info("init log done!")

	//服务器启动时，初始化一个redis连接池
	rpc.InitPool("127.0.0.1:6379", 16, 0, 300*time.Second)

	initUserDao()

	//服务器开始监听
	listen, err := net.Listen("tcp", "0.0.0.0:8888")

	defer listen.Close()

	if err != nil {
		logger.Error("server listen err, err=%s", err.Error())
		return
	}
	logger.Info("service start success, port=8888")

	//监听成功，等待客户端来连接服务器
	for {
		logger.Debug("waiting for client connect...")
		conn, err := listen.Accept()
		if err != nil {
			logger.Error("client connect err, err=%s", err.Error())
			//此处不用return，万一只是一条连接出错，不能将整个服务器退出
			// return
		}

		//连接成功后，启动一个协程与客户端保持通信
		go process(conn)
	}

}

func initLog(path string) (err error) {
	err = logger.SetupLogWithConf(path)
	if err != nil {
		fmt.Printf("log init fail, err=%s\n", err.Error())
		return
	}
	return nil
}
