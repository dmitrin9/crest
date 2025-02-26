package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/net/html"
)

const (
	COMMAND_NOT_RECOGNIZED = "Command was not recognized. Check if you've entered a valid command."
	FLAGS_PLACEMENT        = "Flag or URL placement bad. Please insure URL is at the end of your command. All flags must be somewhere in between command-name (crest) and the argument (url)."
	NON_LOCALHOST_CRAWL    = "You are trying to crawl a site that is not on your localhost. This action is forbidden. \nDon't fret! If your site has hrefs which redirect to other sites, they will be ignored and won't throw errors. However, crawling an entirely different domain is entirely unsupported."
	SCHEME_REQUIRED        = "All URL's must contain their scheme (http, https, etc...)"
	STATUS_ERROR           = "STATUS ERROR!"
)

func splitUrl(raw string) map[string]string {
	parsed, err := url.Parse(raw)
	if err != nil {
		log.Fatal(err)
	}

	urlStructure := make(map[string]string)

	urlStructure["scheme"] = parsed.Scheme
	urlStructure["hostname"] = parsed.Hostname()
	urlStructure["path"] = parsed.Path

	return urlStructure
}

func getLinks(n *html.Node) []string {

	var links []string
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, attr := range n.Attr {
			if attr.Key == "href" {
				link := attr.Val
				/*
					if url[0:7] != "http://" && url[0:8] != "https://" {
						return errors.New(SCHEME_REQUIRED)
					}
				*/
				urlStructure := splitUrl(link)
				if urlStructure["hostname"] == "" {
					links = append(links, link)
				}
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
			-t/--test-http : test
			-v/--verbose : verbose
	*/

	if len(args) == 0 {
		return errors.New(COMMAND_NOT_RECOGNIZED)
	}

	verbose := false

	for i := range args {
		if args[i] == "-v" || args[i] == "--verbose" {
			if i == len(args)-1 {
				return errors.New(FLAGS_PLACEMENT)
			}
			verbose = true
		}
	}

	for i := range args {
		if args[i] == "-t" || args[i] == "--test-http" {
			if i == len(args)-1 {
				return errors.New(FLAGS_PLACEMENT)
			}

			url := args[len(args)-1]
			urlData := splitUrl(url)
			if urlData["scheme"] != "http" && urlData["scheme"] != "https" {
				return errors.New(SCHEME_REQUIRED)
			}
			if urlData["hostname"] != "localhost" && urlData["hostname"] != "127.0.0.1" {
				return errors.New(NON_LOCALHOST_CRAWL)

			}
			if verbose {
				fmt.Println("OK: PROTOCOL DEFINED")
			}

			base_res, err := page(url)
			if err != nil {
				return err
			}
			if verbose {
				fmt.Println("OK: BASE PAGE RESPONSE")
			}
			node, err := html.Parse(base_res.Body)
			if err != nil {
				return err
			}
			if verbose {
				fmt.Println("OK: HTML NODES FOR BASE RESPONSE")
			}
			links := getLinks(node)
			for i := range links {
				r, err := page(url + links[i])
				if err != nil {
					fmt.Println(fmt.Sprintf("ERROR: QUITTED ON LINK %d of %d TOTAL LINKS", i, len(links)))
					return err
				}
				if verbose {
					fmt.Println(fmt.Sprintf("OK: PAGE RESPONSE FOR NODE %d LINK %s", i+1, links[i]))
				}
				r.Body.Close()
				if verbose {
					fmt.Println("OK: CLOSED LINK RESPONSE")
				}
			}
			base_res.Body.Close()
			if verbose {
				fmt.Println("OK: CLOSED BASE RESPONSE")
			}
			fmt.Println("ALL CHECKS PASSED")
		}
	}
	return nil
}

func main() {
	log.Fatal(handle())
}
