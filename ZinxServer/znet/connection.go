package znet

import (
    "fmt"
    "io"
    "net"
    "zinx/utlis"
    "zinx/ziface"
)

// 链接模块
type Connection struct {
    // 当前链接的socket TCP套接字
    Conn *net.TCPConn
    // 链接的ID
    ConnID uint32
    // 当前链接的状态
    isClosed bool
    // 告知当前链接已经停止的channel
    ExitChan chan bool
    // 该链接处理的方法Router
    Router ziface.IRouter
}

// 初始化链接模块的方法

func NewConnection(conn *net.TCPConn, id uint32, router ziface.IRouter) *Connection {
    c := &Connection{
        Conn:     conn,
        ConnID:   id,
        isClosed: false,
        ExitChan: make(chan bool, 1),
        Router:   router,
    }
    return c
}

// 链接的读业务方法
func (c *Connection) StartRead() {
    fmt.Println("Reader Goroutine is running...")
    defer fmt.Println("connID=", c.ConnID, "Reader is exit ,remote addr is", c.RemoteAddr().String())
    defer c.Stop()

    for {
        // 读取客户端的数据到buf中，最大512 字节
        buf := make([]byte, utlis.GlobalObject.MaxPackageSize)
        _, err := c.Conn.Read(buf)
        if err != nil {
            if err == io.EOF {
                fmt.Println("End,", err)
                break
            } else {
                fmt.Println("recv buf err", err)
                continue
            }
        }
        // 从当前conn数据的Request请求数据
        req := Request{
            conn: c,
            data: buf,
        }
        // 从路由中，找到注册绑定过得Conn对应的router调用
        go func(request ziface.IRequest) {
            c.Router.PreHandle(request)
            c.Router.Handle(request)
            c.Router.PostHandle(request)
        }(&req)
    }
}

func (c *Connection) Start() {
    fmt.Println("Conn Start()... ConnID=", c.ConnID)
    // 启动从当前链接的读数据的业务
    go c.StartRead()
    // todo 启动从当前链接协数据的业务
}

func (c *Connection) Stop() {
    fmt.Println("Conn Stop()... ConnID=", c.ConnID)

    // 如果当前链接已经关闭
    if c.isClosed == true {
        return
    }
    c.isClosed = true
    //关闭socket链接
    _ = c.Conn.Close()
    // 关闭chan
    close(c.ExitChan)
}
func (c *Connection) GetTCPConnection() *net.TCPConn {
    return c.Conn
}
func (c *Connection) GetConnID() uint32 {
    return c.ConnID
}
func (c *Connection) RemoteAddr() net.Addr {
    return c.Conn.RemoteAddr()
}
func (c *Connection) Send(data []byte) error {
    return nil
}
