package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"

	"golang.org/x/net/html"
)

const (
	INVALID_AMOUNT_COMMANDLINE_ARGUMENTS = "Invalid amount of commandline arguments."
	FLAGS_PLACEMENT                      = "Flag or URL placement bad. Please insure URL is at the end of your command. All flags must be somewhere in between command-name (crest) and the argument (url)."
	NON_LOCALHOST_CRAWL                  = "You are trying to crawl a site that is not on your localhost. This action is forbidden. \nDon't fret! If your site has hrefs which redirect to other sites, they will be ignored and won't throw errors. However, crawling an entirely different domain is entirely unsupported."
	SCHEME_REQUIRED                      = "All URL's must contain their scheme (http, https, etc...)"
	STATUS_ERROR                         = "STATUS ERROR!"
	INVALID_TEST                         = "Testing type is either invalid or unspecified: specify with '--test-http'/'-t' flag"
	UNRECOGNIZED_COMMAND                 = "Command unrecognized. Please look at the documentation. If you believe there's a problem with crest, feel free to create an issue. Just make sure to read the readme.md file and the issues tab first to see if your issue is already being worked on."
	INCLUDE_PORT                         = "As of now, your URL must include a port."
)

type Context struct {
	quiet        bool
	verbose      bool
	followRobots bool
	exclude      []string
	hooks        map[string]string

	CURRENT string
	CONTENT string
}

type RobotPolicy struct {
	agent    string
	allow    []string
	disallow []string
}

func (p *RobotPolicy) Empty() bool {
	if len(p.allow) == 0 && len(p.disallow) == 0 && len(p.agent) == 0 {
		return true
	}
	return false
}

func (c *Context) computeExcludedLinks(links []string) []string {
	tmp := []string{}
	n := len(c.exclude)
	if n == 0 {
		tmp = append(tmp, links...)
	} else {
		for _, exclude := range c.exclude {
			for _, link := range links {
				if link != exclude {
					tmp = append(tmp, link)
				}
			}
		}
	}
	return tmp
}

func (c *Context) computeHook(hookCode string) error {
	hook := strings.Split(hookCode, " ")
	command := hook[0]
	args := hook[1:]
	cmd := exec.Command(command, args...)
	out, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			fmt.Printf("Exited with code: %d\n", exitErr.ExitCode())
			fmt.Printf("Error output: %s\n", string(exitErr.Stderr))
		} else {
			fmt.Printf("Error: %v\n", err)
		}
	}

	c.printv(os.Stdout, string(out), "")

	return nil
}

func (c *Context) printv(stream io.Writer, out string, longOut string) {
	reset := "\033[0m"
	//debugColor := "\033[93m DEBUG: "
	color := "\033[31m ERROR: "
	if stream == os.Stdout {
		color = "\033[32m OK: "
	}

	if c.verbose {
		if len(longOut) == 0 {
			fmt.Fprintln(stream, color+out+reset) // If verbose arg is empty, print normal text in place of verbose text.
		} else {
			fmt.Fprintln(stream, color+longOut+reset)
		}
	} else if c.quiet {
		if stream == os.Stderr {
			fmt.Fprintln(stream, color+out+reset)
		}
	} else {
		fmt.Fprintln(stream, color+out+reset)
	}
}

func splitUrl(raw string) map[string]string {
	urlStructure := make(map[string]string)

	parsed, err := url.Parse(raw)
	if err != nil {
		log.Fatal(err)
	}
	host := parsed.Host

	urlStructure["port"] = ""
	for i := range host {
		if host[i] == ':' {
			splitted := strings.Split(host, ":")
			urlStructure["port"] = splitted[1]
		}
	}

	urlStructure["scheme"] = parsed.Scheme
	urlStructure["hostname"] = parsed.Hostname()
	urlStructure["path"] = parsed.Path
	urlStructure["fragment"] = parsed.Fragment

	return urlStructure
}

func RobotParser(url string, ctx *Context) ([]RobotPolicy, error) {
	var robotPolicies []RobotPolicy
	robot_path := url + "/robots.txt"

	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, robot_path, nil)
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("status err")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	robotContents := strings.Split(string(body), "\n")

	var policy RobotPolicy
	for i := range robotContents {

		if len(robotContents[i]) > 0 {
			comment := string(robotContents[i][0])
			if comment == "#" || comment == " " {
				continue
			}
		}

		if len(robotContents[i]) >= 8 {
			disallow := robotContents[i][0:8]
			if disallow == "Disallow" {
				content := strings.TrimSpace(robotContents[i][9:])
				policy.disallow = append(policy.disallow, content)
			}
		}

		if len(robotContents[i]) >= 5 {
			allow := robotContents[i][0:5]
			if allow == "Allow" {
				content := strings.TrimSpace(robotContents[i][6:])
				policy.allow = append(policy.allow, content)
			}
		}

		if len(robotContents[i]) >= 10 {
			agent := robotContents[i][0:10]
			if agent == "User-agent" {
				if policy.Empty() {
					content := strings.TrimSpace(robotContents[i][11:])
					policy.agent = content
				} else {
					robotPolicies = append(robotPolicies, policy)
					policy = RobotPolicy{}
					content := strings.TrimSpace(robotContents[i][11:])
					policy.agent = content
				}
			}
		}
		if i == len(robotContents)-1 {
			robotPolicies = append(robotPolicies, policy)
		}

	}
	return robotPolicies, nil
}

