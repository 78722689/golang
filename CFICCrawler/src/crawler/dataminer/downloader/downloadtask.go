package downloader

import "routingpool"

// Implement Task interface
type DownloadTask struct {
	routingpool.Base
}

func NewDownloadTask(name string, call func(int)) *DownloadTask {
	return &DownloadTask{Base : routingpool.Base{Name: name, Call: call, Response: make(chan bool)}}
}

func (c *DownloadTask) Run(id int) {
	c.Call(id)
}
