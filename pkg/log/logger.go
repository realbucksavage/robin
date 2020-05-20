package log

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/op/go-logging"
)

var (
	loggingFormat = logging.MustStringFormatter(
		`%{color}[%{time:02/01/2006 15:04:05.000}] [%{shortpkg}/%{shortfile}] [%{level:.4s}]:%{color:reset} %{message}`,
	)
	moduleName = "robin"

	L         = logging.MustGetLogger(moduleName)
	StdLogger = log.New(new(stdLogWriter), "", log.Lshortfile)
)

type stdLogWriter struct{}

func (s stdLogWriter) Write(p []byte) (n int, err error) {
	return fmt.Printf("[%s] [STD]: %s", time.Now().Format("02/01/2006 15:04:05.000"), string(p))
}

func init() {
	b := logging.NewLogBackend(os.Stdout, "", 0)
	logging.SetBackend(b)
	logging.SetFormatter(loggingFormat)
}

func SetLevel(level string) {
	l, err := logging.LogLevel(level)
	if err != nil {
		L.Warningf("%s is not a valid logging level", level)
		return
	}

	logging.SetLevel(l, moduleName)
	L.Debugf("Logging level set to %s", level)
}
