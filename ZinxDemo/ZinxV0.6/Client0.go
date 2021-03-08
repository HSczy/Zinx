package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx/znet"
)

// 模拟客户端
func main() {
	fmt.Println("Client start...")
	time.Sleep(1 * time.Second)
	// 1. 直接链接远程服务器，得到一个conn链接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("Client 0 start err,exit!")
		return
	}
	for {
		// 2. 链接调用Write 写数据
		// 发送封包的Msg的消息
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMessage(0, []byte("ZinxV0.6 client 0 Test Message")))
		if err != nil {
			fmt.Println("Pack error,", err)
			break
		}
		if _, err := conn.Write(binaryMsg); err != nil {
			fmt.Println("Write error,", err)
			break
		}
		// 服务器就应该给我们回复一个message数据
		// 先读出message的headlen 和 id
		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("Read head err ", err)
			break
		}
		// 将二进制的head拆包到接口体中
		msgHead, err := dp.Unpack(binaryHead)
		if err != nil {
			fmt.Println("client unpack msgHeader error", err)
			break
		}
		// 在读出 message的
		if msgHead.GetMsgLen() > 0 {
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read body err ", err)
				break
			}
			fmt.Println("---->Recv Server Msg ID:", msg.Id, "\t Msg data:", string(msg.Data))
		}

		// 阻塞
		time.Sleep(1 * time.Second)
	}
}
