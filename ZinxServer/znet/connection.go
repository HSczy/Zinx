package znet

import (
	"errors"
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
	// 无缓冲管道，用于读、写Goroutine之间的消息通信
	msgChan chan []byte
	// 该链接处理的方法Router
	MsgHandler ziface.IMsgHandler
}

// 初始化链接模块的方法

func NewConnection(conn *net.TCPConn, id uint32, msgHandler ziface.IMsgHandler) *Connection {
	c := &Connection{
		Conn:       conn,
		ConnID:     id,
		isClosed:   false,
		ExitChan:   make(chan bool, 1),
		msgChan:    make(chan []byte),
		MsgHandler: msgHandler,
	}
	return c
}

// 链接的读业务方法
func (c *Connection) StartRead() {
	fmt.Println("[Reader Goroutine is running...]")
	defer fmt.Println("connID=", c.ConnID, "Reader is exit ,remote addr is", c.RemoteAddr().String())
	defer c.Stop()

	for {
		// 创建一个拆包解包对象
		dp := NewDataPack()
		// 读取和苦短的Msg Head 二进制流 8个字节
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head error ", err)
			break
		}

		// 拆包，得到msgID 和msg DataLen 放到msg消息中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("package msg error ", err)
			break
		}
		// 根据dataLen 再次读取Data 放在msg.Data
		var msgData []byte
		if msg.GetMsgLen() > 0 {
			msgData = make([]byte, msg.GetMsgLen())
			_, err := io.ReadFull(c.GetTCPConnection(), msgData)
			if err != nil {
				fmt.Println("read msg data error ", err)
				break
			}
		}
		msg.SetMsgData(msgData)

		// 从当前conn数据的Request请求数据
		req := Request{
			conn: c,
			msg:  msg,
		}
		if utlis.GlobalObject.WorkerPoolSize > 0 {
			// 已经开启了工作池，将消息发送给Worker工作池即可
			c.MsgHandler.SendMsgToQueue(&req)
		} else {
			// 从路由中，找到注册绑定过得Conn对应的router调用
			// 根据绑定好的MsgID找到对应的方法 并执行
			go c.MsgHandler.DoMsgHandler(&req)
		}
	}
}

// 写消息的Goroutine，专门发送给客户端消息的模块
func (c *Connection) StartWrite() {
	fmt.Println("[Writer Goroutine is running... ]")
	defer fmt.Println(c.RemoteAddr().String(), " [conn Writer exit!]")
	// 不断的阻塞的等待channel消息，进行写给客户端
	for {
		select {
		case data := <-c.msgChan:
			// 有数据要写给客户端
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data error, ", err)
				return
			}
		case <-c.ExitChan:
			//代表 Reader 已经退出，此事Writer也要退出
			return
		}
	}
}

func (c *Connection) Start() {
	fmt.Println("Conn Start()... ConnID=", c.ConnID)
	// 启动从当前链接的读数据的业务
	go c.StartRead()
	// todo 启动从当前链接协数据的业务
	go c.StartWrite()
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
	// 告知Writer关闭
	c.ExitChan <- true
	// 关闭chan
	close(c.ExitChan)
	close(c.msgChan)
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
func (c *Connection) SendMsg(msgID uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg ")
	}
	// 将data 进行封包 MsgDataLen｜MsgId｜Data
	dp := NewDataPack()
	msg := NewMessage(msgID, data)
	binaryMsg, err := dp.Pack(msg)
	if err != nil {
		fmt.Println("Pack errors msg id = ", msgID)
		return errors.New("Pack error msg ")
	}
	c.msgChan <- binaryMsg
	return nil
}
