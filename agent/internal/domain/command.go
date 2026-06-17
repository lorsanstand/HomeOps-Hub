package domain

import (
	"context"

	domainHub "github.com/lorsanstand/HomeOps-Hub/shared/domain"
)

type Command struct {
	domainHub.CapabilityCommand
	Execute func(ctx context.Context, args map[string]string) (string, error)
}
