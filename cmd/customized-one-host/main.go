package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

const NumOfWorkers = 50
const NumOfRequests = 10000

var oneHostClient = &http.Client{
	Transport: createOneHostTransport(),
}

var host = flag.String("host", "0.0.0.0", "host of server")

func main() {
	flag.Parse()

	inputs := make(chan bool, 100)

	start := time.Now()

	wg := sync.WaitGroup{}

	for i := 0; i < NumOfWorkers; i++ {
		wg.Add(1)
		go func() {
			for range inputs {
				err := makeRequest()
				if err != nil {
					log.Fatal(err)
				}
			}
			wg.Done()
		}()
	}

	for i := 0; i < NumOfRequests; i++ {
		inputs <- true
	}
	close(inputs)

	wg.Wait()

	elapsed := time.Since(start)

	fmt.Println("elapsed", elapsed)
}

func makeRequest() error {
	resp, err := http.Get("http://" + *host + ":8080")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func createOneHostTransport() *http.Transport {
	result := http.DefaultTransport.(*http.Transport).Clone()

	result.MaxIdleConnsPerHost = result.MaxIdleConns

	return result
}
