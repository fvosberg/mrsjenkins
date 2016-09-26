package todo

import (
	"time"

	log "github.com/Sirupsen/logrus"
	elastic "gopkg.in/olivere/elastic.v3"
)

type Datastore interface {
	Create(*Todo)
}

type elasticsearchDatastore struct {
	client *elastic.Client
}

func (t *elasticsearchDatastore) Create(todo *Todo) {
	log.Printf("TODO: should create Todo %+v\n", todo)
}

func NewElasticDatastore() Datastore {
	return NewElasticDatastoreWithURL("http://elasticsearch.mrsjenkins.de:9200")
}

func NewElasticDatastoreWithURL(URL string) Datastore {
	datastore := &elasticsearchDatastore{
		client: newElasticsearchClient(URL),
	}
	return datastore
}

func newElasticsearchClient(URL string) *elastic.Client {
	var client *elastic.Client
	var err error
	log.Printf("Trying to initialize an Elasticsearch client on \"%s\"\n", URL)
	for {
		client, err = elastic.NewClient(
			elastic.SetURL(URL),
			elastic.SetInfoLog(log.StandardLogger()),
		)
		if err != nil {
			log.Printf("Error while connecting to elasticsearch on \"%s\": %+v - retrying in 5 seconds\n", URL, err)
			time.Sleep(5 * time.Second)
		} else {
			log.Printf("Initialized Elasticsearch client on \"%s\"\n", URL)
			break
		}
	}
	return client
}
