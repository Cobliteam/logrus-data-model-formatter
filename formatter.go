// Package formatter provides a custom Logrus formatter for structuring logs into
// a data model.
package formatter

import (
	"strings"

	"github.com/sirupsen/logrus"
)

// LogrusDataModelFormatter is a custom formatter for Logrus that organizes log fields
// into structured groups and formats the output as JSON.
type LogrusDataModelFormatter struct {
	BaseFormatter logrus.Formatter
}

func (l *LogrusDataModelFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	dataCopy := make(logrus.Fields)
	for k, v := range entry.Data {
		dataCopy[k] = v
	}

	ddFields := make(map[string]interface{})
	customFields := make(map[string]interface{})
	for k, v := range dataCopy {
		if key, found := strings.CutPrefix(k, "dd."); found {
			ddFields[key] = v
			delete(dataCopy, k)
		} else {
			customFields[k] = v
			delete(dataCopy, k)
		}
	}

	if len(ddFields) > 0 {
		dataCopy["dd"] = ddFields
	}
	if len(customFields) > 0 {
		dataCopy["custom"] = customFields
	}

	jsonFormatter := &logrus.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05.000Z", // ISO 8601 format with milliseconds.
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime: "timestamp",
			logrus.FieldKeyMsg:  "message",
		},
		DisableHTMLEscape: true,
	}
	return jsonFormatter.Format(&logrus.Entry{
		Logger:  entry.Logger,
		Data:    dataCopy,
		Time:    entry.Time.UTC(),
		Level:   entry.Level,
		Message: entry.Message,
	})
}
