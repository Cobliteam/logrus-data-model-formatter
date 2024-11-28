package formatter

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestLogrusDataModelFormatter(t *testing.T) {
	testCases := []struct {
		name           string
		inputFields    logrus.Fields
		expectedOutput map[string]interface{}
	}{
		{
			name: "Standard DD Fields",
			inputFields: logrus.Fields{
				"dd.trace_id": "1234",
				"dd.span_id":  "12345678",
				"dd.service":  "test-service",
				"dd.version":  "1.0.0",
			},
			expectedOutput: map[string]interface{}{
				"dd": map[string]interface{}{
					"trace_id": "1234",
					"span_id":  "12345678",
					"service":  "test-service",
					"version":  "1.0.0",
				},
				"timestamp": "2024-11-28T14:00:00.000Z",
				"message":   "Test message",
				"level":     "info",
			},
		},
		{
			name: "Custom Fields",
			inputFields: logrus.Fields{
				"user_id": "123",
			},
			expectedOutput: map[string]interface{}{
				"custom": map[string]interface{}{
					"user_id": "123",
				},
				"timestamp": "2024-11-28T14:00:00.000Z",
				"message":   "Test message",
				"level":     "info",
			},
		},
		{
			name: "No DD Fields",
			inputFields: logrus.Fields{
				"dd.trace_id":  "1234",
				"dd.span_id":   "12345678",
				"dd.service":   "test-service",
				"dd.version":   "1.0.0",
				"user_id":      "123",
				"request_type": "api",
			},
			expectedOutput: map[string]interface{}{
				"dd": map[string]interface{}{
					"trace_id": "1234",
					"span_id":  "12345678",
					"service":  "test-service",
					"version":  "1.0.0",
				},
				"custom": map[string]interface{}{
					"user_id":      "123",
					"request_type": "api",
				},
				"timestamp": "2024-11-28T14:00:00.000Z",
				"message":   "Test message",
				"level":     "info",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			formatter := &LogrusDataModelFormatter{}
			entry := &logrus.Entry{
				Logger:  logrus.StandardLogger(),
				Data:    tc.inputFields,
				Time:    time.Date(2024, 11, 28, 14, 0, 0, 0, time.UTC),
				Level:   logrus.InfoLevel,
				Message: "Test message",
			}

			formattedBytes, err := formatter.Format(entry)
			assert.NoError(t, err)

			var output map[string]interface{}
			err = json.Unmarshal(formattedBytes, &output)
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedOutput, output)
		})
	}
}
