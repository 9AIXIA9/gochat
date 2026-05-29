package core

// Session 是 Push 上下文对客户端真实物理连接的抽象描述。
// 这里采用领域模型 (纯技术概念的实体)。
type Session interface {
	// ID 当前 Session 的全局唯一标识 (可以是连接 ID)
	ID() string

	// UserID 所关联的最终用户标识
	UserID() string

	// Send 排队推送二进制数据到物理连接
	Send(msg []byte) error

	// Close 关闭该会话/连接
	Close() error
}
