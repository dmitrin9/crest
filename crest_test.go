package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"testing"
	"text/template"
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
	fmt.Println("--------------------------")
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
		log.Println("Template Initialized.")
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		log.Println("Template Executed.")

	})
}

func HttpTestsiteRun(wg *sync.WaitGroup) http.Server {
	handleHtml("/", "views/index.html")
	handleHtml("/ActuallyExists", "views/ActuallyExists.html")
	handleHtml("/AnotherWorkingSite", "views/AnotherWorkingSite.html")
	handleHtml("/DoesNotExist", "views/DoesNotExist.html")
	http.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		content, err := ioutil.ReadFile("views/robots.txt")
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

// Just some basic happy tests, which works for now. I will add more comprehensive tests in the future.
func TestHttpCrawlingFlags(t *testing.T) {
	var wg sync.WaitGroup
	var ctx Context
	var err error

	srv := HttpTestsiteRun(&wg)

	args := []string{"crest", "-tfv", "http://localhost:8080"}
	err = Handle(args, &ctx)
	if err != nil {
		t.Fatalf("%v", err)
	}

	args = []string{"crest", "-tv", "http://localhost:8080"}
	err = Handle(args, &ctx)
	if err != nil {
		t.Fatalf("%v", err)
	}

	srv.Close()
	wg.Wait()
}

func TestParse(t *testing.T) {
	var s State
	var err error

	raw, err := os.ReadFile("views/Crestfile")
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

	fmt.Println("-----------------------------------")

	fmt.Println("TEST PRINT PARSER: ")
	CURRENT_PATH := "/"
	HTML_CONTENT := "<h1>Hello</h1>"
	s.CURRENT_PATH = CURRENT_PATH
	s.HTML_CONTENT = HTML_CONTENT
	err = s.Parser()
	if err != nil {
		t.Fatalf("%v", err)
	}
	fmt.Println(s.parserNodes)

	fmt.Println("-----------------------------------")
	fmt.Println("TEST PRINT COMPILER: ")
	err = s.Compiler()
	if err != nil {
		t.Fatalf("%v", err)
	}
	fmt.Println("VARIABLES: ", s.variable)
	fmt.Println("INSTRUCTION SET: ")
	for i := 2; i < len(s.instructionSet); i += 2 {
		fmt.Println(s.instructionSet[i], " ", s.instructionSet[i-1])
	}

}
