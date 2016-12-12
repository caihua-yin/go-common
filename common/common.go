// Package common provides various common methods for each program
package common

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"time"
)

// Startup does necessary steps in every program
func Startup() {
	// Initialize seed of random number
	rand.Seed(time.Now().UnixNano())

	// Limit the number of operating system threads to execute user-level Go code simultaneously
	// Make it same as system logical CPU number ("CPU(s)" value of lscpu)
	// = "Socket(s)" * "Core(s) per socket" * "Thread(s) per core"
	if os.Getenv("GOMAXPROCS") == "" {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}
	log.Printf("Using %d CPUs...", runtime.GOMAXPROCS(-1))
}

// ParseDuration panics if error is not nil
func ParseDuration(s string) time.Duration {
	duration, err := time.ParseDuration(s)
	if err != nil {
		panic(fmt.Errorf("Unable to parse duration %#v: %s", s, err))
	}
	return duration
}
