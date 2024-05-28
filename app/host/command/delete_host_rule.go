package command

import (
	"context"
	"fmt"
	"github.com/am6737/headnexus/infra/persistence"
	"github.com/am6737/headnexus/pkg/code"
	"github.com/am6737/headnexus/pkg/decorator"
	"github.com/sirupsen/logrus"
)

type DeleteHostRule struct {
	UserID string
	HostID string
	Rules  []string
}

type DeleteHostRuleHandler decorator.CommandHandlerNoneResponse[*DeleteHostRule]

func NewDeleteHostRuleHandler(
	logger *logrus.Logger,
	repos persistence.Repositories,
) DeleteHostRuleHandler {
	return &deleteHostRuleHandler{
		logger: logger,
		repos:  repos,
	}
}

type deleteHostRuleHandler struct {
	logger *logrus.Logger
	repos  persistence.Repositories
}

func (h *deleteHostRuleHandler) Handle(ctx context.Context, cmd *DeleteHostRule) error {
	hs, err := h.repos.HostRepo.Get(ctx, cmd.HostID)
	if err != nil {
		h.logger.WithError(err).Errorf("failed to get host")
		return err
	}
	fmt.Println("hs => ", hs)
	if hs.Owner != cmd.UserID {
		return code.Forbidden
	}

	if err := h.repos.HostRuleRepo.DeleteHostRule(ctx, cmd.HostID, cmd.Rules...); err != nil {
		h.logger.WithError(err).Errorf("failed to delete host rule")
		return err
	}

	return err
}
