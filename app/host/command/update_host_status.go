package command

import (
	"context"
	"github.com/am6737/headnexus/domain/host/entity"
	"github.com/am6737/headnexus/infra/persistence"
	"github.com/am6737/headnexus/pkg/decorator"
	"github.com/sirupsen/logrus"
)

type UpdateHostStatus struct {
	ID              string
	Name            string
	NetworkID       string
	Role            string
	IPAddress       string
	StaticAddresses map[string]string
	Port            int
	IsLighthouse    bool
	Tags            map[string]interface{} `json:"tags"`
}

type UpdateHostStatusHandler decorator.CommandHandler[*UpdateHost, *entity.Host]

func NewUpdateHostStatusHandler(
	logger *logrus.Logger,
	repos persistence.Repositories,
) UpdateHostStatusHandler {
	return &updateHostStatusHandler{
		logger: logger,
		repos:  repos,
	}
}

type updateHostStatusHandler struct {
	logger *logrus.Logger
	repos  persistence.Repositories
}

func (h *updateHostStatusHandler) Handle(ctx context.Context, cmd *UpdateHost) (*entity.Host, error) {
	//TODO implement me
	panic("implement me")
}
