PREFIX = ./bin/crest # Before building, edit your prefix to suite your needs. It is by default set to the repo's bin dir for the sake of testing.

build:
	go build -o $(PREFIX) *.go
