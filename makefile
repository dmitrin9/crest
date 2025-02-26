INSTALL_PATH=~/local/bin/ # change to the directory where you want to install to.
DEST_PATH=./bin/crest

all: build

build:
	go build -o $(DEST_PATH) *.go
	@echo "Successfully built crest"

install: 
	install -m 755 $(DEST_PATH) $(INSTALL_PATH)	
	@echo "Installed crest to your install path"

clean:
	rm -f ./bin/*
	@echo "Cleaning up..."
