package network

import (
	"context"
	"github.com/am6737/headnexus/domain/network/entity"
)

type CreateNetwork struct {
	Name string `json:"name"`
	Cidr string `json:"cidr"`
}

type DeleteNetwork struct {
	ID string `json:"id"`
}

type UpdateNetwork struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CommandHandler interface {
	Create(ctx context.Context, cmd *CreateNetwork) (*entity.Network, error)
	Delete(ctx context.Context, cmd *DeleteNetwork) error
	Update(ctx context.Context, cmd *UpdateNetwork) (*entity.Network, error)

	AllocateAutoAddress(ctx context.Context, networkID string) (string, error)
	AllocateStaticAddress(ctx context.Context, networkID string, addr string) (string, error)
	RecycleAddress(ctx context.Context, networkID string, addr string) error
	UsedAddresses(ctx context.Context, networkID string) ([]string, error)
	AvailableAddresses(ctx context.Context, networkID string) ([]string, error)
}

type GetNetwork struct {
	ID string `json:"id"`
}

type FindNetwork struct {
	Name          string
	Cidr          string
	SortDirection string
	IncludeCounts bool
	PageSize      int
	PageNumber    int
}

type QueryHandler interface {
	Get(ctx context.Context, query *GetNetwork) (*entity.Network, error)
	Find(ctx context.Context, query *FindNetwork) ([]*entity.Network, error)
}

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	Handler CommandHandler
}

type Queries struct {
	Handler QueryHandler
}
