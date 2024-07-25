package build

import "github.com/bavix/vakeel-way/internal/infra/repositories"

// WebhookRepository returns a new instance of the WebhookStubRepository with
// the webhook data loaded from the configuration.
//
// It uses the webhook data from the configuration to create a new instance of
// WebhookStubRepository. The webhook data is loaded from the configuration and
// converted to a map using the AsMap method of the Webhooks type.
//
// Parameters:
//   - None
//
// Returns:
//   - *repositories.WebhookStubRepository: A new instance of WebhookStubRepository
//     with the webhook data loaded from the configuration.
func (b *Builder) WebhookRepository() *repositories.WebhookStubRepository {
	// Load the webhook data from the configuration.
	webhookData := b.config.Webhooks.AsMap()

	// Create a new instance of WebhookStubRepository with the webhook data.
	return repositories.NewWebhookRepository(webhookData)
}
