package ziface

/*
链接管理模块，
*/
type IConnManager interface {
	// 添加链接
	AddConn(conn IConnection)
	// 删除链接
	DelConn(coon IConnection)
	// 根据ConnID 获取链接
	GetConn(connID uint32) (conn IConnection, err error)
	// 所有链接的个数
	LenConn() int
	// 停止所有链接
	StopAllConn()
}
