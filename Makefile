default: run

run:
	go install
	mrsjenkins

test:
	go test -v

clean:
	rm $(GOPATH)/bin/mrsjenkins
