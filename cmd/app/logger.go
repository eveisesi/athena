package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/newrelic/go-agent/v3/integrations/logcontext/nrlogrusplugin"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

func loadLogger(cfg config, command string) (*logrus.Logger, error) {
	logger = logrus.New()

	logger.SetOutput(ioutil.Discard)

	logger.AddHook(&writerHook{
		Writer:    os.Stdout,
		LogLevels: logrus.AllLevels,
	})

	logger.AddHook(&writerHook{
		Writer: &lumberjack.Logger{
			Filename: fmt.Sprintf("logs/error-%s.log", command),
			MaxSize:  10,
			Compress: false,
		},
		LogLevels: []logrus.Level{
			logrus.PanicLevel,
			logrus.FatalLevel,
			logrus.ErrorLevel,
			logrus.WarnLevel,
		},
	})

	logger.AddHook(&writerHook{
		Writer: &lumberjack.Logger{
			Filename:   fmt.Sprintf("logs/info-%s.log", command),
			MaxBackups: 3,
			MaxSize:    10,
			Compress:   false,
		},
		LogLevels: []logrus.Level{
			logrus.InfoLevel,
		},
	})

	level, err := logrus.ParseLevel(cfg.Log.Level)
	if err != nil {
		return logger, errors.Wrap(err, "failed to configure log level")
	}

	logger.SetLevel(level)
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	if cfg.Env == "production" {
		logrus.SetFormatter(nrlogrusplugin.ContextFormatter{})
	}

	return logger, nil
}

type writerHook struct {
	Writer    io.Writer
	LogLevels []logrus.Level
}

func (w *writerHook) Fire(entry *logrus.Entry) error {

	line, err := entry.String()
	if err != nil {
		return err
	}
	_, err = w.Writer.Write([]byte(line))
	return err
}

func (w *writerHook) Levels() []logrus.Level {
	return w.LogLevels
}
