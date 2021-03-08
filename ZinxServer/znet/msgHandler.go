package znet

import (
	"fmt"
	"strconv"
	"zinx/utlis"
	"zinx/ziface"
)

/*
消息处理模块的实现
*/

type MsgHandle struct {
	// 存放每个MsgID的对应的方法
	Apis map[uint32]ziface.IRouter
	// 负责Worker读取任务的消息队列
	TaskQueue []chan ziface.IRequest
	// 业务工作Worker池的数量
	WorkPoolSize uint32
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:         make(map[uint32]ziface.IRouter),
		WorkPoolSize: utlis.GlobalObject.WorkerPoolSize, // 从全局配置中获取
		TaskQueue:    make([]chan ziface.IRequest, utlis.GlobalObject.WorkerPoolSize),
	}
}

func (mh *MsgHandle) DoMsgHandler(req ziface.IRequest) {
	//1. 从request中找到msgID
	handler, ok := mh.Apis[req.GetMsgID()]
	if !ok {
		fmt.Println("api MsgID = ", req.GetMsgID(), "is NOT FOUND! Need Register! ")
	}
	//2. 根据MsgID调取相对应的方法
	handler.PostHandle(req)
	handler.Handle(req)
	handler.PostHandle(req)
}

func (mh *MsgHandle) AddRouter(msgID uint32, router ziface.IRouter) {
	//1 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := mh.Apis[msgID]; ok {
		// id 已经注册
		panic("repeat apii ,msgID = " + strconv.Itoa(int(msgID)))
	}
	// 2 添加msg与API的绑定关系
	mh.Apis[msgID] = router
	fmt.Println("Add api MsgId = ", msgID, "success! ")
}

// 启动一个Worker 工作池(开启工作池的动作只能发生一次，一个zinx框架只能有一个worker工作池)
func (mh *MsgHandle) StartWorkerPool() {
	// 根据workerPoolSize 分别开启Worker，每个Worker用一个go承载
	for i := 0; i < int(mh.WorkPoolSize); i++ {
		// 一个worker被启动
		// 1 当前的worker对应的channel消息队列开辟空间,第i个worker 就用第i个channal
		mh.TaskQueue[i] = make(chan ziface.IRequest, utlis.GlobalObject.MaxWorkerTaskLen)
		// 2 启动当前的worker，阻塞等待消息从channel 传递过来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

// 启动一个Worker工作流程
func (mh *MsgHandle) StartOneWorker(workID int, taskQueue chan ziface.IRequest) {
	fmt.Println("Worker ID = ", workID, "is start....")
	for {
		select {
		// 如果有消息过来，出列的就是一个客户端的Request，执行当前Request所绑定的业务
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

//将消息交给TaskQueue由Worker进行处理
func (mh *MsgHandle) SendMsgToQueue(req ziface.IRequest) {
	// 1 将消息平均分配给不同的worker ,根据ConnID来分配，
	workerID := req.GetConnection().GetConnID() % mh.WorkPoolSize
	fmt.Println(" Add ConnID", req.GetConnection().GetConnID(), "\t Request MsgID =", req.GetMsgID(), "to WorkID = ", workerID)
	// 2 将消息发送给对应的worker的TaskQueue即可
	mh.TaskQueue[workerID] <- req
}
