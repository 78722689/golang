package utility

import (
	"os"
	"github.com/op/go-logging"
	"github.com/spf13/viper"
	"io"
)

var (
	cfic_logger = logging.MustGetLogger("CFICLogger")
	format = logging.MustStringFormatter(
		`%{color}%{time:15:04:05.000} %{shortfile} %{shortfunc} ? %{level:.4s}%{color:reset} %{message}`,)
)

func Init_Logger() {
	logfile := viper.GetString("logger.logfile")

	var out io.Writer
	if logfile != "" {
		out, _ = os.Create("d:\\cficcrawler.log")

	}
	backend := logging.NewLogBackend(out, "", 0)
	backendFormatter := logging.NewBackendFormatter(backend, format)
	backendLeveled := logging.AddModuleLevel(backend)
	backendLeveled.SetLevel(logging.ERROR, "")

	// Set the backends to be used.
	logging.SetBackend(backendLeveled, backendFormatter)
}

func GetLogger()  *logging.Logger{
	return cfic_logger
}
