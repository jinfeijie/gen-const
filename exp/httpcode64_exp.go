package exp

type HttpCode64 int64

//go:generate gen-const -type HttpCode64
const (
	respOk64   = HttpCode64(0)   // 成功
	respFail64 = HttpCode64(500) // 失败
)
