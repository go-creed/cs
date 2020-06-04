package log

import (
	"encoding/json"
	"io"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type general struct {
	Message string    `json:"message"`
	Time    time.Time `json:"timestamp"`
}

func (g general) Marshal(entry *logrus.Entry) io.Reader {
	g.Message = entry.Message
	g.Time = entry.Time
	marshal, _ := json.Marshal(g)
	return strings.NewReader(string(marshal))
}
