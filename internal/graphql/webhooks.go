package graphql

import (
	"context"
	"fmt"

	"github.com/authorizerdev/authorizer/internal/graph/model"
	"github.com/authorizerdev/authorizer/internal/utils"
)

// Webhooks is the method to list webhooks
// Permission: authorizer:admin
func (g *graphqlProvider) Webhooks(ctx context.Context, params *model.PaginatedInput) (*model.Webhooks, error) {
	log := g.Log.With().Str("func", "Webhooks").Logger()
	gc, err := utils.GinContextFromContext(ctx)
	if err != nil {
		log.Debug().Err(err).Msg("Failed to get GinContext")
		return nil, err
	}
	if !g.TokenProvider.IsSuperAdmin(gc) {
		log.Debug().Msg("Not logged in as super admin")
		return nil, fmt.Errorf("unauthorized")
	}

	pagination := utils.GetPagination(params)
	webhooks, pagination, err := g.StorageProvider.ListWebhook(ctx, pagination)
	if err != nil {
		log.Debug().Err(err).Msg("failed ListWebhook")
		return nil, err
	}
	res := make([]*model.Webhook, len(webhooks))
	for i, webhook := range webhooks {
		res[i] = webhook.AsAPIWebhook()
	}
	return &model.Webhooks{
		Pagination: pagination,
		Webhooks:   res,
	}, nil
}
