package znet

import (
    "fmt"
    "io"
    "net"
    "testing"
)

func TestDataPack(t *testing.T) {
    /*
       模拟的服务器
    */
    // 1. 创建socketTCP
    listener, err := net.Listen("tcp", "127.0.0.1:7777")
    if err != nil {
        fmt.Println("server listen err ", err)
        return
    }
    // 创建一个go 承载 负责从客户端处理业务
    go func() {
        // 2. 从客户端读取数据，拆包处理
        for {
            conn, err := listener.Accept()
            if err != nil {
                fmt.Println("server accept error", err)
            }

            go func(conn net.Conn) {
                // 处理客户端的请求
                //------->拆包的过程<-------
                // 定义一个拆包的对象dp
                dp := NewDataPack()
                for {
                    // 1. 第一次从conn读，把包的head读出来
                    headData := make([]byte, dp.GetHeadLen())
                    _, err := io.ReadFull(conn, headData)
                    if err != nil {
                        fmt.Println("read head error")
                        break
                    }
                    msgHead, err := dp.Unpack(headData)
                    if err != nil {
                        fmt.Println("server unpack err", err)
                        break
                    }
                    if msgHead.GetMsgLen() > 0 {
                        // msg是有数据的，需要进行第二次读取
                        // 2. 第二次从conn读，根据head中的dataLen在读取data内容
                        msg := msgHead.(*Message)
                        msg.Data = make([]byte, msg.GetMsgLen())
                        _, err := io.ReadFull(conn, msg.Data)
                        if err != nil {
                            fmt.Println("server unpack data err:", err)
                            break
                        }
                        // 完整的一个消息已经读取完毕
                        fmt.Println("Message len is :", msg.DataLen, "\tMessage ID is :", msg.Id,
                            "Message data is :", string(msg.Data))
                    }

                }

            }(conn)
        }
    }()

    /*
       模拟客户端
    */
    conn, err := net.Dial("tcp", "127.0.0.1:7777")
    if err != nil {
        fmt.Println("client dial err,", err)
        return
    }

    // 创建一个封包对象dp
    dp := NewDataPack()
    // 模拟黏包过程，封装两个msg一同发送
    // 封包第一个msg1包
    msg1 := &Message{
        Id:      1,
        DataLen: 4,
        Data:    []byte{'z', 'i', 'n', 'x'},
    }
    sendData1, err := dp.Pack(msg1)
    if err != nil {
        fmt.Println("client pack msg1 err ", err)
        return
    }
    // 封包第二个msg2包
    msg2 := &Message{
        Id:      1,
        DataLen: 9,
        Data:    []byte{'h','e','l','l','o','z', 'i', 'n', 'x'},
    }
    sendData2, err := dp.Pack(msg2)
    if err != nil {
        fmt.Println("client pack msg1 err ", err)
        return
    }
    // 将两个包粘结在一起
    sendData1 = append(sendData1,sendData2...)
    // 一次性发给服务端
    _,err = conn.Write(sendData1)
    if err != nil {
        fmt.Println("send message err ", err)
        return
    }
    // 客户端阻塞
    select {}
}
