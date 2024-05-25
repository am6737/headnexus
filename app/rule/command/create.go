package command

import (
	"context"
	"errors"
	"github.com/am6737/headnexus/app/rule"
	"github.com/am6737/headnexus/domain/host/entity"
)

func (h *RuleHandler) Create(ctx context.Context, cmd *rule.CreateRule) (*entity.Rule, error) {
	if cmd == nil {
		return nil, errors.New("command is nil")
	}

	//find, err := h.repo.Find(ctx, &entity.RuleFindOptions{
	//	Name: cmd.Name,
	//})
	//if err != nil {
	//	return nil, err
	//}
	//
	//if len(find) > 0 {
	//	return nil, errors.New("rule already exists")
	//}

	action, err := entity.ParseRuleAction(cmd.Action)
	if err != nil {
		return nil, err
	}

	_type, err := entity.ParseRuleType(cmd.Type)
	if err != nil {
		return nil, err
	}

	e := &entity.Rule{
		Type:        _type,
		Name:        cmd.Name,
		Description: cmd.Description,
		Port:        cmd.Port,
		Proto:       entity.ConvertToRuleProto(cmd.Proto),
		Action:      action,
		Host:        cmd.Host,
	}

	create, err := h.repo.Create(ctx, cmd.UserID, e)
	if err != nil {
		return nil, err
	}

	return create, nil
}
