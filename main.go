// 	go run main.go
// 	go run main.go https://sentry.io
// 	go run main.go bad-url
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
)

// go build -ldflags="-X main.release=VALUE"
var release string

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s URL", os.Args[0])
	}

	err := sentry.Init(sentry.ClientOptions{
		// Your project's DSN in Sentry
		Dsn:     "http://f3a0b3fb522349c1a7e32ed35d8153cd@localhost:9000/2",
		Debug:   true,
		Release: release,
	})

	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	defer sentry.Flush(2 * time.Second)

	resp, err := http.Get(os.Args[1])
	if err != nil {
		sentry.CaptureException(err)
		log.Printf("Reported to Sentry: %s", err)
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			sentry.CaptureException(err)
			log.Printf("Reported to Sentry: %s", err)
			return
		}
	}(resp.Body)

	for header, values := range resp.Header {
		for _, value := range values {
			fmt.Printf("%s=%s\n", header, value)
		}
	}
}
