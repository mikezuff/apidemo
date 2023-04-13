package config

import (
	"github.com/mikezuff/apidemo/pkg/log"
	"github.com/rs/zerolog"
)

type AppContext struct {
	Logger *zerolog.Logger
}

func InitAppContext() *AppContext {
	return &AppContext{
		Logger: log.Init(),
	}
}
