package command

import (
	"github.com/am6737/headnexus/app/host"
	"github.com/am6737/headnexus/app/network"
	"github.com/am6737/headnexus/domain/host/repository"
	"github.com/am6737/headnexus/infra/persistence"
	"github.com/sirupsen/logrus"
)

var _ host.CommandHandler = &HostHandler{}

func NewHostHandler(repo repository.HostRepository, hostRuleRepo repository.HostRuleRepository, repos *persistence.Repositories, logger *logrus.Logger, nch network.CommandHandler) *HostHandler {
	return &HostHandler{
		logger:       logger,
		nch:          nch,
		repo:         repo,
		hostRuleRepo: hostRuleRepo,
		repos:        repos,
	}
}

type HostHandler struct {
	logger *logrus.Logger

	nch  network.CommandHandler
	repo repository.HostRepository

	hostRuleRepo repository.HostRuleRepository

	repos *persistence.Repositories
}
