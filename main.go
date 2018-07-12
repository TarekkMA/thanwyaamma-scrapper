package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/TarekkMA/thanwyaamma-scrapper/score"

	"github.com/TarekkMA/thanwyaamma-scrapper/score/score2018"
	flags "github.com/jessevdk/go-flags"
	log "github.com/sirupsen/logrus"
)

func main() {

	var opts struct {
		SeatingStart int `short:"s" description:"The startpoint" required:"true"`
		Concurrency  int `short:"c" description:"Number of concurrent workers" required:"true"`
		Number       int `short:"n" description:"Number of seats to scrap" required:"true"`
	}

	flags.Parse(&opts)

	if opts.Number < opts.Concurrency {
		opts.Concurrency = opts.Number
	}

	results := make(chan *score.Result)
	input := make(chan int)
	errs := make(chan *score.Error)

	s := score2018.NewScrepper()

	var wg sync.WaitGroup
	var handlersWg sync.WaitGroup

	wg.Add(opts.Concurrency)

	go func() {
		for i := opts.SeatingStart; i <= opts.SeatingStart+opts.Number; i++ {
			input <- i
		}
		close(input)
	}()

	for i := 0; i < opts.Concurrency; i++ {
		go func() {
			defer wg.Done()
			for i := range input {
				if res, err := s.Get(int32(i)); err != nil {
					errs <- err
				} else if res != nil {
					results <- res
				}
			}
		}()
	}

	handlersWg.Add(2)
	go handelErr(errs, &handlersWg)
	go handelRes(results, &handlersWg)

	wg.Wait()
	close(results)
	close(errs)
	handlersWg.Wait()
}

func handelErr(errs chan *score.Error, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.OpenFile("errs.json", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Error(err)
	}

	for e := range errs {
		log.Errorf("ERR ==== %+v", e)

		if _, err := file.WriteString(fmt.Sprintf("%d : %v \n", e.SeatingNumber, e.Err)); err != nil {
			log.Fatal(err)
		}
	}
}

func handelRes(results chan *score.Result, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.OpenFile("results.json", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if _, err := file.WriteString("[\n"); err != nil {
		log.Fatal(err)
	}

	defer func() {
		if _, err := file.WriteString("]"); err != nil {
			log.Fatal(err)
		}
	}()

	for r := range results {
		if bts, jsonErr := json.Marshal(r); err != nil {
			log.Error(jsonErr)
		} else if _, err := file.Write(append(bts, []byte(",\n")...)); err != nil {
			log.Error(err)
		}
	}

}
