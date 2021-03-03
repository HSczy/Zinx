package znet

import "zinx/ziface"

// 这里之所以BaseRouter的方法都为空
// 是因为有的router不希望有PreHandle和PostHandle这两个业务
// 所以Router全部继承BaseRouter的好处就是，不需要实现PreHandle和PostHandle方法
// 实现router时，先嵌入BaseRouter基类，然后再根据需要对这个 基类的方法进行重写即可
type BaseRouter struct{}

func (br *BaseRouter) PreHandle(request ziface.IRequest) {}

// 在处理conn业务的主方法Hook
func (br *BaseRouter) Handle(request ziface.IRequest) {}

// 在处理conn业务之后的钩子方法Hook
func (br *BaseRouter) PostHandle(request ziface.IRequest) {}
