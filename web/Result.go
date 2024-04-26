package web

// Result Result
type Result struct {
	Code int         // 状态
	Msg  string      // 消息
	Data interface{} // 数据
}
