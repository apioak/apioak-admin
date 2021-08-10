package enums

const (
	Success = 0  // 成功
	Error   = -1 // 失败
)

var MapMessages = map[int]string{
	Success: "成功",
	Error:   "失败",
}

func CodeMessages(code int) string {
	return MapMessages[code]
}
