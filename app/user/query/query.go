package query

import (
	"context"
	"errors"
	"github.com/am6737/headnexus/app/user"
	"github.com/am6737/headnexus/domain/user/entity"
	"github.com/am6737/headnexus/domain/user/repository"
	"github.com/sirupsen/logrus"
)

var _ user.QueryHandler = &UserHandler{}

func NewUserHandler(repo repository.UserRepository, logger *logrus.Logger) *UserHandler {
	return &UserHandler{
		logger: logger,
		repo:   repo,
	}
}

type UserHandler struct {
	logger *logrus.Logger

	repo repository.UserRepository
}

func (h *UserHandler) Get(ctx context.Context, query *user.GetUser) (*entity.User, error) {
	get, err := h.repo.Get(ctx, query.ID)
	if err != nil {
		return nil, err
	}
	return get, nil
}

func (h *UserHandler) Find(ctx context.Context, query *user.FindUser) ([]*entity.User, error) {
	h.logger.WithField("query", query).Info("Find req")
	if query == nil {
		return nil, errors.New("query is nil")
	}

	find, err := h.repo.Find(ctx, &entity.FindOptions{
		Email:        query.Email,
		Verification: query.Verification,
		Token:        query.Token,
	})
	if err != nil {
		h.logger.Errorf("error finding rule: %v", err)
		return nil, err
	}

	return find, nil
}
