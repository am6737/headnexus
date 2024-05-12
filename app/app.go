package app

import (
	"github.com/am6737/headnexus/app/host"
	"github.com/am6737/headnexus/app/network"
	"github.com/am6737/headnexus/app/rule"
	"github.com/am6737/headnexus/app/user"
	pkgjwt "github.com/am6737/headnexus/pkg/jwt"
)

type Application struct {
	Host      host.Application
	Network   network.Application
	Rule      rule.Application
	User      user.Application
	JwtConfig *pkgjwt.JWTConfig
}
