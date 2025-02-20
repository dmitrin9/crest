package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
)

const (
	FLAGS_PLACEMENT   = "Flag or URL placement bad. Please insure URL is at the end of your command. All flags must be somewhere in between command-name (crest) and the argument (url)."
	PROTOCOL_REQUIRED = "All URL's must contain their protocol (http, https, etc...)"
	STATUS_ERROR      = "STATUS ERROR!"
)

func page(url string) error {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	res, err := client.Do(req)
	defer res.Body.Close()
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return errors.New(fmt.Sprintf("%s in %s | STATUS: %d", STATUS_ERROR, url, res.StatusCode))
	}

	return nil
}

func handle() error {
	args := os.Args

	/*
		crest <FLAGS> <URL>
		Flags:
			-t : test
			-v : verbose
	*/

	for i := range args {
		if args[i] == "-t" {
			if i == len(args)-1 {
				return errors.New(FLAGS_PLACEMENT)
			}

			url := args[len(args)-1]
			if url[0:7] != "http://" && url[0:8] != "https://" {
				fmt.Println(url[0:7])
				return errors.New(PROTOCOL_REQUIRED)
			}
			page(url)
		}
	}
	return nil
}

func main() {
	log.Fatal(handle())
}
