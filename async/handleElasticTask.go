package async

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

type ProductMapping struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (processor *RedisTaskProccessor) HandleSyncProductTask(ctx context.Context, task *asynq.Task) error {
	var payload SyncDataTask
	err := json.Unmarshal(task.Payload(), &payload)
	if err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}
	// var docs string // id:name
	bulkIndexer, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:  "ecommerce_product",
		Client: &processor.elasticClient,
	})
	if err != nil {
		log.Error().
			Err(err).
			Msgf("failed to init bulk indexer: %w", err)
		return fmt.Errorf("err: %w", err)
	}
	for _, id := range payload.IDs {
		name, err := processor.store.GetProductNameByID(ctx, id)
		if err != nil {
			log.Error().
				Err(err).
				Msgf("failed to get product name: %w", err)
			return fmt.Errorf("err: %w", err)
		}
		product := ProductMapping{
			ID:   id,
			Name: name,
		}
		jProduct, err := json.Marshal(product)
		if err != nil {
			log.Error().
				Err(err).
				Msgf("failed to marshal doc: %w", err)
			return fmt.Errorf("err: %w", err)
		}
		err = bulkIndexer.Add(ctx, esutil.BulkIndexerItem{
			Action: "index",
			Body:   bytes.NewReader(jProduct),
		})
		if err != nil {
			log.Error().
				Err(err).
				Msgf("failed to add doc to bulk indexer: %w", err)
			return fmt.Errorf("err: %w", err)
		}
		// docs = fmt.Sprintf("%s\n%s", docs, jProduct)
	}
	if err := bulkIndexer.Close(ctx); err != nil {
		log.Error().
			Err(err).
			Msgf("failed to add doc to bulk indexer: %w", err)
		return fmt.Errorf("failed to sends item: %w", err)
	}
	return nil
}

type TagMapping struct {
	Tag        string   `json:"tag"`
	ProductIDs []string `json:"product_id"`
}

// ecommerce_tags
func (processor *RedisTaskProccessor) HandleSyncTagTask(ctx context.Context, task *asynq.Task) error {
	var payload SyncDataTask
	err := json.Unmarshal(task.Payload(), &payload)
	if err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", err)
	}
	bulkReq, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:  "ecommerce_tag",
		Client: &processor.elasticClient,
	})
	if err != nil {
		log.Error().
			Err(err).
			Msgf("failed to init bulk indexer: %w", err)
		return fmt.Errorf("err: %w", err)
	}
	for _, id := range payload.IDs {
		tagName, err := processor.store.GetTagNameByID(ctx, id)
		if err != nil {
			log.Error().
				Err(err).
				Msgf("failed to get tag name: %w", err)
			return fmt.Errorf("err: %w", err)
		}
		productIDs, err := processor.store.ListProductIDbyTagID(ctx, id)
		if err != nil {
			log.Error().
				Err(err).
				Msgf("failed to list product ids: %w", err)
			return fmt.Errorf("err: %w", err)
		}
		tag := TagMapping{
			Tag:        tagName,
			ProductIDs: productIDs,
		}
		jTag, err := json.Marshal(tag)
		if err != nil {
			log.Error().
				Err(err).
				Msgf("failed to marshal doc: %w", err)
			return fmt.Errorf("err: %w", err)
		}
		err = bulkReq.Add(ctx, esutil.BulkIndexerItem{
			Action: "index",
			Body:   bytes.NewReader(jTag),
		})
		if err != nil {
			log.Error().
				Err(err).
				Msgf("failed to add doc to bulk indexer: %w", err)
			return fmt.Errorf("failed to sends item: %w", err)
		}
	}
	if err := bulkReq.Close(ctx); err != nil {
		log.Error().
			Err(err).
			Msgf("failed to add doc to bulk indexer: %w", err)
		return fmt.Errorf("failed to sends item: %w", err)
	}

	return nil
}
