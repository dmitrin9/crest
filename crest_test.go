package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"testing"
)

const workdir = "/home/dmitri/repos/crawl-tester/"

func printStream(stdout string, stderr string) {
	fmt.Fprintln(os.Stdout, "\033[32m"+stdout+"\033[0m")
	fmt.Fprintln(os.Stderr, "\033[31m"+stderr+"\033[0m")
}

func command(args []string, ctx *Context) {
	err := Handle(args, ctx)
	if err != nil {
		log.Fatal(err)
	}
}

func printTheThing(arg []LexNode) {
	for i := range arg {
		if arg[i].tok_raw != "" {
			fmt.Print(arg[i].tok_raw, ": ", arg[i])
			fmt.Println()
		}
	}
	for i := range arg {
		if arg[i].tok_raw != "" {
			fmt.Print(arg[i].tok_raw)
			fmt.Println()
		}
	}
}

func handleHtml(route string, path string) {
	http.HandleFunc(route, func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("Template Initialized for " + path)
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("Template Executed for " + path)

	})
}

func HttpTestsiteRun(wg *sync.WaitGroup) http.Server {
	handleHtml("/", "test_environment/index.html")
	handleHtml("/ActuallyExists", "test_environment/ActuallyExists.html")
	handleHtml("/AnotherWorkingSite", "test_environment/AnotherWorkingSite.html")
	handleHtml("/DoesNotExist", "test_environment/DoesNotExist.html")
	http.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		content, err := ioutil.ReadFile("test_environment/robots.txt")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Write(content)
	})
	srv := http.Server{Addr: ":8080"}

	wg.Add(1)
	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	wg.Done()

	return srv
}

// Initialize testing server.
func TestInit(t *testing.T) {
	var wg sync.WaitGroup
	srv := HttpTestsiteRun(&wg)
	defer srv.Close()
	wg.Wait()
}

// Testing minimal flags.
// Since an error when crawling some paths is expected, it will not throw.
// An error will only throw if the paths that are supposed to fail are
// crawled successfully if that makes sense.
func TestHttpCrawlingMinimalFlags(t *testing.T) {
	var ctx Context
	var err error
	var args []string

	args = []string{"crest", "-tv", "http://localhost:8080"}
	err = Handle(args, &ctx)
	if err == nil {
		t.Fatalf("%v", err)
	}

}

// Happy testing to test maximum flags.
func TestHttpCrawlingAllFlags(t *testing.T) {
	var ctx Context
	var err error
	var args []string

	args = []string{"crest", "-tfv", "http://localhost:8080"}
	err = Handle(args, &ctx)
	if err != nil {
		t.Fatalf("%v", err)
	}
}

// Test crestfile parser/ir-compiler.
func TestParse(t *testing.T) {
	var s State
	var err error

	raw, err := os.ReadFile("test_environment/Crestfile")
	if err != nil {
		t.Fatalf("%v", err)
	}
	data := string(raw)
	s.raw = data

	fmt.Println("TEST PRINT LEXER: ")
	err = s.Lexer()
	if err != nil {
		t.Fatalf("%v", err)
	}
	printTheThing(s.lexNodes)

	fmt.Println("test print parser: ")
	err = s.Parser()
	if err != nil {
		t.Fatalf("%v", err)
	}
	fmt.Println(s.parserNodes)

	fmt.Println("TEST PRINT COMPILER: ")
	err = s.Compiler()
	if err != nil {
		t.Fatalf("%v", err)
	}
	fmt.Println("VARIABLES: ", s.variable)
	fmt.Println("INSTRUCTION SET: ")
	for i := 0; i < len(s.instructionSet)-2; i += 2 {
		fmt.Println(s.instructionSet[i], " ", s.instructionSet[i+1])
	}

}

// Test crestfile handler.
func TestCrestfileHandler(t *testing.T) {
	args := []string{"crest", "run", "test_environment/Crestfile"}

	ctx := Context{
		verbose:      false,
		quiet:        false,
		followRobots: false,

		exclude: []string{},
	}
	s := State{
		raw:         "",
		lexNodes:    []LexNode{},
		parserNodes: []ParserNode{},

		variable:       make(map[string]string),
		instructionSet: []string{},

		offset: 0,
		row:    0,
		col:    0,
	}

	if err := HandleFile(args, &s, &ctx); err != nil {
		t.Fatalf("%v", err)
	}
}
