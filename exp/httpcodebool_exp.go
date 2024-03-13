package exp

type HttpCodeBool bool

//go:generate gen-const -type HttpCodeBool
const (
	respOkBool   = HttpCodeBool(true)  // 成功
	respFailBool = HttpCodeBool(false) // 失败
)
