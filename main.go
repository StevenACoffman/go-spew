// Sometimes we'd like our Go programs to intelligently
// handle [Unix signals](http://en.wikipedia.org/wiki/Unix_signal).
// For example, we might want a server to gracefully
// shutdown when it receives a `SIGTERM`, or a command-line
// tool to stop processing input if it receives a `SIGINT`.
//
// [Timers](timers) are for when you want to do
// something once in the future - _tickers_ are for when
// you want to do something repeatedly at regular
// intervals.

package main

import (
	"bytes"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"text/template"
	"time" // or "runtime"

	"github.com/satori/go.uuid"
	"strconv"
)


func check(err error) {
	if err != nil {
		fmt.Printf("Error: %+v. \n", err)
		//os.Exit(1)
		panic(err)
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}


func makePayload(data map[string]string) string {
	const payloadTmpl = `
{
  "dests": [
    "k8s-pirate-message"
  ],
  "_lb0": {
    "k8s-pirate-message": "{{.LOG_BUFFER_PARTITION}}"
  },
  "eventid": "{{.EVENT_ID}}",
  "requestid": "{{.REQUEST_ID}}",
  "origin": "local-pirate.{{.ENVIRONMENT}}",
  "eventtype": "foo",
  "tstamp_usec": "{{.TIMESTAMP_USEC}}"
}`

	t := template.Must(template.New("payload").Parse(payloadTmpl))
	buf := &bytes.Buffer{}
	if err := t.Execute(buf, data); err != nil {
		panic(err)
	}
	payload := buf.String()

	return payload
}

func makeData(timeStampMicros int64) map[string]string {
	eventId, _ := uuid.NewV4()
	logBufferPartition, _ := uuid.NewV4()
	requestId, _ := uuid.NewV4()
	kubernetesEnv := getEnv("ENVIRONMENT", "test")

	data := map[string]string{
		"LOG_BUFFER_PARTITION":                logBufferPartition.String(),
		"EVENT_ID":            eventId.String(),
		"REQUEST_ID":            requestId.String(),
		"ENVIRONMENT":             kubernetesEnv,
		"TIMESTAMP_USEC": strconv.FormatInt(timeStampMicros, 10),
	}
	fmt.Println("Configuration:")
	for key, value := range data {
		fmt.Println("Key:", key, "Value:", value)
	}
	return data
}

func main() {
	// Go signal notification works by sending `os.Signal`
	// values on a channel. We'll create a channel to
	// receive these notifications (we'll also make one to
	// notify us when the program can exit).
	sigs := make(chan os.Signal, 2)
	done := make(chan bool, 1)

	interval, _ := strconv.Atoi(getEnv("INTERVAL", "30"))

	// Tickers use a similar mechanism to timers: a
	// channel that is sent values. Here we'll use the
	// `range` builtin on the channel to iterate over
	// the values as they arrive every 30s.
	ticker := time.NewTicker(time.Second * time.Duration(interval))

	var payload string

	go func() {
		for t := range ticker.C {
			// Nanosecond / 1,000 = microsecpmds
			fmt.Println("Tick at", t)
			payload = makePayload(makeData(t.UnixNano() / 1000))
			fmt.Println(payload)
		}
	}()
	// `signal.Notify` registers the given channel to
	// receive notifications of the specified signals.
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	// This goroutine executes a blocking receive for
	// signals. When it gets one it'll print it out
	// and then notify the program that it can finish.
	go func() {
		sig := <-sigs
		fmt.Println(sig, " received")

		// Tickers can be stopped like timers.
		ticker.Stop()
		fmt.Println("Ticker stopped")

		done <- true
	}()
	// The program will wait here until it gets the
	// expected signal (as indicated by the goroutine
	// above sending a value on `done`) and then exit.
	fmt.Println("awaiting signals")
	<-done
	fmt.Println("exiting")
	os.Exit(1)
}
