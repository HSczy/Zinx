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
	// 当前server的链接管理器
	ConnMgr ziface.IConnManager
	// 该Server 创建之后自动调用的Hook函数
	OnConnStart func(conn ziface.IConnection)
	// 该Server 销毁之前自动调用的Hook函数
	OnConnStop func(conn ziface.IConnection)
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
			// 判断当前已有链接数是否超过最大链接个数
			// 如果超过就关闭此链接
			if s.ConnMgr.LenConn() >= utlis.GlobalObject.MaxConn {
				// TODO 客户端响应一个超出最大链接的错误包
				fmt.Println("Too Many Connections MxcConn = ", utlis.GlobalObject.MaxConn, "\t Connections now is ", s.ConnMgr.LenConn())
				_ = conner.Close()
				continue
			}

			// 将处理新链接的业务方 法和conn进行绑定 得到我们的链接模块
			dealConn := NewConnection(s, conner, cid, s.MsgHandler)
			cid++
			go dealConn.Start()
		}
	}()

}

// 停止服务器
func (s *Server) Stop() {
	// 将一些服务器的资源、状态或者一些已经开辟的链接信息进行停止或回收
	fmt.Println("[STOP] Zinx server name ", s.Name)
	s.ConnMgr.StopAllConn()
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

func (s *Server) GetConnMgr() ziface.IConnManager {
	return s.ConnMgr
}

// 注册OnConnStart钩子函数的方法
func (s *Server) SetOnConnStart(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStart = hookFunc
}

// 注册OnConnStop钩子函数的方法
func (s *Server) SetOnConnStop(hookFunc func(connection ziface.IConnection)) {
	s.OnConnStop = hookFunc
}

// 调用OnConnStart钩子函数的方法
func (s *Server) CallOnConnStart(connection ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("---> Call OnConnStart... ")
		s.OnConnStart(connection)
	}
}

// 调用OnConnStop钩子函数的方法
func (s *Server) CallOnConnStop(connection ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("---> Call OnConnStop... ")
		s.OnConnStop(connection)
	}
}

//初始化Server模块的方法
func NewServer() ziface.IServer {
	s := &Server{
		Name:       utlis.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utlis.GlobalObject.Host,
		Port:       utlis.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(),
	}
	return s
}