func GetAllowedRobots(url string, links []string, ctx *Context) ([]string, error) {
	/*
	 * Compare the policy provided by robots.txt file
	 * in order to determine which paths to permit.
	 * The final list of allowed paths is called the
	 * delta.
	 */
	policies, err := RobotParser(url, ctx)
	ctx.printv(os.Stdout, "Generated robot policies", "")
	if err != nil {
		return nil, err
	}
	var delta []string

	disallows := []string{}
	for p := range policies {
		policy := policies[p]
		if policy.agent == "Crestbot" || policy.agent == "*" {
			disallows = append(disallows, policy.disallow...)
		}
	}

	for l := range links {
		linkDisallow := false
		link := links[l]
		for d := range disallows {
			if link == disallows[d] {
				linkDisallow = true
				break
			}
		}
		if !linkDisallow {
			delta = append(delta, link)
		}
	}

	return delta, nil
}

func Page(url string, path string, ctx *Context) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url+path, nil)
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		ctx.printv(os.Stderr, "Recieved status error", "")
		return nil, errors.New(fmt.Sprintf("%s in %s | STATUS: %d", STATUS_ERROR, url, res.StatusCode))
	}

	if path == "" {
		path = "/"
	}
	hook := ctx.hooks[path]
	if len(ctx.hooks) > 0 && len(hook) > 0 {
		err := ctx.computeHook(hook)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func getPageLinksTask(n *html.Node) []string {
	var links []string

	for c := range n.Descendants() {
		var linkBuffer []string
		if c.Type == html.ElementNode && c.Data == "a" {
			for _, attr := range c.Attr {
				if attr.Key == "href" {
					link := attr.Val
					urlStructure := splitUrl(link)
					if urlStructure["hostname"] == "" && urlStructure["scheme"] == "" && urlStructure["fragment"] == "" {
						linkBuffer = append(linkBuffer, link)
					}
				}
			}
		}
		links = append(links, linkBuffer...)
	}
	return links
}

func RecursiveLinkCheck(url string, links []string, ctx *Context, depth int) error {
	/*
	 * Crawl non-fragment URLs declared in anchor tags
	 * to a depth not exceeding 20 stack frames. This
	 * is the meat and potatoes of the --test-http flag
	 * and by extension the entirety of crest.
	 */
	if depth == 0 {
		path := splitUrl(url)["path"]
		res, err := Page(url, path, ctx)
		if err != nil {
			return err
		}
		ctx.printv(os.Stdout, "Created base response for "+url, "")
		node, err := html.Parse(res.Body)
		if err != nil {
			return err
		}
		ctx.printv(os.Stdout, "Got HTML nodes", "")

		links = getPageLinksTask(node)
		tmp := ctx.computeExcludedLinks(links)
		links = tmp

		if ctx.followRobots {
			tmp, err := GetAllowedRobots(url, links, ctx)
			if err != nil {
				return err
			}
			links = tmp
		}
		ctx.printv(os.Stdout, "Closing base page response closure", "Closing base page response closure for "+url)
		res.Body.Close()
	}

	newLinks := []string{}
	for i := range links {

		r, err := Page(url, links[i], ctx)
		if err != nil {
			ctx.printv(os.Stderr, fmt.Sprintf("Quitted at %s which is link %d of %d total links at link recursion depth %d", links[i], i, len(links), depth), "")
			return err
		}
		ctx.printv(os.Stdout, "Response open", fmt.Sprintf("Response opened at depth %d", depth))
		node, err := html.Parse(r.Body)
		if err != nil {
			ctx.printv(os.Stderr, "Problem getting nodes", "Problem getting HTML nodes from request")
			return err
		}
		if ctx.followRobots {
			tmp, err := GetAllowedRobots(url, ctx.computeExcludedLinks(getPageLinksTask(node)), ctx)
			if err != nil {
				return err
			}
			newLinks = append(newLinks, tmp...)
		} else {
			newLinks = append(newLinks, ctx.computeExcludedLinks(getPageLinksTask(node))...)
		}
		r.Body.Close()
		ctx.printv(os.Stdout, "Response closed", fmt.Sprintf("Response closed at depth %d", depth))
	}

	if depth > 20 {
		return nil
	}

	return RecursiveLinkCheck(url, newLinks, ctx, depth+1)

}

func Handle(args []string, ctx *Context) error {
	/*
	 * Handle commandline stuff.
	 * The code is very self explanitory.
	 * Will probably refactor at some point
	 * if crest's commandline interface becomes
	 * expressive enough.
	 */
	var test string

	if len(args) <= 1 {
		return errors.New(UNRECOGNIZED_COMMAND)
	}

	for i := range args {
		if string(args[i][0]) == "-" {
			if i == len(args)-1 {
				return errors.New(FLAGS_PLACEMENT)
			}
			for j := range args[i] {
				c := string(args[i][j])
				if c == "v" {
					ctx.verbose = true
				}
				if c == "q" {
					ctx.quiet = true
				}
				if c == "f" {
					ctx.followRobots = true
				}
				if c == "t" {
					test = "test-http"
				}
			}
		}
		if args[i] == "--verbose" {
			if i == len(args)-1 {
				return errors.New(FLAGS_PLACEMENT)
			}
			ctx.verbose = true
		}
		if args[i] == "--quiet" {
			if i == len(args)-1 {
				return errors.New(FLAGS_PLACEMENT)
			}
			ctx.quiet = true
		}
		if args[i] == "--follow-robots" {
			if i == len(args)-1 {
				return errors.New(FLAGS_PLACEMENT)
			}
			ctx.followRobots = true
		}
		if args[i] == "--test-http" {
			if i == len(args)-1 {
				return errors.New(FLAGS_PLACEMENT)
			}
			test = "test-http"
		}
	}

	if test == "test-http" {
		url := args[len(args)-1]
		urlData := splitUrl(url)
		if urlData["scheme"] != "http" && urlData["scheme"] != "https" {
			return errors.New(SCHEME_REQUIRED)
		}
		if urlData["hostname"] != "localhost" && urlData["hostname"] != "127.0.0.1" {
			return errors.New(NON_LOCALHOST_CRAWL)
		}
		if len(urlData["port"]) == 0 {
			return errors.New(INCLUDE_PORT)
		}
		if err := RecursiveLinkCheck(url, []string{}, ctx, 0); err != nil {
			return err
		}
		ctx.printv(os.Stdout, "Got links", "Recursive link check done")
	} else {
		return errors.New(INVALID_TEST)
	}
	return nil
}

func HandleFile(args []string, s *State, ctx *Context) error {
	if len(args) != 3 {
		return errors.New(INVALID_AMOUNT_COMMANDLINE_ARGUMENTS)
	}
	filename := args[2]
	raw, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	data := string(raw)
	s.raw = data

	err = s.Lexer()
	if err != nil {
		return err
	}

	err = s.Parser()
	if err != nil {
		return err
	}

	err = s.Compiler()
	if err != nil {
		return err
	}

	var test string
	var url string
	instructions := s.instructionSet

	// [type testHTTP verbose true followRobots true exclude hello exclude hello/world some-unrelated-tool ]
	for i := 0; i < len(instructions); i += 2 {
		current := instructions[i]
		next := instructions[i+1]

		if current == "type" {
			test = next
		} else if current == "verbose" {
			if next == "true" {
				ctx.verbose = true
			}
			if next == "false" {
				ctx.verbose = false
			}
		} else if current == "followRobots" {
			if next == "true" {
				ctx.followRobots = true
			}
			if next == "false" {
				ctx.followRobots = false
			}
		} else if current == "exclude" {
			ctx.exclude = append(ctx.exclude, next)
		} else if current == "quiet" {
			if next == "true" {
				ctx.quiet = true
			}
			if next == "false" {
				ctx.quiet = false
			}
		} else if current == "url" {
			url = next
		}
		ctx.hooks = s.hooks
	}
	ctx.printv(os.Stdout, "Successfully compiled crestfile instruction set", "")

	if test == "testHTTP" {
		urlData := splitUrl(url)
		if urlData["scheme"] != "http" && urlData["scheme"] != "https" {
			return errors.New(SCHEME_REQUIRED)
		}
		if urlData["hostname"] != "localhost" && urlData["hostname"] != "127.0.0.1" {
			return errors.New(NON_LOCALHOST_CRAWL)
		}
		if len(urlData["port"]) == 0 {
			return errors.New(INCLUDE_PORT)
		}
		if err := RecursiveLinkCheck(url, []string{}, ctx, 0); err != nil {
			return err
		}
		ctx.printv(os.Stdout, "Got links", "Recursive link check done")
	} else {
		return errors.New(INVALID_TEST)
	}

	return nil
}
