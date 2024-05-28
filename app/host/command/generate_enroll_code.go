package command

import (
	"context"
	"github.com/am6737/headnexus/common/code"
	ctime "github.com/am6737/headnexus/common/time"
	"github.com/am6737/headnexus/domain/host/entity"
	"github.com/am6737/headnexus/infra/persistence"
	pcode "github.com/am6737/headnexus/pkg/code"
	"github.com/am6737/headnexus/pkg/decorator"
	"github.com/sirupsen/logrus"
)

type GenerateEnrollCode struct {
	UserID string
	HostID string
}

type GenerateEnrollCodeResponse struct {
	Code            string
	LifetimeSeconds int64
}

type GenerateEnrollCodeHandler decorator.CommandHandler[*GenerateEnrollCode, *GenerateEnrollCodeResponse]

func NewGenerateEnrollCodeHandler(
	logger *logrus.Logger,
	repos persistence.Repositories,
) GenerateEnrollCodeHandler {
	return &generateEnrollCodeHandler{
		logger: logger,
		repos:  repos,
	}
}

type generateEnrollCodeHandler struct {
	logger *logrus.Logger
	repos  persistence.Repositories
}

func (h *generateEnrollCodeHandler) Handle(ctx context.Context, cmd *GenerateEnrollCode) (*GenerateEnrollCodeResponse, error) {
	host, err := h.repos.HostRepo.Get(ctx, cmd.HostID)
	if err != nil {
		return nil, err
	}
	if host.Owner != cmd.UserID {
		return nil, pcode.Forbidden
	}

	//r, err := h.repos.HostRepo.GetEnrollHost(ctx, &entity.GetEnrollHost{
	//	HostID: cmd.HostID,
	//})
	//if err != nil {
	//	h.logger.Error("failed to get host", err)
	//	return nil, err
	//}

	//if r.EnrollAt > 0 {
	//	return nil, errors.New("host already enrolled")
	//}

	ecode := code.GenEnrollCode()

	if err := h.repos.HostRepo.EnrollHost(ctx, &entity.EnrollHost{
		HostID:              cmd.HostID,
		Code:                ecode,
		EnrollCodeExpiredAt: ctime.CurrentTimestampMillis() + lifetimeSeconds,
	}); err != nil {
		h.logger.Errorf("failed to enroll host: %v", err)
		return nil, err
	}

	return &GenerateEnrollCodeResponse{
		Code:            ecode,
		LifetimeSeconds: lifetimeSeconds,
	}, nil
}
