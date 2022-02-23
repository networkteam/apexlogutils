package pgx_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/apex/log"
	"github.com/apex/log/handlers/memory"
	"github.com/jackc/pgx/v4"

	logutilspgx "github.com/networkteam/apexlogutils/pgx"
)

func TestLogger_Log(t *testing.T) {
	type args struct {
		level pgx.LogLevel
		msg   string
		data  map[string]interface{}
	}
	tests := []struct {
		name     string
		args     args
		opts     []logutilspgx.LoggerOpt
		expected *log.Entry
	}{
		{
			name: "pgx trace is logged as debug",
			args: args{
				level: pgx.LogLevelTrace,
				msg:   "Hey, it's a test",
				data: map[string]interface{}{
					"foo": "bar",
				},
			},
			expected: &log.Entry{
				Level:   log.DebugLevel,
				Message: "Hey, it's a test",
				Fields:  log.Fields{"foo": "bar", "component": "db.driver", "level": "trace"},
			},
		},
		{
			name: "pgx debug is logged as debug",
			args: args{
				level: pgx.LogLevelDebug,
				msg:   "Hey, it's a test",
				data: map[string]interface{}{
					"foo": "bar",
				},
			},
			expected: &log.Entry{
				Level:   log.DebugLevel,
				Message: "Hey, it's a test",
				Fields:  log.Fields{"foo": "bar", "component": "db.driver", "level": "debug"},
			},
		},
		{
			name: "pgx info is logged as debug",
			args: args{
				level: pgx.LogLevelInfo,
				msg:   "Hey, it's a test",
				data: map[string]interface{}{
					"foo": "bar",
				},
			},
			expected: &log.Entry{
				Level:   log.DebugLevel,
				Message: "Hey, it's a test",
				Fields:  log.Fields{"foo": "bar", "component": "db.driver", "level": "info"},
			},
		},
		{
			name: "pgx warn is logged as warn",
			args: args{
				level: pgx.LogLevelWarn,
				msg:   "Hey, it's a test",
				data: map[string]interface{}{
					"foo": "bar",
				},
			},
			expected: &log.Entry{
				Level:   log.WarnLevel,
				Message: "Hey, it's a test",
				Fields:  log.Fields{"foo": "bar", "component": "db.driver", "level": "warn"},
			},
		},
		{
			name: "pgx error is logged as error and included in fields as error",
			args: args{
				level: pgx.LogLevelError,
				msg:   "Hey, there was an error",
				data: map[string]interface{}{
					"err": fmt.Errorf("test error"),
					"sql": "SELECT * FROM users",
				},
			},
			expected: &log.Entry{
				Level:   log.ErrorLevel,
				Message: "Hey, there was an error",
				Fields:  log.Fields{"component": "db.driver", "level": "error", "error": "test error", "sql": "SELECT * FROM users"},
			},
		},
		{
			name: "pgx error is ignored if it matches the ignore error option",
			opts: []logutilspgx.LoggerOpt{
				logutilspgx.WithIgnoreErrors(func(err error) bool {
					return err.Error() == "ignored error"
				}),
			},
			args: args{
				level: pgx.LogLevelError,
				msg:   "Hey, there was an error",
				data: map[string]interface{}{
					"err": fmt.Errorf("ignored error"),
					"sql": "SELECT * FROM users",
				},
			},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := memory.New()
			logger := &log.Logger{
				Handler: handler,
				Level:   log.DebugLevel,
			}

			p := logutilspgx.NewLogger(logger, tt.opts...)
			p.Log(context.Background(), tt.args.level, tt.args.msg, tt.args.data)

			if tt.expected == nil {
				if len(handler.Entries) > 0 {
					t.Errorf("expected no log entries, got: %d", len(handler.Entries))
				}
				return
			}

			if len(handler.Entries) != 1 {
				t.Errorf("Expected 1 entry, got %d", len(handler.Entries))
			}
			entry := handler.Entries[0]
			if entry.Level != tt.expected.Level {
				t.Errorf("Expected level %s, got %s", tt.expected.Level, entry.Level)
			}
			if entry.Message != tt.expected.Message {
				t.Errorf("Expected message %s, got %s", tt.expected.Message, entry.Message)
			}
			if len(entry.Fields) != len(tt.expected.Fields) {
				t.Errorf("Expected %d fields, got %d: %v", len(tt.expected.Fields), len(entry.Fields), entry.Fields)
			}
			for k, v := range tt.expected.Fields {
				if entry.Fields[k] != v {
					t.Errorf("Expected field %s to be %s, got %s", k, v, entry.Fields[k])
				}
			}
		})
	}
}
