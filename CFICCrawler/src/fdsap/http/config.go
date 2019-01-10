package http

import "github.com/op/go-logging"

type HttpHandlerConfig struct {
	RoutingNumber int
	RoutingCapacity int

	Log *logging.Logger
}
