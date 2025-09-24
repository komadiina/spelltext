package logging

import (
	"os"
	"sync"
	"time"

	"github.com/charmbracelet/log"
)

var (
	once sync.Once
	inst *log.Logger
)

type Logger = log.Logger

func Init(level log.Level) {
	once.Do(func() {
		l := *log.NewWithOptions(os.Stdout, log.Options{
			ReportCaller:    true,
			ReportTimestamp: true,
			TimeFormat:      time.DateTime,
		})

		l.SetOutput(os.Stdout)
		l.SetLevel(level)

		inst = &l
	})
}

func Get() *log.Logger {
	if inst == nil {
		Init(log.InfoLevel)
	}
	return inst
}
