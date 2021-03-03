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

// Test PreHandle

func (this *PingRouter) PreHandle(req ziface.IRequest) {
	fmt.Println("Call Router PreHandle..")
	_, err := req.GetConnection().GetTCPConnection().Write([]byte("before ping...\n"))
	if err != nil {
		fmt.Println("call back before ping error")
	}
}

// Test Handle
func (this *PingRouter) Handle(req ziface.IRequest) {
	fmt.Println("Call Router Handle..")
	_, err := req.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping...\n"))
	if err != nil {
		fmt.Println("call backping...ping...ping...error")
	}
}

// Test PostHadnle
func (this *PingRouter) PostHandle(req ziface.IRequest) {
	fmt.Println("Call Router PostHandle..")
	_, err := req.GetConnection().GetTCPConnection().Write([]byte("after ping...\n"))
	if err != nil {
		fmt.Println("call back after ping error")
	}
}

func main() {
	s := znet.NewServer("[zinx V0.3]...")
	s.AddRouter(&PingRouter{})
	s.Serve()
}
