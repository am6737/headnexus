package command

import (
	"context"
	"github.com/am6737/headnexus/config"
	"github.com/am6737/headnexus/domain/host/entity"
	"github.com/am6737/headnexus/infra/persistence"
	"github.com/am6737/headnexus/pkg/code"
	"github.com/am6737/headnexus/pkg/decorator"
	"github.com/sirupsen/logrus"
	"strings"
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
	host, err := h.repos.HostRepo.Get(ctx, cmd.HostID)
	if err != nil {
		h.logger.WithError(err).Errorf("failed to get host")
		return err
	}
	if host.Owner != cmd.UserID {
		return code.Forbidden
	}

	if err := h.repos.HostRuleRepo.DeleteHostRule(ctx, cmd.HostID, cmd.Rules...); err != nil {
		h.logger.WithError(err).Errorf("failed to delete host rule")
		return err
	}

	hostRules, err := h.repos.HostRuleRepo.ListHostRule(ctx, &entity.ListHostRuleOptions{
		HostID: cmd.HostID,
	})
	if err != nil {
		h.logger.WithError(err).Errorf("failed to find rule")
		return err
	}

	var ruleids []string
	for _, rule := range hostRules {
		ruleids = append(ruleids, rule.RuleID)
	}

	var rules []*entity.Rule
	if len(ruleids) > 0 {
		rules, err = h.repos.RuleRepo.Gets(ctx, cmd.UserID, ruleids)
		if err != nil {
			h.logger.WithError(err).Errorf("failed to find rule")
			return err
		}
	}

	host.Config.Outbound = append([]config.OutboundRule{}, config.DefaultOutbound...)
	host.Config.Inbound = []config.InboundRule{}

	for _, rule := range rules {
		switch rule.Type {
		case entity.RuleTypeOutbound:
			host.Config.Outbound = append(host.Config.Outbound, config.OutboundRule{
				Port:   rule.Port,
				Proto:  rule.Proto.String(),
				Host:   strings.Split(rule.Host, ","),
				Action: rule.Action.String(),
			})
		case entity.RuleTypeInbound:
			host.Config.Inbound = append(host.Config.Inbound, config.InboundRule{
				Port:   rule.Port,
				Proto:  rule.Proto.String(),
				Host:   strings.Split(rule.Host, ","),
				Action: rule.Action.String(),
			})
		}
	}

	_, err = h.repos.HostRepo.Update(ctx, host)
	if err != nil {
		h.logger.WithError(err).Error("update host failed")
		return nil
	}

	return err
}
