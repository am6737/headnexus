package command

import (
	"context"
	"errors"
	"fmt"
	ctime "github.com/am6737/headnexus/common/time"
	"github.com/am6737/headnexus/config"
	"github.com/am6737/headnexus/domain/host/entity"
	"github.com/am6737/headnexus/infra/persistence"
	"github.com/am6737/headnexus/pkg/decorator"
	"github.com/sirupsen/logrus"
	"strings"
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

	host, err := h.repos.HostRepo.Get(ctx, cmd.HostID)
	if err != nil {
		return nil, err
	}

	fmt.Println("更新前配置 host.Config.Outbound => ", host.Config.Outbound)

	if err := h.repos.HostRuleRepo.AddHostRule(ctx, host.ID, cmd.Rules...); err != nil {
		h.logger.WithError(err).Error("add host rule failed")
		return nil, err
	}

	hostRules, err := h.repos.HostRuleRepo.ListHostRule(ctx, &entity.ListHostRuleOptions{
		HostID: host.ID,
	})
	if err != nil {
		h.logger.WithError(err).Error("list host rule failed")
		return nil, err
	}

	host.Config.Inbound = make([]config.InboundRule, 0)

	var rules []*entity.HostRule
	for _, v := range hostRules {
		rule, err := h.repos.RuleRepo.Get(ctx, cmd.UserID, v.RuleID)
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
			Host:        strings.Split(rule.Host, ","),
		})
		if rule.Type == entity.RuleTypeOutbound {
			host.Config.Outbound = append(host.Config.Outbound, config.OutboundRule{
				Port:   rule.Port,
				Proto:  rule.Proto.String(),
				Host:   strings.Split(rule.Host, ","),
				Action: rule.Action.String(),
			})
		}
		if rule.Type == entity.RuleTypeInbound {
			host.Config.Inbound = append(host.Config.Inbound, config.InboundRule{
				Port:   rule.Port,
				Proto:  rule.Proto.String(),
				Host:   strings.Split(rule.Host, ","),
				Action: rule.Action.String(),
			})
		}
	}

	fmt.Println("更新后配置 host.Config.Outbound => ", host.Config.Outbound)

	_, err = h.repos.HostRepo.Update(ctx, host)
	if err != nil {
		h.logger.WithError(err).Error("update host failed")
		return nil, err
	}

	return rules, nil
}
