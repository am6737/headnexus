package command

import (
	"context"
	"github.com/am6737/headnexus/app/user"
	"github.com/am6737/headnexus/domain/user/entity"
)

func (h *UserHandler) Update(ctx context.Context, cmd *user.UpdateUser) (*entity.User, error) {
	u, err := h.repo.Get(ctx, cmd.ID)
	if err != nil {
		return nil, err
	}

	e := &entity.User{
		ID:           u.ID,
		Name:         cmd.Name,
		Email:        cmd.Email,
		Token:        cmd.Token,
		Verification: cmd.Verification,
		Status:       entity.UserStatus(cmd.Status),
		Password:     cmd.Password,
	}

	err = h.repo.Update(ctx, e)
	if err != nil {
		return nil, err
	}
	return e, nil
}
