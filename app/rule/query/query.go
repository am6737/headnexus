package query

import (
	"context"
	"errors"
	"github.com/am6737/headnexus/app/rule"
	"github.com/am6737/headnexus/domain/rule/entity"
	"github.com/am6737/headnexus/domain/rule/repository"
	"github.com/sirupsen/logrus"
)

var _ rule.QueryHandler = &RuleHandler{}

func NewRuleHandler(repo repository.RuleRepository, logger *logrus.Logger) *RuleHandler {
	return &RuleHandler{
		logger: logger,
		repo:   repo,
	}
}

type RuleHandler struct {
	logger *logrus.Logger

	repo repository.RuleRepository
}

func (h *RuleHandler) Get(ctx context.Context, query *rule.GetRule) (*entity.Rule, error) {
	//TODO implement me
	panic("implement me")
}

func (h *RuleHandler) Find(ctx context.Context, query *rule.FindRule) ([]*entity.Rule, error) {
	h.logger.WithField("query", query).Info("Find req")
	if query == nil {
		return nil, errors.New("query is nil")
	}

	find, err := h.repo.Find(ctx, &entity.FindOptions{
		Name:     query.Name,
		HostID:   query.HostID,
		PageSize: query.PageSize,
		PageNum:  query.PageNum,
	})
	if err != nil {
		h.logger.Errorf("error finding rule: %v", err)
		return nil, err
	}

	return find, nil
}
