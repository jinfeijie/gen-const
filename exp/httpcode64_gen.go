// Code generated by "gen-const -type HttpCode64"; DO NOT EDIT.
// 代码文件由 "gen-const -type HttpCode64" 生成; 不要编辑本代码文件。

package exp

import "fmt"

type HttpCode64Type struct {
	Val int
	Msg string
}

func (receiver *HttpCode64Type) GetVal() int {
	return receiver.Val
}

func (receiver *HttpCode64Type) GetMsg() string {
	return receiver.Msg
}

func HttpCode64Func(val int, msg string) *HttpCode64Type {
	return &HttpCode64Type{
		Val: val,
		Msg: msg,
	}
}

func (receiver *HttpCode64Type) String() string {
	return "HttpCode64Type (Val: " + fmt.Sprintf("%+v", receiver.Val) + ", Msg: " + receiver.Msg + ")"
}

var (
	RespOk64HttpCode64   = HttpCode64Func(0, "成功")
	RespFail64HttpCode64 = HttpCode64Func(500, "失败")
)
