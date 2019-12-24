package logging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
)

var (
	timeFormat = "2006-01-02 15:04:05.000"
	logFormat  = "%s %s %s [%s] "
)

type formatter struct{}

var _ logrus.Formatter = (*formatter)(nil)

func (formatter) Format(entry *logrus.Entry) ([]byte, error) {
	level := strings.ToUpper(entry.Level.String())
	time := entry.Time.Format(timeFormat)
	fileLineno := entry.Data[fileLinenoStr]
	hostname := entry.Data[hostnameStr]
	delete(entry.Data, fileLinenoStr)
	delete(entry.Data, hostnameStr)

	buffer := bytes.NewBufferString(fmt.Sprintf(logFormat, time, level, hostname, fileLineno))
	if extra, ok := entry.Data[extraStr]; ok {
		extraJson, err := json.Marshal(extra)
		if err == nil {
			buffer.WriteString("[")
			buffer.WriteString(string(extraJson))
			buffer.WriteString("] ")
		}
	}
	buffer.WriteString(entry.Message)
	for key, value := range entry.Data {
		buffer.WriteString(fmt.Sprintf(", %s=%v", key, value))
	}
	buffer.WriteByte('\n')
	return buffer.Bytes(), nil
}

func NewFormatter() *formatter {
	return &formatter{}
}
