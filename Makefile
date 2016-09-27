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

test-coverage:
	go test -cover

test-coverage-profile:
	go test -coverprofile cover.out
	go tool cover -html=cover.out -o test-coverage.html
	rm -f cover.out

test-elasticsearch-container:
	# YOU NEED DOCKER TO RUN ON 192.168.99.100
	eval $(docker-machine env)
	docker run -d -p 9200:9200 --name mrsjenkins_elasticsearch_test elasticsearch elasticsearch -Des.network.publish_host="192.168.99.100"

clean:
	rm -f $(GOPATH)/bin/mrsjenkins
	docker rm -f mrsjenkins_elasticsearch_test

help:
	# make run - starts a docker container with elastic search and runs mrsjenkins locally
	# make docker-image - builds a docker image (mrsjenkins)
	# make test - runs the tests
	# make test-coverage - shows the test coverage of this package
	# make test-coverage-profile - generates a test coverage profile under test-coverage.html
	# make test-elasticsearch-container - runs an elasticsearch docker container for the integration tests
	# make clean - cleans up the test container and the binary
