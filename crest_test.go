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

func printHooks(hooks map[string]string) {
	for key, value := range hooks {
		fmt.Println(key, " ", value)
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

func TestCrestfileHandler(t *testing.T) {
	args := []string{"crest", "run", "test_environment/Crestfile"}

	ctx := Context{
		verbose:      false,
		quiet:        false,
		followRobots: false,
		hooks:        make(map[string]string),

		exclude: []string{},
	}
	s := State{
		raw:         "",
		lexNodes:    []LexNode{},
		parserNodes: []ParserNode{},

		variable:       make(map[string]string),
		hooks:          make(map[string]string),
		instructionSet: []string{},

		offset: 0,
		row:    0,
		col:    0,
	}

	if err := HandleFile(args, &s, &ctx); err != nil {
		t.Fatalf("%v", err)
	}
	fmt.Println("Hooks")
	printHooks(ctx.hooks)
}

func TestUtils(t *testing.T) {
	ctx := Context{
		verbose:      false,
		quiet:        false,
		followRobots: false,
		hooks:        make(map[string]string),

		exclude: []string{},
	}
	err := ctx.computeHook("echo Hello")
	if err != nil {
		t.Fatalf("%v", err)
	}
}
