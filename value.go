package main

type Value struct {
	originalName string
	Val          interface{}
	Msg          string
	signed       bool
	str          string
}

func (v *Value) String() string {
	return v.str
}
