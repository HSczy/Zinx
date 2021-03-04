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
	fmt.Println("Call Router Handle..")
	// 先读取客户端的数据 再回写ping...ping...ping...
	fmt.Println("recv from client: MsgID = ", req.GetMsgID(), "\tdata = ", string(req.GetData()))
	if err := req.GetConnection().SendMsg(1, []byte("ping...ping...ping...")); err != nil {
		fmt.Println(err)
	}
}

func main() {
	s := znet.NewServer()
	s.AddRouter(&PingRouter{})
	s.Serve()
}
