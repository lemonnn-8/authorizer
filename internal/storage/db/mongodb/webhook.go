package mongodb

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/authorizerdev/authorizer/internal/graph/model"
	"github.com/authorizerdev/authorizer/internal/storage/schemas"
)

// AddWebhook to add webhook
func (p *provider) AddWebhook(ctx context.Context, webhook *schemas.Webhook) (*schemas.Webhook, error) {
	if webhook.ID == "" {
		webhook.ID = uuid.New().String()
	}
	webhook.Key = webhook.ID
	webhook.CreatedAt = time.Now().Unix()
	webhook.UpdatedAt = time.Now().Unix()
	// Add timestamp to make event name unique for legacy version
	webhook.EventName = fmt.Sprintf("%s-%d", webhook.EventName, time.Now().Unix())
	webhookCollection := p.db.Collection(schemas.Collections.Webhook, options.Collection())
	_, err := webhookCollection.InsertOne(ctx, webhook)
	if err != nil {
		return nil, err
	}
	return webhook, nil
}

// UpdateWebhook to update webhook
func (p *provider) UpdateWebhook(ctx context.Context, webhook *schemas.Webhook) (*schemas.Webhook, error) {
	webhook.UpdatedAt = time.Now().Unix()
	// Event is changed
	if !strings.Contains(webhook.EventName, "-") {
		webhook.EventName = fmt.Sprintf("%s-%d", webhook.EventName, time.Now().Unix())
	}
	webhookCollection := p.db.Collection(schemas.Collections.Webhook, options.Collection())
	_, err := webhookCollection.UpdateOne(ctx, bson.M{"_id": bson.M{"$eq": webhook.ID}}, bson.M{"$set": webhook}, options.MergeUpdateOptions())
	if err != nil {
		return nil, err
	}
	return webhook, nil
}

// ListWebhooks to list webhook
func (p *provider) ListWebhook(ctx context.Context, pagination *model.Pagination) ([]*schemas.Webhook, *model.Pagination, error) {
	webhooks := []*schemas.Webhook{}
	opts := options.Find()
	opts.SetLimit(pagination.Limit)
	opts.SetSkip(pagination.Offset)
	opts.SetSort(bson.M{"created_at": -1})
	paginationClone := pagination
	webhookCollection := p.db.Collection(schemas.Collections.Webhook, options.Collection())
	count, err := webhookCollection.CountDocuments(ctx, bson.M{}, options.Count())
	if err != nil {
		return nil, nil, err
	}
	paginationClone.Total = count
	cursor, err := webhookCollection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var webhook *schemas.Webhook
		err := cursor.Decode(&webhook)
		if err != nil {
			return nil, nil, err
		}
		webhooks = append(webhooks, webhook)
	}
	return webhooks, paginationClone, nil
}

// GetWebhookByID to get webhook by id
func (p *provider) GetWebhookByID(ctx context.Context, webhookID string) (*schemas.Webhook, error) {
	var webhook *schemas.Webhook
	webhookCollection := p.db.Collection(schemas.Collections.Webhook, options.Collection())
	err := webhookCollection.FindOne(ctx, bson.M{"_id": webhookID}).Decode(&webhook)
	if err != nil {
		return nil, err
	}
	return webhook, nil
}

// GetWebhookByEventName to get webhook by event_name
func (p *provider) GetWebhookByEventName(ctx context.Context, eventName string) ([]*schemas.Webhook, error) {
	webhooks := []*schemas.Webhook{}
	webhookCollection := p.db.Collection(schemas.Collections.Webhook, options.Collection())
	opts := options.Find()
	opts.SetSort(bson.M{"created_at": -1})
	cursor, err := webhookCollection.Find(ctx, bson.M{"event_name": bson.M{
		"$regex": fmt.Sprintf("^%s", eventName),
	}}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var webhook *schemas.Webhook
		err := cursor.Decode(&webhook)
		if err != nil {
			return nil, err
		}
		webhooks = append(webhooks, webhook)
	}
	return webhooks, nil
}

// DeleteWebhook to delete webhook
func (p *provider) DeleteWebhook(ctx context.Context, webhook *schemas.Webhook) error {
	webhookCollection := p.db.Collection(schemas.Collections.Webhook, options.Collection())
	_, err := webhookCollection.DeleteOne(ctx, bson.M{"_id": webhook.ID}, options.Delete())
	if err != nil {
		return err
	}
	webhookLogCollection := p.db.Collection(schemas.Collections.WebhookLog, options.Collection())
	_, err = webhookLogCollection.DeleteMany(ctx, bson.M{"webhook_id": webhook.ID}, options.Delete())
	if err != nil {
		return err
	}
	return nil
}
