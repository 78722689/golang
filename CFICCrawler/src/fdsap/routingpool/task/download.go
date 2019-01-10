package task

import "fdsap/routingpool"

type HttpGetTask struct {
	routingpool.Base
}

func NewHttpGetTask(name string, call func(int)) *HttpGetTask {
	return &HttpGetTask{Base : routingpool.Base{Name: name, Call: call, Response: make(chan bool)}}
}

func (c *HttpGetTask) Run(id int) {
	c.Call(id)
}
