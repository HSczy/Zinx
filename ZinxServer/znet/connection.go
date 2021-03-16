package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"zinx/utlis"
	"zinx/ziface"
)

// 链接模块
type Connection struct {
	// 当前conn属于哪个server
	TcpServer ziface.IServer
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
	// 链接属性集合
	property map[string]interface{}
	// 保护链接属性的互斥锁
	propertyLock sync.RWMutex
}

// 初始化链接模块的方法

func NewConnection(server ziface.IServer, conn *net.TCPConn, id uint32, msgHandler ziface.IMsgHandler) *Connection {
	c := &Connection{
		TcpServer:  server,
		Conn:       conn,
		ConnID:     id,
		isClosed:   false,
		ExitChan:   make(chan bool, 1),
		msgChan:    make(chan []byte),
		MsgHandler: msgHandler,
		property:   make(map[string]interface{}),
	}
	// 将conn加入到ConnManager中
	c.TcpServer.GetConnMgr().AddConn(c)
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
	// 启动从当前链接写数据的业务
	go c.StartWrite()

	// 按照开发者传进来的hook函数，在链接创建时调用hook函数
	c.TcpServer.CallOnConnStart(c)
}

func (c *Connection) Stop() {
	fmt.Println("Conn Stop()... ConnID=", c.ConnID)

	// 如果当前链接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	//关闭socket链接
	// 按照开发者传进来的hook函数，在链接销毁前调用hook函数
	c.TcpServer.CallOnConnStop(c)

	_ = c.Conn.Close()
	// 告知Writer关闭
	c.ExitChan <- true
	// 将当前链接从ConnMgr中删除
	c.TcpServer.GetConnMgr().DelConn(c)
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

// 设置链接属性
func (c *Connection) SetProperty(key string, value interface{}) {
	// 写锁
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	c.property[key] = value
}

// 获取链接属性
func (c *Connection) GetProperty(key string) (interface{}, error) {
	// 读锁
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()
	if value, ok := c.property[key]; ok {
		return value, nil
	}
	return nil, errors.New("Property :" + key + "not found.")
}

// 移除链接属性
func (c *Connection) RemoveProperty(key string) {
	// 写锁
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.property, key)
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
