package http

import (
	"fdsap/routingpool"
	"fdsap/routingpool/task"
)

type HttpGetInfo struct {
	Url string
	Path string
	Overwrite bool
}

type HttpHandler struct {
	httpRoutingPool *routingpool.ThreadPool
	cfg *HttpHandlerConfig
}

func (handler *HttpHandler) Wait() {
	handler.httpRoutingPool.Wait()
}

func (handler *HttpHandler) Get(info *HttpGetInfo) {
	caller := func(id int) {
		r := &Request{Url: info.Url, File: info.Path, OverWrite: info.Overwrite}
		_, err := r.Get()
		if err != nil {
			//logger.Errorf("Request failure, %s", err)
			return
		}
	}

	handler.httpRoutingPool.PutTask(task.NewHttpGetTask("HTTP-HttpHandler", caller))
}

func NewHttpHandler(cfg *HttpHandlerConfig) (*HttpHandler, error) {

	handler := &HttpHandler{}
	handler.httpRoutingPool = routingpool.GetPool(128, 128)
	handler.cfg = cfg
	//handler.context = context

	return handler,nil
}