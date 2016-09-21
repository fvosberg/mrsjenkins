# Make

After getting new packages or updating dependencies we have to update the Godeps/Godeps.json via godep save from the root folder

## run [default]

builds and runs mrsjenkins locally

## docker-image

Builds a docker image with all needed code for building inside. Maybe we will change this behaviour to a container without the sources, just an alpine image with the binary inside (todo)

## test

Runs tests. Needs a running elasticsearch on 192.168.99.100 (see test-elasticsearch-container)

## test-elasticsearch-container

Starts an elasticsearch container for test purposes.

Requirements:

* docker-machine default running on 192.168.99.100
* a host entry with elasticsearch.mrsjenkins.test to the docker-machine

## clean

Deletes the mrsjenkins binary and the elasticsearch test container.
