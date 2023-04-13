package log

import (
	"os"
	"sync"

	"github.com/rs/zerolog"
)

var (
	loggerInit sync.Once
	logger     zerolog.Logger
)

func Init() *zerolog.Logger {
	loggerInit.Do(func() {
		logger = zerolog.New(os.Stderr).With().Timestamp().Logger()
	})
	return &logger
}
