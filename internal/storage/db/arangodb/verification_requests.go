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

// AddVerification to save verification request in database
func (p *provider) AddVerificationRequest(ctx context.Context, verificationRequest *schemas.VerificationRequest) (*schemas.VerificationRequest, error) {
	if verificationRequest.ID == "" {
		verificationRequest.ID = uuid.New().String()
		verificationRequest.Key = verificationRequest.ID
	}
	verificationRequest.CreatedAt = time.Now().Unix()
	verificationRequest.UpdatedAt = time.Now().Unix()
	verificationRequestRequestCollection, _ := p.db.Collection(ctx, schemas.Collections.VerificationRequest)
	meta, err := verificationRequestRequestCollection.CreateDocument(ctx, verificationRequest)
	if err != nil {
		return nil, err
	}
	verificationRequest.Key = meta.Key
	verificationRequest.ID = meta.ID.String()
	return verificationRequest, nil
}

// GetVerificationRequestByToken to get verification request from database using token
func (p *provider) GetVerificationRequestByToken(ctx context.Context, token string) (*schemas.VerificationRequest, error) {
	var verificationRequest *schemas.VerificationRequest
	query := fmt.Sprintf("FOR d in %s FILTER d.token == @token LIMIT 1 RETURN d", schemas.Collections.VerificationRequest)
	bindVars := map[string]interface{}{
		"token": token,
	}
	cursor, err := p.db.Query(ctx, query, bindVars)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()
	for {
		if !cursor.HasMore() {
			if verificationRequest == nil {
				return verificationRequest, fmt.Errorf("verification request not found")
			}
			break
		}
		_, err := cursor.ReadDocument(ctx, &verificationRequest)
		if err != nil {
			return nil, err
		}
	}
	return verificationRequest, nil
}

// GetVerificationRequestByEmail to get verification request by email from database
func (p *provider) GetVerificationRequestByEmail(ctx context.Context, email string, identifier string) (*schemas.VerificationRequest, error) {
	var verificationRequest *schemas.VerificationRequest
	query := fmt.Sprintf("FOR d in %s FILTER d.email == @email FILTER d.identifier == @identifier LIMIT 1 RETURN d", schemas.Collections.VerificationRequest)
	bindVars := map[string]interface{}{
		"email":      email,
		"identifier": identifier,
	}
	cursor, err := p.db.Query(ctx, query, bindVars)
	if err != nil {
		return nil, err
	}
	defer cursor.Close()
	for {
		if !cursor.HasMore() {
			if verificationRequest == nil {
				return verificationRequest, fmt.Errorf("verification request not found")
			}
			break
		}
		_, err := cursor.ReadDocument(ctx, &verificationRequest)
		if err != nil {
			return nil, err
		}
	}
	return verificationRequest, nil
}

// ListVerificationRequests to get list of verification requests from database
func (p *provider) ListVerificationRequests(ctx context.Context, pagination *model.Pagination) ([]*schemas.VerificationRequest, *model.Pagination, error) {
	var verificationRequests []*schemas.VerificationRequest
	sctx := arangoDriver.WithQueryFullCount(ctx)
	query := fmt.Sprintf("FOR d in %s SORT d.created_at DESC LIMIT %d, %d RETURN d", schemas.Collections.VerificationRequest, pagination.Offset, pagination.Limit)
	cursor, err := p.db.Query(sctx, query, nil)
	if err != nil {
		return nil, nil, err
	}
	defer cursor.Close()
	paginationClone := pagination
	paginationClone.Total = cursor.Statistics().FullCount()
	for {
		var verificationRequest *schemas.VerificationRequest
		meta, err := cursor.ReadDocument(ctx, &verificationRequest)

		if arangoDriver.IsNoMoreDocuments(err) {
			break
		} else if err != nil {
			return nil, nil, err
		}

		if meta.Key != "" {
			verificationRequests = append(verificationRequests, verificationRequest)
		}

	}
	return verificationRequests, paginationClone, nil
}

// DeleteVerificationRequest to delete verification request from database
func (p *provider) DeleteVerificationRequest(ctx context.Context, verificationRequest *schemas.VerificationRequest) error {
	collection, _ := p.db.Collection(ctx, schemas.Collections.VerificationRequest)
	_, err := collection.RemoveDocument(ctx, verificationRequest.Key)
	if err != nil {
		return err
	}
	return nil
}
