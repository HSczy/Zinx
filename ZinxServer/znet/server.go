package znet

import (
    "errors"
    "fmt"
    "net"
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
}

func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
    fmt.Println("[Conn Handle] CallBackToClient...")
    if _, err := conn.Write(data[:cnt]); err != nil {
        fmt.Println("write back buf err", err)
        return errors.New("CallBackToClient error")
    }
    return nil
}

// 启动服务器
func (s *Server) Start() {
    fmt.Printf("[Start] Server Listenner at IP :%s,Port %d, is starting\n", s.IP, s.Port)
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
            // 将处理新链接的业务方法和conn进行绑定 得到我们的链接模块
            dealConn := NewConnection(conner, cid, CallBackToClient)
            cid++
            go dealConn.Start()
            //// 已经与客户端建立连接，做一些简单的业务，最大512字节的回显业务
            //go func() {
            //    // 不断的从客户端获取数据
            //    for {
            //        buf := make([]byte, 512)
            //        cnt, err := conner.Read(buf)
            //        if err != nil {
            //            fmt.Println("recv buf err ", err)
            //            continue
            //        }
            //        // 回显功能
            //        fmt.Printf("receive data:%s,cnt=%d\n", buf, cnt)
            //        if _, err := conner.Write(buf[:cnt]); err != nil {
            //            fmt.Println("write back buf err ", err)
            //            continue
            //        }
            //    }
            //}()
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
    // 阻塞状态
    select {}
}

//初始化Server模块的方法
func NewServer(name string) ziface.IServer {
    s := &Server{
        Name:      name,
        IPVersion: "tcp4",
        IP:        "0.0.0.0",
        Port:      8999,
    }
    return s
}
