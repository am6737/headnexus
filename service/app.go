package service

import (
	"context"
	"github.com/am6737/headnexus/app"
	"github.com/am6737/headnexus/app/host"
	hostCommand "github.com/am6737/headnexus/app/host/command"
	hostsQuery "github.com/am6737/headnexus/app/host/query"
	"github.com/am6737/headnexus/app/network"
	networkCommand "github.com/am6737/headnexus/app/network/command"
	networkQuery "github.com/am6737/headnexus/app/network/query"

	ruleCommand "github.com/am6737/headnexus/app/rule/command"
	ruleQuery "github.com/am6737/headnexus/app/rule/query"

	"github.com/am6737/headnexus/app/rule"
	"github.com/am6737/headnexus/config"
	"github.com/am6737/headnexus/infra/persistence"
	"github.com/sirupsen/logrus"
)

func NewApplication(ctx context.Context, cfg *config.Config, logger *logrus.Logger) *app.Application {

	mongodbConn, err := persistence.ConnectMongoDB(cfg.Persistence.Url)
	if err != nil {
		panic(err)
	}

	networkRepo := persistence.NewNetworkMongoDBRepository(mongodbConn, cfg.Persistence.DB)
	nc := networkCommand.NewNetworkHandler(networkRepo, logger)
	nq := networkQuery.NewNetworkHandler(networkRepo, logger)

	hostRepo := persistence.NewHostMongodbRepository(mongodbConn, cfg.Persistence.DB)
	hc := hostCommand.NewHostHandler(hostRepo, logger, nc)
	hq := hostsQuery.NewHostHandler(hostRepo, logger, nc)

	ruleRepo := persistence.NewRuleMongodbRepository(mongodbConn, cfg.Persistence.DB)
	rc := ruleCommand.NewRuleHandler(ruleRepo, logger)
	rq := ruleQuery.NewRuleHandler(ruleRepo, logger)

	return &app.Application{
		Host: host.Application{
			Commands: host.Commands{
				Handler: hc,
			},
			Queries: host.Queries{
				Handler: hq,
			},
		},
		Network: network.Application{
			Commands: network.Commands{
				Handler: nc,
			},
			Queries: network.Queries{
				Handler: nq,
			},
		},
		Rule: rule.Application{
			Commands: rule.Commands{
				Handler: rc,
			},
			Queries: rule.Queries{
				Handler: rq,
			},
		},
	}
}
