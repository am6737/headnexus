package repository

import (
	"context"
	"github.com/am6737/headnexus/domain/host/entity"
)

type HostRuleRepository interface {
	AddHostRule(ctx context.Context, hostID string, ruleID ...string) error
	ListHostRule(ctx context.Context, opts *entity.ListHostRuleOptions) ([]*entity.HostRuleRelation, error)
}
