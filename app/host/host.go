package host

import (
	"context"
	ctime "github.com/am6737/headnexus/common/time"
	"github.com/am6737/headnexus/domain/host/entity"
)

type CreateHost struct {
	Name            string
	NetworkID       string
	Role            string
	StaticAddresses []string
	Rules           []string
	IPAddress       string
	Port            int
	IsLighthouse    bool
	Tags            map[string]interface{}
}

type DeleteHost struct {
	UserID string
	ID     string
}

type UpdateHost struct {
	ID              string
	Name            string
	NetworkID       string
	Role            string
	IPAddress       string
	StaticAddresses []string
	Port            int
	IsLighthouse    bool
	Tags            map[string]interface{} `json:"tags"`
}

type GenerateEnrollCode struct {
	UserID string
	HostID string
}

type GenerateEnrollCodeResponse struct {
	Code            string
	LifetimeSeconds int64
}

type GenerateEnrollCodeCheck struct {
	Code   string
	HostID string
	UserID string
}

type GenerateEnrollCodeCheckResponse struct {
	Exists bool
}

type EnrollHost struct {
	HostID     string
	EnrollCode string
	UserID     string
}

type EnrollHostResponse struct {
	EnrollAt   int64
	LastSeenAt string
	Online     bool
}

type AddHostRule struct {
	UserID string
	HostID string
	Rules  []string
}

type CommandHandler interface {
	Create(ctx context.Context, cmd *CreateHost) (*Host, error)
	Delete(ctx context.Context, cmd *DeleteHost) error
	Update(ctx context.Context, cmd *UpdateHost) (*Host, error)

	AddHostRule(ctx context.Context, cmd *AddHostRule) ([]*entity.Rule, error)

	EnrollHost(ctx context.Context, cmd *EnrollHost) (*EnrollHostResponse, error)
	GenerateEnrollCode(ctx context.Context, cmd *GenerateEnrollCode) (*GenerateEnrollCodeResponse, error)
	EnrollCodeCheck(ctx context.Context, cmd *GenerateEnrollCodeCheck) (*GenerateEnrollCodeCheckResponse, error)
}

type GetHost struct {
	ID     string
	UserID string
}

type FindHost struct {
	Limit        int
	Offset       int
	IsLighthouse bool
	NetworkID    string
	IPAddress    string
	Role         string
	Name         string
	Filters      map[string]interface{}
	Sort         map[string]interface{}
}

// Host 表示一个主机
type Host struct {
	ID              string
	Name            string
	NetworkID       string
	Role            string
	IPAddress       string
	StaticAddresses []string

	Port         int
	IsLighthouse bool
	Online       bool
	CreatedAt    string
	LastSeenAt   string
	EnrollAt     string

	Tags   map[string]interface{}
	Config string
}

type ListHostRules struct {
	UserID   string
	HostID   string
	PageSize int
	PageNum  int
}

type QueryHandler interface {
	Get(ctx context.Context, query *GetHost) (*Host, error)
	Find(ctx context.Context, query *FindHost) ([]*Host, error)

	ListHostRules(ctx context.Context, query *ListHostRules) ([]*entity.Rule, error)
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

func ConvertEntityToHost(h *entity.Host) *Host {
	var online bool
	if h.Status == entity.Online {
		online = true
	}
	return &Host{
		ID:              h.ID,
		Name:            h.Name,
		IPAddress:       h.IPAddress,
		Port:            h.Port,
		IsLighthouse:    h.IsLighthouse,
		Online:          online,
		Role:            h.Role,
		LastSeenAt:      ctime.FormatTimeSince(h.LastSeenAt),
		NetworkID:       h.NetworkID,
		StaticAddresses: h.StaticAddresses,
		Tags:            h.Tags,
		CreatedAt:       ctime.FormatTimestamp(h.CreatedAt),
	}
}
