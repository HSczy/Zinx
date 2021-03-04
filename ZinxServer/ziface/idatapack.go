package ziface

/*
封包拆包的模块
直接面向TCP链接中的数据流，用于处理TCP黏包问题
 */

type IDataPack interface {
    // 获取包的头的长度方法
    GetHeadLen() uint32
    // 封包方法
    Pack(msg IMessage) ([]byte,error)
    // 拆包方法
    Unpack(data []byte) (IMessage,error)
}
