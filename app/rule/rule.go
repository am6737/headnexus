package rule

import (
	"context"
	"encoding/json"
	"github.com/am6737/headnexus/domain/rule/entity"
)

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

type CreateRule struct {
	Type        string
	Action      string
	UserID      string
	Name        string
	Description string
	Port        string
	Proto       string
	Host        []string
}

type DeleteRule struct {
	UserID string
	ID     string
}

type UpdateRule struct {
	Type        uint8
	ID          int64
	UserID      string
	Description string
	Port        string
	Proto       string
	Action      string
	Host        []string
}

type CommandHandler interface {
	Create(ctx context.Context, cmd *CreateRule) (*entity.Rule, error)
	Delete(ctx context.Context, cmd *DeleteRule) error
	Update(ctx context.Context, cmd *UpdateRule) (*entity.Rule, error)
}

type GetRule struct {
	ID int64
}

type FindRule struct {
	Name     string
	HostID   string
	PageSize int
	PageNum  int
}

func (fr FindRule) String() string {
	marshal, err := json.Marshal(fr)
	if err != nil {
		return ""
	}
	return string(marshal)
}

type QueryHandler interface {
	Get(ctx context.Context, query *GetRule) (*entity.Rule, error)
	Find(ctx context.Context, query *FindRule) ([]*entity.Rule, error)
}
