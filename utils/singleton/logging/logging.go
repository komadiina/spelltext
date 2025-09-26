package logging

import (
	"bufio"
	"fmt"
	"io"
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

func Init(level log.Level, name string) {
	once.Do(func() {
		l := *log.NewWithOptions(os.Stdout, log.Options{
			ReportCaller:    true,
			ReportTimestamp: true,
			TimeFormat:      time.DateTime,
			Prefix:          name,
		})

		_, err := os.Stat("var/log")
		if os.IsNotExist(err) {
			os.Mkdir("var/log", 0755)
		}

		fd, _ := os.OpenFile(fmt.Sprintf("var/log/%s-%s.log", name, time.Now().Format(time.RFC3339)), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		l.SetOutput(io.MultiWriter(os.Stdout, bufio.NewWriter(fd)))
		l.SetLevel(level)

		inst = &l
	})
}

func Get(name string) *log.Logger {
	if inst == nil {
		Init(log.InfoLevel, name)
	}
	return inst
}
