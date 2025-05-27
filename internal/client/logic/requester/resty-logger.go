package requester

import (
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"go-keeper/pkg/logging"
)

var _ resty.Logger = (*RestyLogger)(nil)

type RestyLogger struct {
	logger *logging.ZapLogger
}

func NewRestyLogger(logger *logging.ZapLogger) *RestyLogger {
	return &RestyLogger{
		logger: logger,
	}
}

func (l *RestyLogger) Errorf(format string, v ...interface{}) {
	l.logger.ErrorCtx(context.Background(), fmt.Sprintf(format, v...))
}
func (l *RestyLogger) Warnf(format string, v ...interface{}) {
	l.logger.WarnCtx(context.Background(), fmt.Sprintf(format, v...))
}
func (l *RestyLogger) Debugf(format string, v ...interface{}) {
	l.logger.DebugCtx(context.Background(), fmt.Sprintf(format, v...))
}
