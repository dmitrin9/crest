package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

const (
	FLAGS_PLACEMENT   = "Flag or URL placement bad. Please insure URL is at the end of your command. All flags must be somewhere in between command-name (crest) and the argument (url)."
	PROTOCOL_REQUIRED = "All URL's must contain their protocol (http, https, etc...)"
	STATUS_ERROR      = "STATUS ERROR!"
)

func getLinks(n *html.Node) []string {

	var links []string
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				links = append(links, attr.Val)
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = append(links, getLinks(c)...)
	}
	return links
}

func page(url string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	//defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.New(fmt.Sprintf("%s in %s | STATUS: %d", STATUS_ERROR, url, res.StatusCode))
	}

	return res, nil
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
				return errors.New(PROTOCOL_REQUIRED)
			}

			base_res, err := page(url)
			if err != nil {
				return err
			}
			node, err := html.Parse(base_res.Body)
			if err != nil {
				return err
			}
			links := getLinks(node)
			for i := range links {
				r, err := page(url + links[i])
				if err != nil {
					fmt.Println(fmt.Sprintf("QUITTED ON LINK %d of %d TOTAL LINKS", i, len(links)))
					return err
				}
				r.Body.Close()
			}
			base_res.Body.Close()
			fmt.Println("ALL LINKS PASSED!")
		}
	}
	return nil
}

func main() {
	log.Fatal(handle())
}
