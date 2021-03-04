package znet

import (
    "bytes"
    "encoding/binary"
    "errors"
    "zinx/utlis"
    "zinx/ziface"
)

/*
封包、拆包的具体模块
*/
type DataPack struct{}

// 拆包封包实例的初始化方法

func NewDataPack() *DataPack {
    return &DataPack{}
}

func (dp *DataPack) GetHeadLen() uint32 {
    // DataLen uint32(4字节) + ID uint32(4字节)
    return 8
}

func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
    // 创建一个存放bytes自己的缓冲
    dataBuff := bytes.NewBuffer([]byte{})

    //将dataLen写进dataBuff中
    // binary.LittleEndian 大端写 大端读 小端写小端读
    if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgLen()); err != nil {
        return nil, err
    }
    // 将dataId写进dataBuff中
    if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
        return nil, err
    }
    // 将dataData写进dataBuff中
    if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgData()); err != nil {
        return nil, err
    }
    return dataBuff.Bytes(), nil
}

// 拆包方法 将包的Head的信息读出来 之后再根据Head信息里data长度，在进行一次读
func (dp *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
    // 创建一个从输入的二进制数据的ioReader
    dataBuff := bytes.NewReader(binaryData)

    // 只解压header信息，得到dataLen和Ms个ID
    msg := &Message{

    }
    // 读dataLen
    if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
        return nil, err
    }
    // 读MsgID
    if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
        return nil, err
    }

    if utlis.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utlis.GlobalObject.MaxPackageSize {
       return nil, errors.New("Too Large msg data recv! ")
    }
    return msg,nil
}
