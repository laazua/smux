package smux

// Handler 消息处理器接口
type Handler interface {
	Handle(conn *Conn) error
}

// HandlerFunc 是一个适配器，允许普通函数作为Handler使用
type HandlerFunc func(conn *Conn) error

// Handle 实现Handler接口
func (fun HandlerFunc) Handle(conn *Conn) error {
	return fun(conn)
}
