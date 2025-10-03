package logging

import (
	"bufio"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/charmbracelet/log"
)

var (
	once sync.Once
	inst map[string]*log.Logger
)

type Logger = log.Logger

func Init(level log.Level, name string, toFile bool) {
	inst = make(map[string]*log.Logger)

	once.Do(func() {
		l := *log.NewWithOptions(os.Stdout, log.Options{
			ReportCaller:    true,
			ReportTimestamp: true,
			TimeFormat:      time.DateTime,
			Prefix:          name,
		})

		if toFile {
			_, err := os.Stat("var/log")
			if os.IsNotExist(err) {
				os.Mkdir("var/log", 0755)
			}

			fd, _ := os.OpenFile(
				fmt.Sprintf(
					"var/log/%s-%s.log", name, time.Now().Format(time.RFC3339)),
				os.O_APPEND|os.O_CREATE|os.O_WRONLY,
				0644,
			)
			l.SetOutput(bufio.NewWriter(fd))
		} else {
			l.SetOutput(os.Stdout)
		}

		l.SetLevel(level)

		inst[name] = &l
	})
}

func Get(name string, toFile bool) *log.Logger {
	if inst == nil || inst[name] == nil {
		Init(log.InfoLevel, name, toFile)
	}

	return inst[name]
}
