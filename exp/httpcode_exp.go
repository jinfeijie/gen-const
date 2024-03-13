package exp

type HttpCode int

//go:generate gen-const -type HttpCode
const (
	respOk   = HttpCode(0)   // 成功
	respFail = HttpCode(500) // 失败
)
