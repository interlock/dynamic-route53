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
	log.Printf("Frequency: %v", viper.GetInt64("frequency"))
	var dur time.Duration = time.Duration(viper.GetInt64("frequency")) * time.Second

	timer := time.NewTimer(dur)
	log.Printf("Start timer (%d)", dur)
	var errChan chan error
	var exit bool = false
	for {
		select {
		case t := <-timer.C:
			log.Printf("update call: %v", t)
			errChan = update()
		case err := <-errChan:
			if err != nil {
				log.Printf("Error: %s", err.Error())
			}
			errChan = nil
			if dur == 0 {
				timer.Stop()
				log.Printf("Finished single run")
				exit = true
				break
			}
			log.Print("Restarting Timer")
			timer.Reset(dur)
		}
		if exit == true {
			break
		}
	}
}

func update() chan error {
	errChan := make(chan error)
	go func() {
		defer close(errChan)
		err := doUpdate()
		log.Printf("done update")
		errChan <- err
	}()

	return errChan
}
