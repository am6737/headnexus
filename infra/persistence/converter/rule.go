package converter

import (
	"github.com/am6737/headnexus/domain/host/entity"
	"github.com/am6737/headnexus/infra/persistence/po"
)

func RuleEntityToPO(e *entity.Rule) *po.Rule {
	m := &po.Rule{}
	m.ID = e.ID
	m.Name = e.Name
	m.Description = e.Description
	m.Port = e.Port
	m.Proto = string(e.Proto)
	m.Action = uint8(e.Action)
	m.Host = e.Host
	return m
}

func RulePOToEntity(po *po.Rule) *entity.Rule {
	return &entity.Rule{
		UserID:      po.UserID,
		ID:          po.ID,
		Type:        entity.RuleType(po.Type),
		CreatedAt:   po.CreatedAt,
		Name:        po.Name,
		Description: po.Description,
		Port:        po.Port,
		Proto:       entity.RuleProto(po.Proto),
		Action:      entity.RuleAction(po.Action),
		Host:        po.Host,
	}
}
