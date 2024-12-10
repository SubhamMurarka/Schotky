package es

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	models "github.com/SubhamMurarka/Schotky/Models"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

// InitializeElasticsearch sets up the Elasticsearch client.
func InitializeElasticsearch(url string) (*elasticsearch.Client, error) {
	url = "http://" + url
	cfg := elasticsearch.Config{
		Addresses: []string{url},
	}
	client, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	fmt.Println("properly connected")
	return client, nil
}

// var indexCreated bool

// CreateIndexWithMapping ensures the index and its mapping are created only once.
func CreateIndexWithMapping(esClient *elasticsearch.Client) error {
	// Check if the index has already been created
	// if indexCreated {
	// 	log.Println("Index already created, skipping mapping.")
	// 	return nil
	// }

	// Check if the index exists
	res, err := esClient.Indices.Exists([]string{"analytics_data"})
	if err != nil {
		return fmt.Errorf("error checking if index exists: %v", err)
	}
	defer res.Body.Close()

	// If the index exists, no need to create it
	if res.StatusCode == 200 {
		log.Println("Index already exists, skipping index creation.")
		return nil
	}

	// Define the mapping for the index
	mapping := `{
		"mappings": {
			"properties": {
				"country": { 
					"type": "keyword"
				},
				"city": { 
					"type": "keyword"
				},
				"os": { 
					"type": "keyword"
				},
				"device": { 
					"type": "keyword"
				},
				"browser": { 
					"type": "keyword"
				},
				"timestamp": { 
					"type": "date",
					"ignore_malformed": true
				},
				"referrer": { 
					"type": "text", 
					"analyzer": "standard"
				},
				"short_url": { 
					"type": "keyword"
				}
			}
    	}
    }`

	// Create the index with the specified mapping
	resCreate, err := esClient.Indices.Create(
		"analytics_data",
		esClient.Indices.Create.WithBody(strings.NewReader(mapping)),
	)
	if err != nil {
		return fmt.Errorf("error creating index with mapping: %v", err)
	}
	defer resCreate.Body.Close()

	// Check for successful index creation
	if resCreate.IsError() {
		return fmt.Errorf("error creating index with mapping: %s", resCreate.String())
	}

	log.Println("Index created with mapping")
	return nil
}

func WriteBatchToElasticsearch(batch []models.ProcessedData, esClient *elasticsearch.Client) error {
	// Ensure the index and mapping are created first, but only once
	err := CreateIndexWithMapping(esClient)
	if err != nil {
		log.Printf("Error creating index with mapping: %v", err)
		return err
	}

	// Proceed with bulk indexing
	bulkIndexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Client:     esClient,
		Index:      "analytics_data", // Default index
		FlushBytes: 5 * 1024 * 1024,  // Flush every 5MB
	})
	if err != nil {
		return err
	}
	defer bulkIndexer.Close(context.Background())

	// Add items to the BulkIndexer
	for _, record := range batch {
		jsonData, err := json.Marshal(record)
		if err != nil {
			log.Printf("Error marshaling record: %v", err)
			continue
		}

		err = bulkIndexer.Add(context.Background(), esutil.BulkIndexerItem{
			Action: "index",                   // Index action
			Body:   bytes.NewReader(jsonData), // Document to index
			OnSuccess: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem) {
				log.Printf("Successfully indexed document: %s", res.DocumentID)
			},
			OnFailure: func(ctx context.Context, item esutil.BulkIndexerItem, res esutil.BulkIndexerResponseItem, err error) {
				if err != nil {
					log.Printf("Error indexing document: %v", err)
				} else {
					log.Printf("Failed to index document: %s, error: %s", res.DocumentID, res.Error.Reason)
				}
			},
		})
		if err != nil {
			log.Printf("Error adding document to BulkIndexer: %v", err)
		}
	}

	// Report bulk indexing statistics
	stats := bulkIndexer.Stats()
	log.Printf("Bulk indexing completed. Indexed: %d, Created: %d, Errors: %d",
		stats.NumIndexed,
		stats.NumCreated,
		stats.NumFailed)

	return nil
}
