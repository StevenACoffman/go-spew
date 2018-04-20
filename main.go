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
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time" // or "runtime"

	"github.com/satori/go.uuid"
	"strconv"
	"encoding/json"
)

var isDebug bool

func init() {
	isDebug,_ = strconv.ParseBool(os.Getenv("DEBUG"))
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}


type IthakaStructuredLogMessage struct {
	EventType  string   `json:"eventtype,omitempty"`
	EventId    string   `json:"eventid,omitempty"`
	Origin     string   `json:"origin,omitempty"`
	RequestId  string   `json:"requestid,omitempty"`
	//Sometimes an int64, sometimes a string, json.Number is either
	TstampUsec int64    `json:"tstamp_usec,omitempty"`
	Dests      []string `json:"dests,omitempty"`
	NodeName   string   `json:"node_name,omitempty"`
}

func makeMessage(environment string, nodeName string, timestampMicros int64) string {
	var payload IthakaStructuredLogMessage

	eventId, _ := uuid.NewV4()
	requestId, _ := uuid.NewV4()


	payload.Dests = []string{ "k8s-fluent-bit-watermark"}
	payload.EventId = eventId.String()
	payload.RequestId = requestId.String()
	payload.Origin = "watermark."+environment
	payload.EventType = "watermark"
	payload.TstampUsec = timestampMicros
	payload.NodeName  = nodeName
	str, _ := json.Marshal(payload)
	return string(str)
}

func main() {
	// Go signal notification works by sending `os.Signal`
	// values on a channel. We'll create a channel to
	// receive these notifications (we'll also make one to
	// notify us when the program can exit).
	sigs := make(chan os.Signal, 2)
	done := make(chan bool, 1)

	kubernetesEnv := getEnv("ENVIRONMENT", "test")
	nodeNameEnv := getEnv("NODE_NAME", "unknown")

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
			if isDebug {
				fmt.Println("Tick at", t)
			}
			payload = makeMessage(kubernetesEnv, nodeNameEnv, t.UnixNano() / 1000)
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
