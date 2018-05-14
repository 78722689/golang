package utility

import (
	"os"
	"github.com/op/go-logging"
)

var (
	cfic_logger = logging.MustGetLogger("CFICLogger")
	format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{shortfile} %{shortfunc} ? %{level:.4s} %{id:03x}%{color:reset} %{message}`,)
)

func Init_Logger() {
	backend := logging.NewLogBackend(os.Stderr, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	backendLeveled := logging.AddModuleLevel(backend)
	backendLeveled.SetLevel(logging.ERROR, "")

	// Set the backends to be used.
	logging.SetBackend(backendLeveled, backendFormatter)
}

func GetLogger()  *logging.Logger{
	return cfic_logger
}
