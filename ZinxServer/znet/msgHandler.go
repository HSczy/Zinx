package znet

import (
    "fmt"
    "strconv"
    "zinx/ziface"
)

/*
消息处理模块的实现
*/

type MsgHandle struct {
    // 存放每个MsgID的对应的方法
    Apis map[uint32]ziface.IRouter
}

func NewMsgHandle() *MsgHandle {
    return &MsgHandle{
        Apis: make(map[uint32]ziface.IRouter),
    }
}

func (mh *MsgHandle) DoMsgHandler(req ziface.IRequest) {
    //1. 从request中找到msgID
    handler,ok := mh.Apis[req.GetMsgID()]
    if !ok {
        fmt.Println("api MsgID = ", req.GetMsgID(),"is NOT FOUND! Need Register! ")
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
    fmt.Println("Add api MsgId = ",msgID,"success! ")
}
