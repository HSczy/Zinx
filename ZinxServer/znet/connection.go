package znet

import (
    "fmt"
    "net"
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

    // 当前的链接状态
    handleAPI ziface.HandleFunc

    // 告知当前链接已经停止的channel
    ExitChan chan bool
}

// 初始化链接模块的方法

func NewConnection(conn *net.TCPConn, id uint32, callbackAPI ziface.HandleFunc) *Connection {
    c := &Connection{
        Conn:      conn,
        ConnID:    id,
        isClosed:  false,
        handleAPI: callbackAPI,
        ExitChan:  make(chan bool, 1),
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
        buf := make([]byte, 512)
        cnt, err := c.Conn.Read(buf)
        if err != nil {
            fmt.Println("recv buf err", err)
            continue
        }
        // 调用当前链接所绑定的HandleAPI
        if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
            fmt.Println("ConnID", c.ConnID, "handle is error", err)
            break
        }
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
    c.Conn.Close()
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
