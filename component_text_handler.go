package apexlogutils

import (
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/apex/log"
)

// start time.
var start = time.Now()

// colors.
const (
	none   = 0
	red    = 31
	green  = 32
	yellow = 33
	blue   = 34
	gray   = 37
)

var colors = [...]int{
	log.DebugLevel: gray,
	log.InfoLevel:  blue,
	log.WarnLevel:  yellow,
	log.ErrorLevel: red,
	log.FatalLevel: red,
}

var strings = [...]string{
	log.DebugLevel: "DEBUG",
	log.InfoLevel:  "INFO",
	log.WarnLevel:  "WARN",
	log.ErrorLevel: "ERROR",
	log.FatalLevel: "FATAL",
}

const componentFieldName = "component"

// ComponentTextHandler with additional handling of "component" field for identification of sub components
type ComponentTextHandler struct {
	mu sync.Mutex
	w  io.Writer
}

// NewComponentTextHandler builds new handler.
func NewComponentTextHandler(w io.Writer) *ComponentTextHandler {
	return &ComponentTextHandler{
		w: w,
	}
}

// HandleLog implements log.Handler.
func (h *ComponentTextHandler) HandleLog(e *log.Entry) error {
	color := colors[e.Level]
	level := strings[e.Level]
	names := e.Fields.Names()

	h.mu.Lock()
	defer h.mu.Unlock()

	component := e.Fields.Get(componentFieldName)
	if component == nil {
		component = "global"
	}

	ts := time.Since(start) / time.Second
	_, _ = fmt.Fprintf(h.w, "\033[%dm%6s\033[0m[%04d] \033[1;30m%-10s\033[0m %-25s", color, level, ts, component, e.Message)

	for _, name := range names {
		if name == componentFieldName {
			continue
		}
		_, _ = fmt.Fprintf(h.w, " \033[%dm%s\033[0m=%v", color, name, e.Fields.Get(name))
	}

	_, _ = fmt.Fprintln(h.w)

	return nil
}
