package command

import (
	"github.com/am6737/headnexus/app/rule"
	"github.com/am6737/headnexus/domain/rule/repository"
	"github.com/sirupsen/logrus"
)

var _ rule.CommandHandler = &RuleHandler{}

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
