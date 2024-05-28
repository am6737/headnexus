package query

import (
	"context"
	"github.com/am6737/headnexus/domain/host/entity"
	"github.com/am6737/headnexus/infra/persistence"
	"github.com/am6737/headnexus/pkg/decorator"
	"github.com/sirupsen/logrus"
)

type GetHost struct {
	UserID string
	HostID string
}

// Host 表示一个主机
type Host struct {
	ID              string
	Name            string
	NetworkID       string
	Role            string
	IPAddress       string
	PublicIP        string
	StaticAddresses map[string]string
	Port            int
	IsLighthouse    bool
	Online          bool
	CreatedAt       string
	LastSeenAt      string
	EnrollAt        string

	Tags   map[string]interface{}
	Config string
}

type GetHostHandler decorator.CommandHandler[*GetHost, *entity.Host]

func NewGetHostHandler(
	logger *logrus.Logger,
	repos persistence.Repositories,
) GetHostHandler {
	return &getHostHandler{
		logger: logger,
		repos:  repos,
	}
}

type getHostHandler struct {
	logger *logrus.Logger
	repos  persistence.Repositories
}

func (h *getHostHandler) Handle(ctx context.Context, query *GetHost) (*entity.Host, error) {
	r1, err := h.repos.HostRepo.Get(ctx, query.HostID)
	if err != nil {
		h.logger.Errorf("error getting host: %v", err)
		return nil, err
	}

	//r2 := ConvertEntityToHost(r1)

	return r1, nil
}
