default: run

run:
	go install
	mrsjenkins
	docker run -d -p 9200:9200 --name mrsjenkins_elasticsearch elasticsearch

docker-image:
	docker build -t mrsjenkins .
	echo "Start by entering docker run -it --rm --name mrsjenkins mrsjenkins"

test:
	go test -v

test-elasticsearch-container:
	# YOU NEED DOCKER TO RUN ON 192.168.99.100
	eval $(docker-machine env)
	docker run -d -p 9200:9200 --name mrsjenkins_elasticsearch_test elasticsearch elasticsearch -Des.network.publish_host="192.168.99.100"

clean:
	rm -f $(GOPATH)/bin/mrsjenkins
	docker rm -f mrsjenkins_elasticsearch_test
