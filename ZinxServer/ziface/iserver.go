package ziface

// 为znet/server的抽象层

//定义一个服务器接口
type IServer interface {
	// 启动服务器
	Start()
	// 停止服务器
	Stop()
	// 运行服务器
	Serve()
	//路由功能：给当前服务器添加路由方法，共客户端的链接处理使用
	AddRouter(router IRouter)
}
