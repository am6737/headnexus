package query

import (
	"github.com/am6737/headnexus/app/host"
	"github.com/am6737/headnexus/app/network"
	"github.com/am6737/headnexus/domain/host/repository"
	"github.com/sirupsen/logrus"
)

var _ host.QueryHandler = &HostHandler{}

func NewHostHandler(repo repository.HostRepository, logger *logrus.Logger, nch network.CommandHandler) *HostHandler {
	return &HostHandler{
		repo:   repo,
		logger: logger,
		nch:    nch,
	}
}

type HostHandler struct {
	logger *logrus.Logger

	nch  network.CommandHandler
	repo repository.HostRepository
}
