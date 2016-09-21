# Make

## run [default]

builds and runs mrsjenkins

## test

Runs tests. Needs a running elasticsearch on elasticsearch.mrsjenkins.test (see test-elasticsearch-container)

## test-elasticsearch-container

Starts an elasticsearch container for test purposes.

Requirements:

* docker-machine default running
* a host entry with elasticsearch.mrsjenkins.test to the docker-machine

## clean

Deletes the mrsjenkins binary and the elasticsearch test container.
