package logger

import (
	"bambamload/constant"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/sirupsen/logrus"
)

var (
	Logger        = logrus.New()
	RequestLogger = logrus.New()
)

func InitLogger(appEnv string) {

	Logger.SetLevel(logrus.InfoLevel)
	Logger.SetReportCaller(true)

	Logger.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint: true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return "", fmt.Sprintf("%s:%d", filepath.Base(f.File), f.Line)
		},
	})

	logWriter, err := NewDailyFileWriter(constant.ErrorLogsDir, constant.Messages)
	if err != nil {
		Logger.Fatalf("Failed to create daily file writer: %v", err)
	}

	if appEnv == constant.Production {
		Logger.SetOutput(logWriter)
	} else {
		multiWriter := io.MultiWriter(os.Stdout, logWriter)
		Logger.SetOutput(multiWriter)
	}

	//set up request logger
	RequestLogger.SetLevel(logrus.InfoLevel)
	RequestLogger.SetReportCaller(true)

	RequestLogger.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint: true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return "", fmt.Sprintf("%s:%d", filepath.Base(f.File), f.Line)
		},
	})

	requestLogWriter, err := NewDailyFileWriter(constant.RequestLogsDir, constant.Messages)
	if err != nil {
		RequestLogger.Fatalf("Failed to create daily file writer: %v", err)
	}

	requestMultiWriter := io.MultiWriter(requestLogWriter)
	RequestLogger.SetOutput(requestMultiWriter)
}
