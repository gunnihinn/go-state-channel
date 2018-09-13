// This example only works on UNIX derivatives, as it uses the SIGHUP signal.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type State struct {
	frobinate bool
}

var state State

func handler(w http.ResponseWriter, r *http.Request) {
	start := state.frobinate

	time.Sleep(100 * time.Millisecond)

	if start != state.frobinate {
		http.Error(w, "Great non-success", http.StatusInternalServerError)
	} else {
		fmt.Fprintln(w, "Great success")
	}
}

func main() {
	toggle := make(chan os.Signal, 1)
	signal.Notify(toggle, syscall.SIGHUP)

	go func() {
		for {
			<-toggle
			state.frobinate = !state.frobinate
		}
	}()

	http.HandleFunc("/", handler)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
