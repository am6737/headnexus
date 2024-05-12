package command

import (
	"context"
	"errors"
	"fmt"
	"github.com/am6737/headnexus/app/user"
	"github.com/am6737/headnexus/domain/user/entity"
)

func (h *UserHandler) Create(ctx context.Context, cmd *user.CreateUser) (*entity.User, error) {
	if cmd == nil {
		return nil, errors.New("command is nil")
	}

	find, err := h.repo.Find(ctx, &entity.FindOptions{
		Email: cmd.Email,
	})
	if err != nil {
		return nil, err
	}

	if len(find) > 0 {
		return nil, errors.New("user already exists")
	}

	e := &entity.User{
		Name:     cmd.Name,
		Email:    cmd.Email,
		Status:   entity.UserStatus(cmd.Status),
		Password: cmd.Password,
	}

	fmt.Println("创建user", e)
	create, err := h.repo.Create(ctx, e)
	if err != nil {
		h.logger.Error("创建err：", err)
		return nil, err
	}

	return create, nil
}
