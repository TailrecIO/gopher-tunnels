package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

func main() {

	var port int
	flag.IntVar(&port, "port", 9000, "listening port")
	flag.Parse()

	var counter uint64

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		defer r.Body.Close()

		atomic.AddUint64(&counter, 1)

		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}
		body := string(bodyBytes)

		counterFinal := atomic.LoadUint64(&counter)

		fmt.Printf("[Request: %v ] ========================================\n", counterFinal)
		fmt.Printf("Received at: %v\n", time.Now().String())
		fmt.Printf("Path: %v\n", r.URL.String())
		fmt.Printf("Method: %v\n", r.Method)
		fmt.Println("Headers:")
		for k, v := range r.Header {
			fmt.Printf("\t- %v: %v\n", k, v)
		}
		fmt.Println("Body:")
		fmt.Println(body)
		fmt.Print("\n\n")
		w.Write(bodyBytes)
	})

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}