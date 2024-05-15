package interrupter

import (
	"github.com/rs/zerolog/log"
	"kirill5k/reqfol/internal/config"
	"syscall"
	"time"
)

type Interrupter interface {
	StartupTime() time.Time
	Interrupt() bool
}

type signallingInterrupter struct {
	initialDelay time.Duration
	startupTime  time.Time
}

func NewSignallingInterrupter(conf *config.Interrupter) Interrupter {
	return &signallingInterrupter{
		initialDelay: conf.InitialDelay,
		startupTime:  time.Now(),
	}
}

func (si *signallingInterrupter) Interrupt() bool {
	currentTime := time.Now()
	runTime := currentTime.Sub(si.startupTime)
	if runTime > si.initialDelay {
		log.Info().Msg("Sending termination signal to shutdown the app")
		err := syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		if err != nil {
			log.Err(err).Msg("Error sending SIGINT from Interrupter")
		}
		return true
	} else {
		log.Info().Msgf("Delaying termination as the app started %.1f minutes ago", runTime.Minutes())
		return false
	}
}

func (si *signallingInterrupter) StartupTime() time.Time {
	return si.startupTime
}
