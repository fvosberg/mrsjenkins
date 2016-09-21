default: run

run:
	go install
	mrsjenkins

test:
	go test -v

elasticsearch-container:
	eval $(docker-machine env)
	docker run -d -p 9200:9200 --name mrsjenkins_elasticsearch elasticsearch

clean:
	rm $(GOPATH)/bin/mrsjenkins
	docker rm -f mrsjenkins_elasticsearch
