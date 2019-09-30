package utils

import (
	"time"

	"github.com/getsentry/sentry-go"
)

func ConfigureReportingEngine() {
	sentry.Init(sentry.ClientOptions{
		Dsn:   "https://0bdafa2142ab49919874b76ba0ce7379@sentry.io/1757523",
		Debug: true,
	})
}

func Report(err error) {
	sentry.CaptureException(err)
}

func ReportSync(err error) {
	Report(err)
	sentry.Flush(time.Second * 5)
}
