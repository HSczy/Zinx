package znet

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"zinx/ziface"
)

/*
链接管理实现模块
*/

type ConnManager struct {
	// 已经创建的Connection 的map集合
	connMap map[uint32]ziface.IConnection
	// 针对map的互斥锁 保护链接集合
	connLock sync.RWMutex
}

// 创建当前链接的方法

func NewConnManager() *ConnManager {
	return &ConnManager{
		connMap: make(map[uint32]ziface.IConnection),
	}
}

// 添加链接
func (cm *ConnManager) AddConn(conn ziface.IConnection) {
	// 保护共享资源map，加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	// 将conn 加入到 ConnManager 中
	cm.connMap[conn.GetConnID()] = conn
	fmt.Println("connId = ", conn.GetConnID(), "add to ConnManager successfully: conn num now is ", cm.LenConn())
}

// 删除链接
func (cm *ConnManager) DelConn(conn ziface.IConnection) {
	// 保护共享资源map，加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	// 删除链接
	delete(cm.connMap, conn.GetConnID())
	fmt.Println("ConnManager delete connId = ", conn.GetConnID(), "  successfully: conn num now is ", cm.LenConn())
}

// 根据ConnID 获取链接
func (cm *ConnManager) GetConn(connID uint32) (conn ziface.IConnection, err error) {
	// 保护共享资源map，加读锁
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()
	if conn, ok := cm.connMap[connID]; ok {
		// 找到数据
		return conn, nil
	} else {
		return nil, errors.New("conn id =" + strconv.Itoa(int(connID)) + "is not found! ")
	}
}

// 所有链接的个数
func (cm *ConnManager) LenConn() int {
	return len(cm.connMap)
}

// 停止所有链接
func (cm *ConnManager) StopAllConn() {
	// 保护共享资源map，加写锁
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	for connId, conn := range cm.connMap {
		conn.Stop()
		delete(cm.connMap, connId)
	}
	fmt.Println("Clear All connections success! conn num ", cm.LenConn())
}
