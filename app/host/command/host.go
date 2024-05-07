package command

import (
	"github.com/am6737/headnexus/app/host"
	"github.com/am6737/headnexus/app/network"
	"github.com/am6737/headnexus/domain/host/repository"
	"github.com/sirupsen/logrus"
)

var _ host.CommandHandler = &HostHandler{}

func NewHostHandler(repo repository.HostRepository, logger *logrus.Logger, nch network.CommandHandler) *HostHandler {
	return &HostHandler{
		logger: logger,
		nch:    nch,
		repo:   repo,
	}
}

type HostHandler struct {
	logger *logrus.Logger

	nch  network.CommandHandler
	repo repository.HostRepository
}
