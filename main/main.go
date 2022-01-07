package main

import (
	"bytes"
	"fmt"
	logrus "github.com/sirupsen/logrus"
	"strings"
)

var (
	LogFormatter logrus.Formatter
)

// Log Level Constants
const (
	defaultLogLevel = logrus.InfoLevel
	panicLevel      = "PANC"
	fatalLevel      = "FATL"
	errorLevel      = "ERRO"
	warnLevel       = "WARN"
	infoLevel       = "INFO"
	debugLevel      = "DEBG"
)

type ErrorHook struct {
}

type plainFormatter struct {
	TimestampFormat string
	LevelDesc       []string
}

func (h *ErrorHook) Levels() []logrus.Level {
	// fire only on ErrorLevel (.Error(), .Errorf(), etc.)
	return []logrus.Level{logrus.ErrorLevel}
}

func (h *ErrorHook) Fire(e *logrus.Entry) error {
	// e.Data is a map with all fields attached to entry
	if _, ok := e.Data["severity"]; !ok {
		e.Data["severity"] = "normal"
	}
	return nil
}

func main() {
	log := logrus.New()
	log.SetReportCaller(true)
	formatter := new(plainFormatter)
	formatter.TimestampFormat = "2006-01-02 15:04:05"
	formatter.LevelDesc = []string{
		panicLevel,
		fatalLevel,
		errorLevel,
		warnLevel,
		infoLevel,
		debugLevel}
	log.SetFormatter(formatter)
	log.AddHook(&ErrorHook{})
	log.Info("sample message")
	log.Errorf("fededefdhejgde")
	log.WithFields(logrus.Fields{"severity": "critical", "err_code": 1234}).Errorf("Error message")
}

// formatFilePath retrieves only the last part from a path.
func formatFilePath(path string) string {
	arr := strings.Split(path, "/")
	return arr[len(arr)-1]
}

func createKeyValuePairs(m logrus.Fields) string {
	b := new(bytes.Buffer)
	for key, value := range m {
		fmt.Fprintf(b, "%s=%s ", key, value)
	}
	return b.String()
}

func (f *plainFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := fmt.Sprintf(entry.Time.Format(f.TimestampFormat))
	//entry.Caller.File = "main.go"
	//entry.Caller.Function = "main"
	//entry.Caller.Line = 4
	return []byte(fmt.Sprintf("%s %s [%s:%d] - [%s] [-] %s [%s]\n",
		timestamp,
		f.LevelDesc[entry.Level],
		formatFilePath(entry.Caller.File),
		entry.Caller.Line,
		formatFilePath(entry.Caller.Function),
		entry.Message,
		createKeyValuePairs(entry.Data))), nil
}
