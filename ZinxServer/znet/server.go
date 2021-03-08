package znet

import (
	"fmt"
	"net"
	"zinx/utlis"
	"zinx/ziface"
)

// 为ziface/Iserver 的实例层

// IServer的接口实现，定义一个Server的服务器模块

type Server struct {
	// 服务器的名称
	Name string
	// 服务器绑定的IP版本
	IPVersion string
	// 服务器绑定的IP地址
	IP string
	// 服务器绑定的IP端口
	Port int
	// 当前的server的消息管理模块，用来绑定MsgID和对应的处理业务API关系
	MsgHandler ziface.IMsgHandler
}

// 启动服务器
func (s *Server) Start() {
	fmt.Printf("[Zinx] ServerName:%s\n Listener at IP : %s\nPort:%d is Starting...\n",
		utlis.GlobalObject.Name, utlis.GlobalObject.Host, utlis.GlobalObject.TcpPort)
	go func() {
		// 1. 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error :", err)
			return
		}
		// 2. 监听服务器的地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, "err", err)
			return
		}

		fmt.Println("start Zinx server succ,", s.Name, "success,Listening.....")
		var cid uint32
		cid = 0
		// 3. 阻塞的等待客户端链接，处理客户端链接业务（读写）
		for {
			// 如果有客户端链接过来，阻塞返回
			conner, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}
			// 将处理新链接的业务方 法和conn进行绑定 得到我们的链接模块
			dealConn := NewConnection(conner, cid, s.MsgHandler)
			cid++
			go dealConn.Start()
		}
	}()

}

// 停止服务器
func (s *Server) Stop() {
	// TODO 将一些服务器的资源、状态或者一些已经开辟的链接信息进行停止或回收
}

// 运行服务器
func (s *Server) Serve() {
	// 启动server的服务功能
	s.Start()

	// TODO 做一些启动服务之后的额外业务
	// 启动消息队列
	s.MsgHandler.StartWorkerPool()
	// 阻塞状态
	select {}
}

// 添加路由方法
func (s *Server) AddRouter(MsgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(MsgID, router)
	fmt.Println("Add router Success!")
	// TODO 做一些启动服务之后的额外业务
	return
}

//初始化Server模块的方法
func NewServer() ziface.IServer {
	s := &Server{
		Name:       utlis.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utlis.GlobalObject.Host,
		Port:       utlis.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
	}
	return s
}
