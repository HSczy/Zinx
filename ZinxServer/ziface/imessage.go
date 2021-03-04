package ziface

/*
将强求的消息封装到Message中，定义抽象的接口
 */
type IMessage interface {
    // 获取消息的id
    GetMsgId() uint32
    //获取消息的内容
    GetMsgData() []byte
    //获取消息的长度
    GetMsgLen() uint32
    //设置消息的id
    SetMsgId(uint32)
    //设置消息的内容
    SetMsgData([] byte)
    // 设置消息的长度
    SetMsgLen(uint32)
}