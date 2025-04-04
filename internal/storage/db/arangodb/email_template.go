package arangodb

import (
	"context"
	"fmt"
	"time"

	arangoDriver "github.com/arangodb/go-driver"
	"github.com/google/uuid"

	"github.com/authorizerdev/authorizer/internal/graph/model"
	"github.com/authorizerdev/authorizer/internal/storage/schemas"
)

// AddEmailTemplate to add EmailTemplate
func (p *provider) AddEmailTemplate(ctx context.Context, emailTemplate *schemas.EmailTemplate) (*schemas.EmailTemplate, error) {
	if emailTemplate.ID == "" {
		emailTemplate.ID = uuid.New().String()
		emailTemplate.Key = emailTemplate.ID
	}
	emailTemplate.Key = emailTemplate.ID
	emailTemplate.CreatedAt = time.Now().Unix()
	emailTemplate.UpdatedAt = time.Now().Unix()
	emailTemplateCollection, _ := p.db.Collection(ctx, schemas.Collections.EmailTemplate)
	meta, err := emailTemplateCollection.CreateDocument(ctx, emailTemplate)
	if err != nil {
		return nil, err
	}

	emailTemplate.Key = meta.Key
	emailTemplate.ID = meta.ID.String()
	return emailTemplate, nil
}

// UpdateEmailTemplate to update EmailTemplate
func (p *provider) UpdateEmailTemplate(ctx context.Context, emailTemplate *schemas.EmailTemplate) (*schemas.EmailTemplate, error) {
	emailTemplate.UpdatedAt = time.Now().Unix()
	emailTemplateCollection, _ := p.db.Collection(ctx, schemas.Collections.EmailTemplate)
	meta, err := emailTemplateCollection.UpdateDocument(ctx, emailTemplate.Key, emailTemplate)
	if err != nil {
		return nil, err
	}
	emailTemplate.Key = meta.Key
	emailTemplate.ID = meta.ID.String()
	return emailTemplate, nil
}

// ListEmailTemplates to list EmailTemplate
func (p *provider) ListEmailTemplate(ctx context.Context, pagination *model.Pagination) ([]*schemas.EmailTemplate, *model.Pagination, error) {
	emailTemplates := []*schemas.EmailTemplate{}
	query := fmt.Sprintf("FOR d in %s SORT d.created_at DESC LIMIT %d, %d RETURN d", schemas.Collections.EmailTemplate, pagination.Offset, pagination.Limit)
	sctx := arangoDriver.WithQueryFullCount(ctx)
	cursor, err := p.db.Query(sctx, query, nil)
	if err != nil {
		return nil, nil, err
	}
	defer cursor.Close()
	paginationClone := pagination
	paginationClone.Total = cursor.Statistics().FullCount()
	for {
		var emailTemplate *schemas.EmailTemplate
		meta, err := cursor.ReadDocument(ctx, &emailTemplate)
		if arangoDriver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, nil, err
		}
		if meta.Key != "" {
			emailTemplates = append(emailTemplates, emailTemplate)
		}
	}
	return emailTemplates, paginationClone, nil
}

// GetEmailTemplateByID to get EmailTemplate by id
func (p *provider) GetEmailTemplateByID(ctx context.Context, emailTemplateID string) (*schemas.EmailTemplate, error) {
	var emailTemplate *schemas.EmailTemplate
	query := fmt.Sprintf("FOR d in %s FILTER d._id == @id RETURN d", schemas.Collections.EmailTemplate)
	bindVars := map[string]interface{}{
		"id": emailTemplateID,
	}
	cursor, err := p.db.Query(ctx, query, bindVars)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()
	for {
		if !cursor.HasMore() {
			if emailTemplate == nil {
				return nil, fmt.Errorf("email template not found")
			}
			break
		}
		_, err := cursor.ReadDocument(ctx, &emailTemplate)
		if err != nil {
			return nil, err
		}
	}
	return emailTemplate, nil
}

// GetEmailTemplateByEventName to get EmailTemplate by event_name
func (p *provider) GetEmailTemplateByEventName(ctx context.Context, eventName string) (*schemas.EmailTemplate, error) {
	var emailTemplate *schemas.EmailTemplate
	query := fmt.Sprintf("FOR d in %s FILTER d.event_name == @event_name RETURN d", schemas.Collections.EmailTemplate)
	bindVars := map[string]interface{}{
		"event_name": eventName,
	}
	cursor, err := p.db.Query(ctx, query, bindVars)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()
	for {
		if !cursor.HasMore() {
			if emailTemplate == nil {
				return nil, fmt.Errorf("email template not found")
			}
			break
		}
		_, err := cursor.ReadDocument(ctx, &emailTemplate)
		if err != nil {
			return nil, err
		}
	}
	return emailTemplate, nil
}

// DeleteEmailTemplate to delete EmailTemplate
func (p *provider) DeleteEmailTemplate(ctx context.Context, emailTemplate *schemas.EmailTemplate) error {
	eventTemplateCollection, _ := p.db.Collection(ctx, schemas.Collections.EmailTemplate)
	_, err := eventTemplateCollection.RemoveDocument(ctx, emailTemplate.Key)
	if err != nil {
		return err
	}
	return nil
}
