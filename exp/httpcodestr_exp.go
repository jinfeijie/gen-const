package exp

type HttpCodeStr string

//go:generate gen-const -type HttpCodeStr
const (
	respOkStr   = HttpCodeStr("0")   // 成功
	respFailStr = HttpCodeStr("500") // 失败
)
