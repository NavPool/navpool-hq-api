package logger

import (
	"github.com/getsentry/raven-go"
	"log"
)

func LogError(err error) {
	log.Print(err)
	raven.CaptureErrorAndWait(err, nil)
}
