package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/olivere/elastic/v7"
	"log"
)

const mapping = `
{
	"settings": {
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings": {
		"properties": {
			"id": {
				"type": "keyword"
			},
			"source": {
				"type": "keyword"
			},
			"type": {
				"type": "keyword"
			},
			"method": {
				"type": "keyword"
			},
			"time": {
				"type": "date"
			},
			"request": {
				"type": "text"
			},
			"response": {
				"type": "text"
			},
			"is_error": {
				"type": "keyword"
			},
			"error_message": {
				"type": "text"
			}
		}
	}
}`

type Elastic struct {
	client *elastic.Client
}

func NewElasticSearch(url string) (*Elastic, error) {
	ctx := context.Background()
	client, err := elastic.NewClient(elastic.SetURL(url), elastic.SetSniff(false))
	if err != nil {
		return nil, err
	}

	info, code, err := client.Ping(url).Do(ctx)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)
	exists, err := client.IndexExists("history").Do(ctx)
	if err != nil {
		return nil, err
	}
	if exists {
		resultDelete, err := client.DeleteIndex("history").Do(ctx)
		if err != nil {
			return nil, err
		}
		if !resultDelete.Acknowledged {
			return nil, errors.New("index history was not deleted")
		}

	}

	createIndex, err := client.CreateIndex("history").BodyString(mapping).Do(ctx)
	if err != nil {
		return nil, err
	}
	if !createIndex.Acknowledged {
		// Not acknowledged
	}
	log.Println("Elasticsearch created history index")

	el := &Elastic{
		client: client,
	}
	return el, nil
}

func (el *Elastic) Save(ctx context.Context, msg *History) error {
	log.Printf("History Id: %s recevied, saving to elastic.\n", msg.Id)
	_, err := el.client.Index().
		Index("history").
		Id(msg.Id).
		BodyJson(msg).
		Do(ctx)
	if err != nil {
		return err
	}
	return nil
}
