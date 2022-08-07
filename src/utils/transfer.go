package utils

import (
	"common"
	"encoding/binary"
	"encoding/json"
	logger "github.com/shengkehua/xlog4go"
	"io"
	"net"
)

// Transfer 定义一个结构体，将传输的函数都绑定在它身上
type Transfer struct {
	Conn net.Conn

	//传输时的缓存，先定义数组，使用时当做byte切片
	buf [8096]byte
}

// ReadPkg 接收函数，用于反序列化接收的数据，并将反序列化的结果返回
func (t *Transfer) ReadPkg() (mes common.Message, err error) {
	//创建byte切片，用于接收用户输入
	//buf := make([]byte, 4096)
	/*
		conn.Read阻塞的前提是conn没有被关闭的前提下，才会阻塞
		如果客户端关闭了conn连接，则不会阻塞
	*/
	n, err := t.Conn.Read(t.buf[:4])

	if n != 4 || err != nil {
		if err == io.EOF {
			logger.Info("client close connect, service also close connect")
			return
		}
		// err = errors.New("读取数据长度失败")
		logger.Error("read data length err, err=%s", err.Error())
		return
	}
	//fmt.Println("读到的buf=",t.buf[:4])
	dataLen := binary.BigEndian.Uint32(t.buf[:4])
	n, err = t.Conn.Read(t.buf[:dataLen])
	if n != int(dataLen) || err != nil {
		// err = errors.New("读取数据失败")
		logger.Error("read data length err, err=%s", err.Error())
		return
	}
	//将读取到的消息反序列化
	err = json.Unmarshal(t.buf[:dataLen], &mes)
	if err != nil {
		logger.Error("mes unmarshal err, err=%s", err.Error())
		return
	}
	return
}

// WritePkg 发送函数，用于向目标发送已经序列化好的数据
func (t *Transfer) WritePkg(data []byte) (err error) {

	//为了确保tcp发送消息的准确性，先发送mes的长度，再发送mes消息本体
	//先获取data这个byte切片的长度，然后将长度数据转化为一个byte切片
	var dataLen uint32
	dataLen = uint32(len(data))
	//var dataLenbytes [4]byte
	//此函数可以将传入的一个uint32的数值转化为一个byte切片
	binary.BigEndian.PutUint32(t.buf[:4], dataLen)
	//发送长度数据
	n, err := t.Conn.Write(t.buf[:4])
	if n != 4 || err != nil {
		logger.Error("send data length err, err=%s", err.Error())
		return
	}
	//发送数据
	_, err = t.Conn.Write(data)
	if err != nil {
		logger.Error("send data err, err=%s", err.Error())
		return
	}
	return
}
