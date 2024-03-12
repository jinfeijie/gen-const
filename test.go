package main

type HttpCode int

//go:generate gen-const -type HttpCode -linecomment -output httpcode.go
const (
	respOk   = HttpCode(0)   // 成功
	respFail = HttpCode(500) // 失败
)

// GetInt()
// GetStr()
