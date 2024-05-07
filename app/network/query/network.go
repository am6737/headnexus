package query

import (
	"github.com/am6737/headnexus/app/network"
	"github.com/am6737/headnexus/domain/network/repository"
	"github.com/sirupsen/logrus"
)

var _ network.QueryHandler = &NetworkHandler{}

func NewNetworkHandler(repo repository.NetworkRepository, logger *logrus.Logger) *NetworkHandler {
	return &NetworkHandler{
		repo:   repo,
		logger: logger,
	}
}

type NetworkHandler struct {
	logger *logrus.Logger

	repo repository.NetworkRepository
}
