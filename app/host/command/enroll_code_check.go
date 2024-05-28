package command

import (
	"context"
	"github.com/am6737/headnexus/domain/host/entity"
	"github.com/am6737/headnexus/infra/persistence"
	"github.com/am6737/headnexus/pkg/decorator"
	"github.com/sirupsen/logrus"
)

type GenerateEnrollCodeCheck struct {
	Code   string
	HostID string
	UserID string
}

type GenerateEnrollCodeCheckResponse struct {
	Exists bool
}

type EnrollCodeCheckHandler decorator.CommandHandler[*GenerateEnrollCodeCheck, *GenerateEnrollCodeCheckResponse]

func NewEnrollCodeCheckHandler(
	logger *logrus.Logger,
	repos persistence.Repositories,
) EnrollCodeCheckHandler {
	return &enrollCodeCheckHandler{
		logger: logger,
		repos:  repos,
	}
}

type enrollCodeCheckHandler struct {
	logger *logrus.Logger
	repos  persistence.Repositories
}

func (h *enrollCodeCheckHandler) Handle(ctx context.Context, cmd *GenerateEnrollCodeCheck) (*GenerateEnrollCodeCheckResponse, error) {
	r, err := h.repos.HostRepo.GetEnrollHost(ctx, &entity.GetEnrollHost{
		HostID: cmd.HostID,
	})
	if err != nil {
		h.logger.Error("failed to get host", err)
		return nil, err
	}

	if r.EnrollAt > 0 {
		return &GenerateEnrollCodeCheckResponse{Exists: false}, nil
	}

	return &GenerateEnrollCodeCheckResponse{Exists: true}, nil
}
