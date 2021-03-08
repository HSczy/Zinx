package ziface

/*
消息管理抽象层
*/
type IMsgHandler interface {
	// 调度/执行对应的Router消息处理方法
	DoMsgHandler(req IRequest)
	// 为消息添加具体的处理路由
	AddRouter(msgID uint32, router IRouter)
	// 启动Worker工作池
	StartWorkerPool()
	//将消息交给TaskQueue由Worker进行处理
	SendMsgToQueue(req IRequest)
}
