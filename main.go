package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	addr := ":" + os.Getenv("PORT")
	http.HandleFunc("/", handle)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func handle(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Error parsing form.", http.StatusBadRequest)
		return
	}

	resp, err := http.Get("https://api.coindesk.com/v1/bpi/currentprice.json")
	if err != nil {
		http.Error(w, "Error requesting url", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading response body", http.StatusInternalServerError)
		return
	}

	type Price struct {
		Rate string `json:"rate"`
		Code string `json:"code"`
	}

	type Message struct {
		Prices struct {
			USD Price `json:"USD"`
			GBP Price `json:"GBP"`
			EUR Price `json:"EUR"`
		} `json:"bpi"`
	}

	var f Message
	json.Unmarshal(body, &f)

	ps := f.Prices
	fmt.Fprintf(w, "%v@%v\n%v@%v\n%v@%v", ps.GBP.Code, ps.GBP.Rate, ps.USD.Code, ps.USD.Rate, ps.EUR.Code, ps.EUR.Rate)
}
