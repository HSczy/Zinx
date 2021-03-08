package main

import (
	"fmt"
	"zinx/ziface"
	"zinx/znet"
)

// ping test 自定义路由
type PingRouter struct {
	znet.BaseRouter
}

// Test Handle
func (this *PingRouter) Handle(req ziface.IRequest) {
	fmt.Println("Call Ping Router Handle..")
	// 先读取客户端的数据 再回写ping...ping...ping...
	fmt.Println("recv from client: MsgID = ", req.GetMsgID(), "\tdata = ", string(req.GetData()))
	if err := req.GetConnection().SendMsg(1, []byte("ping...ping...ping...")); err != nil {
		fmt.Println(err)
	}
}

type HelloRouter struct {
	znet.BaseRouter
}

func (this *HelloRouter) Handle(req ziface.IRequest) {
	fmt.Println("Call Hello Router Handle..")
	// 先读取客户端的数据 再回写ping...ping...ping...
	fmt.Println("recv from client: MsgID = ", req.GetMsgID(), "\tdata = ", string(req.GetData()))
	if err := req.GetConnection().SendMsg(1, []byte("Hello Welcome Zinx")); err != nil {
		fmt.Println(err)
	}
}

// 创建链接之后执行钩子函数
func DoConnectionBegin(conn ziface.IConnection) {
	fmt.Println("====> DoConnectionBegin is Called ...")
	if err := conn.SendMsg(202, []byte("DoConnection Begin")); err != nil {
		fmt.Println(err)
	}
}

func DoConnectionClose(conn ziface.IConnection) {
	fmt.Println("====> DoConnectionClose is Called ...")
}

func main() {
	s := znet.NewServer()
	// 注册链接Hook函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionClose)
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})

	s.Serve()
}
