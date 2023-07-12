package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"os"
	"sync/atomic"

	es "github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
)

type IndexWriter[T any] interface {
	WriteToIndex(ctx context.Context, items []T) error
}

type IndexClient[T any] struct {
	index  string
	client *es.Client
}

func (r IndexClient[T]) WriteToIndex(ctx context.Context, items []T) error {
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:  r.index,
		Client: r.client,
	})
	if err != nil {
		return err
	}
	successful := uint64(0)
	for _, event := range items {
		doc, err := json.Marshal(event)
		if err != nil {
			return err
		}
		err = bi.Add(
			ctx,
			esutil.BulkIndexerItem{
				Action: "index",
				Body:   bytes.NewReader(doc),
				OnSuccess: func(ctx context.Context, bii esutil.BulkIndexerItem, biri esutil.BulkIndexerResponseItem) {
					atomic.AddUint64(&successful, 1)
				},
				OnFailure: func(ctx context.Context, bii esutil.BulkIndexerItem, biri esutil.BulkIndexerResponseItem, err error) {
					if err == nil {
						log.Printf("Error: %s: %s\n", biri.Error.Type, biri.Error.Reason)
					}
					log.Printf("Error: %s", err)
				},
			},
		)
		if err != nil {
			return err
		}
	}
	if err = bi.Close(ctx); err != nil {
		return err
	}
	bi.Stats()
	log.Printf("Bulk commit %d out of %d successful...\n: %#v", successful, len(items), bi.Stats())

	return nil
}
func NewIndexClient[T any](index string) (IndexClient[T], error) {
	var (
		w      IndexClient[T]
		err    error
		config es.Config
		ok     bool
	)
	if config.CloudID, ok = os.LookupEnv("ES_CLOUD_ID"); !ok {
		log.Panic("ES_CLOUD_ID setting is mandatory")
	}

	if config.APIKey, ok = os.LookupEnv("ES_API_KEY"); !ok {
		log.Panic("ES_TOKEN setting is mandatory")
	}
	w.index = index
	if w.client, err = es.NewClient(config); err != nil {
		return w, err
	}

	return w, nil
}
