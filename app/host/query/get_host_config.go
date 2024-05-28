package query

import (
	"context"
	"github.com/am6737/headnexus/infra/persistence"
	"github.com/am6737/headnexus/pkg/code"
	"github.com/am6737/headnexus/pkg/decorator"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type GetHostConfig struct {
	UserID string
	HostID string
}

type GetHostConfigResponse struct {
	Config string
}

type GetHostConfigHandler decorator.CommandHandler[*GetHostConfig, *GetHostConfigResponse]

func NewGetHostConfigHandler(
	logger *logrus.Logger,
	repos persistence.Repositories,
) GetHostConfigHandler {
	return &getHostConfigHandler{
		logger: logger,
		repos:  repos,
	}
}

type getHostConfigHandler struct {
	logger *logrus.Logger
	repos  persistence.Repositories
}

func (h *getHostConfigHandler) Handle(ctx context.Context, cmd *GetHostConfig) (*GetHostConfigResponse, error) {
	host, err := h.repos.HostRepo.Get(ctx, cmd.HostID)
	if err != nil {
		h.logger.WithError(err).Error("failed to get host")
		return nil, err
	}
	if host.Owner != cmd.UserID {
		return nil, code.Forbidden
	}

	if host.EnrollAt == 0 {
		return nil, code.StatusNotAvailable.SetMessage("主机未注册，无法获取配置")
	}

	marshal, err := yaml.Marshal(&host.Config)
	if err != nil {
		h.logger.WithError(err).Error("failed to marshal host config")
		return nil, err
	}

	return &GetHostConfigResponse{
		Config: string(marshal),
	}, nil
}
