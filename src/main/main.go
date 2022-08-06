package main

import (
	"fmt"
	"io"
	"model"
	"net"
	"time"
)

/*
	处理与客户端通信，参数是连接信息
*/
func process(conn net.Conn){
	defer conn.Close()

	//调用总控
	processor := Processor{
		Conn :conn,
	}

	err := processor.process2()

	if err != nil{
		if err == io.EOF{
			//fmt.Println("客户端下线")
			return
		}
		fmt.Println("客户端和服务器通信协程出错,err=", err)
		return
	}
}

//定义一个函数，完成对UserDao的初始化
func initUserDao(){

	/*
		pool是一个全局变量，在redis.go中定义的
		注意，必须先初始化Pool，在执行该函数，否则没有一个pool对象
	*/
	model.MyUserDao = model.NewUserDao(pool)
}


func main(){

	//服务器启动时，初始化一个redis连接池
	initPool("127.0.0.1:6379", 16, 0, 300*time.Second)

	initUserDao()
	
	//服务器开始监听
	listen, err := net.Listen("tcp","0.0.0.0:8888")

	defer listen.Close()

	if err != nil{
		fmt.Println("服务器监听出错,err=", err)
		return
	}
	fmt.Println("服务器在8888端口监听...")

	//监听成功，等待客户端来连接服务器
	for{
		fmt.Println("等待客户端来连接")
		conn, err := listen.Accept()
		if err != nil{
			fmt.Println("客户端连接出错,err=", err)
			//此处不用return，万一只是一条连接出错，不能将整个服务器退出
			// return
		}

		//连接成功后，启动一个协程与客户端保持通信
		go process(conn)
	}

}

