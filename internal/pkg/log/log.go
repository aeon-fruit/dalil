package log

import (
	"fmt"
	"io"
	"time"

	"github.com/aeon-fruit/dalil.git/internal/pkg/config"
	"github.com/go-logr/logr"
	"github.com/go-logr/zerologr"
	"github.com/rs/zerolog"
)

type Logger interface {
	Enabled() bool
	Error(err error, msg string, keysAndValues ...interface{})
	GetSink() logr.LogSink
	Info(msg string, keysAndValues ...interface{})
	V(level int) logr.Logger
	WithCallDepth(depth int) logr.Logger
	WithCallStackHelper() (func(), logr.Logger)
	WithName(name string) logr.Logger
	WithSink(sink logr.LogSink) logr.Logger
	WithValues(keysAndValues ...interface{}) logr.Logger
	Print(v ...interface{})
}

type loggerImpl struct {
	logr.Logger
}

func (l loggerImpl) Print(v ...interface{}) {
	msg := fmt.Sprint(v...)
	l.V(0).Info(msg)
}

func New(appConfig config.AppConfig, w io.Writer) Logger {
	zerologr.SetMaxV(2)

	var logger zerolog.Logger
	if appConfig.AppEnv == config.AppEnvLocal {
		logger = zerolog.New(zerolog.ConsoleWriter{Out: w, TimeFormat: time.RFC3339})
	} else {
		logger = zerolog.New(w)
	}
	logger = logger.With().Caller().Timestamp().Stack().Logger()

	lg := logr.New(logSink{
		impl:  *zerologr.NewLogSink(&logger).WithCallDepth(4).(*zerologr.LogSink),
		level: appConfig.Logging.GetGlobalVerbosity(),
		cfg:   appConfig.Logging,
	})

	return loggerImpl{
		Logger: lg,
	}
}

type logSink struct {
	impl  zerologr.LogSink
	level int
	cfg   config.LoggingConfig
}

func (ls logSink) Init(info logr.RuntimeInfo) {
	ls.impl.Init(info)
}

func (ls logSink) Enabled(level int) bool {
	return ls.impl.Enabled(level)
}

func (ls logSink) Info(level int, msg string, keysAndValues ...interface{}) {
	if level <= ls.level {
		ls.impl.Info(level, msg, keysAndValues...)
	}
}

func (ls logSink) Error(err error, msg string, keysAndValues ...interface{}) {
	ls.impl.Error(err, msg, keysAndValues...)
}

func (ls logSink) WithValues(keysAndValues ...interface{}) logr.LogSink {
	return logSink{
		impl:  *ls.impl.WithValues(keysAndValues...).(*zerologr.LogSink),
		level: ls.level,
		cfg:   ls.cfg,
	}
}

func (ls logSink) WithName(name string) logr.LogSink {
	return logSink{
		impl:  *ls.impl.WithName(name).(*zerologr.LogSink),
		level: ls.cfg.GetVerbosity(name),
		cfg:   ls.cfg,
	}
}

func (ls logSink) WithCallDepth(depth int) logr.LogSink {
	return logSink{
		impl:  *ls.impl.WithCallDepth(depth).(*zerologr.LogSink),
		level: ls.level,
		cfg:   ls.cfg,
	}
}
