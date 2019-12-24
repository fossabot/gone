package logging

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestDefaultFields(t *testing.T) {
	log := NewLog()
	assert.Equal(t, log.defaultFields(logrus.DebugLevel)[fileLinenoStr], "unknown:0")
	assert.Equal(t, log.defaultFields(logrus.InfoLevel)[fileLinenoStr], "unknown:0")

	assert.NotEqual(t, log.defaultFields(logrus.WarnLevel)[fileLinenoStr], "unknown:0")
	assert.NotEqual(t, log.defaultFields(logrus.ErrorLevel)[fileLinenoStr], "unknown:0")
}

func BenchmarkLogWithoutFile(b *testing.B) {
	log := NewLog()
	log.SetOutput(ioutil.Discard)
	extra := map[string]interface{}{"Status Code": 200, "Remote Address:": "11.11.11.11:443", "Request Method": "GET", "Referrer Policy": "no-referrer-when-downgrade"}
	for n := 0; n < b.N; n++ {
		log.Info(extra)
	}
}

func BenchmarkLogWithFile(b *testing.B) {
	log := NewLog()
	log.SetOutput(ioutil.Discard)
	extra := map[string]interface{}{"Status Code": 200, "Remote Address:": "11.11.11.11:443", "Request Method": "GET", "Referrer Policy": "no-referrer-when-downgrade"}
	for n := 0; n < b.N; n++ {
		log.Warn(extra)
	}
}
