package ziface

/*
IRequest接口：
实际上是把客户端的请求的链接信息，和请求的数据包装到一个Request中
 */

type IRequest interface {
    // 得到当前链接
    GetConnection() IConnection
    // 得到请求的数据
    GetData() []byte
}
