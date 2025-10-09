package logging

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/log"
)

var (
	once sync.Once
	inst map[string]*log.Logger
	fd   map[string]*os.File
)

type Logger = log.Logger

func Init(level log.Level, name string, toFile bool) {
	inst = make(map[string]*log.Logger)
	fd = make(map[string]*os.File)

	once.Do(func() {
		l := *log.NewWithOptions(os.Stdout, log.Options{
			ReportCaller:    true,
			ReportTimestamp: true,
			TimeFormat:      time.DateTime,
		})

		if toFile {
			_, err := os.Stat("logs")
			if os.IsNotExist(err) {
				os.MkdirAll("logs", 0755)
			}

			t := time.Now().Format(time.DateTime)
			t = strings.ReplaceAll(t, ":", "-")
			t = strings.ReplaceAll(t, " ", "_")
			filename := fmt.Sprintf("logs/%s-%s.log", name, t)

			f, err := os.Create(filename)
			if err != nil {
				log.Fatalf("failed to create file: %v", err)
			}

			f, err = os.OpenFile(
				filename,
				os.O_CREATE|os.O_APPEND|os.O_WRONLY,
				0644,
			)
			if err != nil {
				log.Fatalf("error opening file: %v", err)
			}

			fd[name] = f
			l.SetOutput(f)
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

func Close(name string) {
	fd[name].Close()
}
