package main

import (
	"errors"
	"fmt"
	"io"
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

type Context struct {
	quiet   bool
	verbose bool
}

func (c *Context) printv(stream io.Writer, out string, longOut string) {
	if c.verbose {
		fmt.Fprintln(stream, longOut)
	} else if c.quiet {
		if stream == os.Stderr {
			fmt.Fprintln(stream, out)
			return
		}
		return
	} else {
		fmt.Fprintln(stream, out)
	}
}

func splitUrl(raw string) map[string]string {
	parsed, err := url.Parse(raw)
	if err != nil {
		log.Fatal(err)
	}

	urlStructure := make(map[string]string)

	urlStructure["scheme"] = parsed.Scheme
	urlStructure["hostname"] = parsed.Hostname()
	urlStructure["path"] = parsed.Path
	urlStructure["fragment"] = parsed.Fragment

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
				if urlStructure["hostname"] == "" && urlStructure["scheme"] == "" && urlStructure["fragment"] == "" {
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

func recursiveLinkCheck(url string, links []string, ctx Context, depth int) error {
	newLinks := []string{}

	for i := range links {
		r, err := page(url + links[i])
		if err != nil {
			ctx.printv(os.Stderr, fmt.Sprintf("ERROR: Quitted at %s", links[i]), fmt.Sprintf("ERROR: Quitted at %s which is link %d of %d total links at link recursion depth %d", links[i], i, len(links), depth))
			return err
		}
		ctx.printv(os.Stdout, "OK: Response open", fmt.Sprintf("OK: recursive response opened, depth %d", depth))
		node, err := html.Parse(r.Body)
		if err != nil {
			ctx.printv(os.Stderr, "ERROR: Error getting nodes.", "ERROR: Error getting nodes from recursive response.")
			return err
		}
		newLinks = append(newLinks, getLinks(node)...)
		r.Body.Close()
	}
	if depth >= 20 {
		return nil
	}
	return recursiveLinkCheck(url, newLinks, ctx, depth+1)

}

func handle(ctx Context) error {
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

	for i := range args {
		if args[i] == "-v" || args[i] == "--verbose" {
			if i == len(args)-1 {
				return errors.New(FLAGS_PLACEMENT)
			}
			ctx.verbose = true
		}
		if args[i] == "-q" || args[i] == "--quiet" {
			if i == len(args)-1 {
				return errors.New(FLAGS_PLACEMENT)
			}
			ctx.quiet = true
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

			ctx.printv(os.Stdout, "OK: Scheme and host verified", "OK: Scheme is verified to be http or https, and host is verified to be localhost or 127.0.0.1")

			base_res, err := page(url)
			if err != nil {
				return err
			}
			ctx.printv(os.Stdout, "OK: Base page response", "OK: Response of the base link "+url+" has been created.")

			node, err := html.Parse(base_res.Body)
			if err != nil {
				return err
			}
			ctx.printv(os.Stdout, "OK: HTML Nodes", "OK: Acquired HTML nodes from base page response.")

			links := getLinks(node)
			if err = recursiveLinkCheck(url, links, ctx, 0); err != nil {
				return err
			}

			ctx.printv(os.Stdout, "OK: Get links", "OK: Recursive link check done")
			base_res.Body.Close()

			ctx.printv(os.Stdout, "OK: Response closed", "OK: Base response closed")
			ctx.printv(os.Stdout, "OK: Complete", "OK: Complete")
		}
	}
	return nil
}

func main() {
	ctx := Context{
		verbose: false,
		quiet:   false,
	}
	if err := handle(ctx); err != nil {
		log.Fatal(err)
	}
}
