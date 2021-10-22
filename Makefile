PREFIX = /usr/local

pactop:
	go build pactop.go

all: pactop

install:
	install -D pactop $(PREFIX)/bin/pactop
	
clean:
	rm pactop
