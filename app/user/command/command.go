package command

import (
	"github.com/am6737/headnexus/app/user"
	"github.com/am6737/headnexus/config"
	"github.com/am6737/headnexus/domain/user/repository"
	"github.com/am6737/headnexus/pkg/email"
	pkgjwt "github.com/am6737/headnexus/pkg/jwt"
	"github.com/sirupsen/logrus"
)

var _ user.CommandHandler = &UserHandler{}

func NewUserHandler(repo repository.UserRepository, logger *logrus.Logger, jwtConfig *pkgjwt.JWTConfig, emailClient *email.EmailClient, listen config.ListenConfig) *UserHandler {
	return &UserHandler{
		logger:      logger,
		repo:        repo,
		jwtConfig:   jwtConfig,
		emailClient: emailClient,
		listen:      listen,
	}
}

type UserHandler struct {
	logger      *logrus.Logger
	repo        repository.UserRepository
	jwtConfig   *pkgjwt.JWTConfig
	emailClient *email.EmailClient
	listen      config.ListenConfig
}
