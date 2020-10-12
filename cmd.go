package main

import (
	"fmt"
	"log"
	"os"
	
	"time"

	"net/http"
	_ "net/http/pprof"

	
	"github.com/spf13/viper"
)

func main() {

	err := rootCmd.Execute()
	if err != nil {
		panic(err)
	}
	if viper.GetBool("profile") {
		go func() {
			log.Printf("Starting Profiler localhost:%d", viper.GetInt32("profile-port"))
			log.Println(http.ListenAndServe(fmt.Sprintf("localhost:%d", viper.GetInt32("profile-port")), nil))
		}()
	}
	errs := validateFlags()
	if len(errs) != 0 {
		for _, e := range errs {
			log.Print(e)
		}
		os.Exit(1)
	}
	dur, err := time.ParseDuration("5s")
	if err != nil {
		panic(err)
	}
	timer := time.NewTimer(dur)
	log.Print("Start timer")
	var errChan chan error
	for {
		select {
		case <-timer.C:
			log.Printf("update call")
			errChan = update()
		case err := <-errChan:
			if err != nil {
				log.Print(err)
			}
			log.Print("Restarting Timer")
			timer.Reset(dur)
			errChan = nil
		}

	}

}

func update() chan error {
	errChan := make(chan error)
	go func() {
		defer close(errChan)
		log.Printf("done update")
		errChan <- nil
	}()
	log.Printf("errChan: %v", errChan)
	return errChan
}

