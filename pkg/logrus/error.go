package logrus

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"

	"github.com/pkg/errors"
)

// PkgErrorEntry enables stack frame extraction directly into the log fields.
type PkgErrorEntry struct {
	*logrus.Entry
	// Depth defines how much of the stacktrace you want.
	Depth int
}

// This is dirty pkg/errors.
type stackTracer interface {
	StackTrace() errors.StackTrace
}

func (e *PkgErrorEntry) WithError(err error) *logrus.Entry {
	out := e.Entry

	common := func(pError stackTracer) {
		st := pError.StackTrace()
		depth := len(st)
		if depth > 10 {
			depth = 10
		}

		if e.Depth != 0 {
			depth = e.Depth
		}
		valued := fmt.Sprintf("%+v", st[0:depth])
		valued = strings.Replace(valued, "\t", "", -1)
		stack := strings.Split(valued, "\n")
		out = out.WithField("stack", stack[2:])
	}

	if err2, ok := err.(stackTracer); ok {
		common(err2)
	}

	if err2, ok := errors.Cause(err).(stackTracer); ok {
		common(err2)
	}

	return out.WithError(err)
}
