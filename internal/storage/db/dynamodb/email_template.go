package dynamodb

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/guregu/dynamo"

	"github.com/authorizerdev/authorizer/internal/graph/model"
	"github.com/authorizerdev/authorizer/internal/storage/schemas"
)

// AddEmailTemplate to add EmailTemplate
func (p *provider) AddEmailTemplate(ctx context.Context, emailTemplate *schemas.EmailTemplate) (*schemas.EmailTemplate, error) {
	collection := p.db.Table(schemas.Collections.EmailTemplate)
	if emailTemplate.ID == "" {
		emailTemplate.ID = uuid.New().String()
	}

	emailTemplate.Key = emailTemplate.ID
	emailTemplate.CreatedAt = time.Now().Unix()
	emailTemplate.UpdatedAt = time.Now().Unix()
	err := collection.Put(emailTemplate).RunWithContext(ctx)

	if err != nil {
		return nil, err
	}

	return emailTemplate, nil
}

// UpdateEmailTemplate to update EmailTemplate
func (p *provider) UpdateEmailTemplate(ctx context.Context, emailTemplate *schemas.EmailTemplate) (*schemas.EmailTemplate, error) {
	collection := p.db.Table(schemas.Collections.EmailTemplate)
	emailTemplate.UpdatedAt = time.Now().Unix()
	err := UpdateByHashKey(collection, "id", emailTemplate.ID, emailTemplate)
	if err != nil {
		return nil, err
	}
	return emailTemplate, nil
}

// ListEmailTemplates to list EmailTemplate
func (p *provider) ListEmailTemplate(ctx context.Context, pagination *model.Pagination) ([]*schemas.EmailTemplate, *model.Pagination, error) {
	var emailTemplate *schemas.EmailTemplate
	var iter dynamo.PagingIter
	var lastEval dynamo.PagingKey
	var iteration int64 = 0
	collection := p.db.Table(schemas.Collections.EmailTemplate)
	emailTemplates := []*schemas.EmailTemplate{}
	paginationClone := pagination
	scanner := collection.Scan()
	count, err := scanner.Count()
	if err != nil {
		return nil, nil, err
	}
	for (paginationClone.Offset + paginationClone.Limit) > iteration {
		iter = scanner.StartFrom(lastEval).Limit(paginationClone.Limit).Iter()
		for iter.NextWithContext(ctx, &emailTemplate) {
			if paginationClone.Offset == iteration {
				emailTemplates = append(emailTemplates, emailTemplate)
			}
		}
		lastEval = iter.LastEvaluatedKey()
		iteration += paginationClone.Limit
	}
	paginationClone.Total = count
	return emailTemplates, paginationClone, nil
}

// GetEmailTemplateByID to get EmailTemplate by id
func (p *provider) GetEmailTemplateByID(ctx context.Context, emailTemplateID string) (*schemas.EmailTemplate, error) {
	collection := p.db.Table(schemas.Collections.EmailTemplate)
	var emailTemplate *schemas.EmailTemplate
	err := collection.Get("id", emailTemplateID).OneWithContext(ctx, &emailTemplate)
	if err != nil {
		return nil, err
	}
	return emailTemplate, nil
}

// GetEmailTemplateByEventName to get EmailTemplate by event_name
func (p *provider) GetEmailTemplateByEventName(ctx context.Context, eventName string) (*schemas.EmailTemplate, error) {
	collection := p.db.Table(schemas.Collections.EmailTemplate)
	var emailTemplates []*schemas.EmailTemplate
	var emailTemplate *schemas.EmailTemplate
	err := collection.Scan().Index("event_name").Filter("'event_name' = ?", eventName).Limit(1).AllWithContext(ctx, &emailTemplates)
	if err != nil {
		return nil, err
	}
	if len(emailTemplates) == 0 {
		return nil, errors.New("no record found")

	}
	emailTemplate = emailTemplates[0]
	return emailTemplate, nil
}

// DeleteEmailTemplate to delete EmailTemplate
func (p *provider) DeleteEmailTemplate(ctx context.Context, emailTemplate *schemas.EmailTemplate) error {
	collection := p.db.Table(schemas.Collections.EmailTemplate)
	err := collection.Delete("id", emailTemplate.ID).RunWithContext(ctx)
	if err != nil {
		return err
	}

	return nil
}
