package command

import (
	"context"
	"errors"
	ctime "github.com/am6737/headnexus/common/time"
	"github.com/am6737/headnexus/domain/host/entity"
	"github.com/am6737/headnexus/infra/persistence"
	"github.com/am6737/headnexus/pkg/decorator"
	"github.com/sirupsen/logrus"
)

type AddHostRule struct {
	UserID string
	HostID string
	Rules  []string
}

type AddHostRuleHandler decorator.CommandHandler[*AddHostRule, []*entity.HostRule]

func NewAddHostRuleHandler(
	logger *logrus.Logger,
	repos persistence.Repositories,
) AddHostRuleHandler {
	return &addHostRuleHandler{
		logger: logger,
		repos:  repos,
	}
}

type addHostRuleHandler struct {
	logger *logrus.Logger
	repos  persistence.Repositories
}

func (h *addHostRuleHandler) Handle(ctx context.Context, cmd *AddHostRule) ([]*entity.HostRule, error) {
	h.logger.WithField("cmd", cmd).Info("add host rule")
	if len(cmd.Rules) == 0 {
		return nil, errors.New("no rule ids provided")
	}

	// TODO 判断用户权限

	r, err := h.repos.HostRepo.Get(ctx, cmd.HostID)
	if err != nil {
		h.logger.Error("error getting host")
		return nil, err
	}

	if err := h.repos.HostRuleRepo.AddHostRule(ctx, r.ID, cmd.Rules...); err != nil {
		h.logger.WithError(err).Error("add host rule failed")
		return nil, err
	}

	var rules []*entity.HostRule
	for _, ruleID := range cmd.Rules {
		rule, err := h.repos.RuleRepo.Get(ctx, cmd.UserID, ruleID)
		if err != nil {
			h.logger.WithError(err).Error("get host rule failed")
			return nil, err
		}
		rules = append(rules, &entity.HostRule{
			Type:        rule.Type,
			CreatedAt:   ctime.FormatTimestamp(rule.CreatedAt),
			ID:          rule.ID,
			HostID:      cmd.HostID,
			UserID:      rule.UserID,
			Name:        rule.Name,
			Description: rule.Description,
			Port:        rule.Port,
			Proto:       rule.Proto,
			Action:      rule.Action,
			Host:        rule.Host,
		})
	}

	return rules, nil
}
