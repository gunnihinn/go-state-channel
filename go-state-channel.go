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

// Copy may be arbitrarily complicated if State contains slices, maps,
// pointers, or other structs.
func (s State) Copy() State {
	return s
}

func stateManager(stateCh chan State, toggle chan os.Signal) {
	state := State{}

	for {
		select {
		case stateCh <- state.Copy():

		case <-toggle:
			state.frobinate = !state.frobinate
		}
	}
}

func handler(stateCh chan State) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		s := <-stateCh

		start := s.frobinate

		time.Sleep(100 * time.Millisecond)

		if start != s.frobinate {
			http.Error(w, "Great non-success", http.StatusInternalServerError)
		} else {
			fmt.Fprintln(w, "Great success")
		}
	}
}

func help() string {
	return fmt.Sprintf(`Try me out by running:

$ curl localhost:8080
$ kill -s HUP %d
$ curl localhost:8080`, syscall.Getpid())
}

func main() {
	stateCh := make(chan State)

	toggle := make(chan os.Signal, 1)
	signal.Notify(toggle, syscall.SIGHUP)

	go stateManager(stateCh, toggle)

	http.HandleFunc("/", handler(stateCh))

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
