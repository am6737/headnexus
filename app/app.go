package app

import (
	"github.com/am6737/headnexus/app/host"
	"github.com/am6737/headnexus/app/network"
	"github.com/am6737/headnexus/app/rule"
)

type Application struct {
	Host    host.Application
	Network network.Application
	Rule    rule.Application
}
