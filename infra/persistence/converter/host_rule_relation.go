package converter

import (
	"github.com/am6737/headnexus/domain/host/entity"
	"github.com/am6737/headnexus/infra/persistence/po"
)

func HostRuleEntityToPO(e *entity.HostRuleRelation) *po.HostRuleRelation {
	return &po.HostRuleRelation{
		ID:        e.ID,
		HostID:    e.HostID,
		RuleID:    e.RuleID,
		CreatedAt: e.CreatedAt,
	}
}

func HostRulePOToEntity(po *po.HostRuleRelation) *entity.HostRuleRelation {
	return &entity.HostRuleRelation{
		ID:        po.ID,
		HostID:    po.HostID,
		RuleID:    po.RuleID,
		CreatedAt: po.CreatedAt,
	}
}
