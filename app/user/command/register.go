package command

import (
	"context"
	"fmt"
	"github.com/am6737/headnexus/app/user"
	"github.com/am6737/headnexus/domain/user/entity"
	pkgjwt "github.com/am6737/headnexus/pkg/jwt"
	pkgstring "github.com/am6737/headnexus/pkg/string"
	"github.com/dgrijalva/jwt-go"
	"time"
)

func (h *UserHandler) Register(ctx context.Context, cmd *user.CreateUser) (*entity.User, error) {
	find, err := h.repo.Find(ctx, &entity.FindOptions{
		Email: cmd.Email,
	})
	if err != nil {
		return nil, err
	}
	if len(find) > 0 {
		return nil, fmt.Errorf("用户已存在")
	}

	create, err := h.repo.Create(ctx, &entity.User{
		Name:     cmd.Name,
		Email:    cmd.Email,
		Password: pkgstring.Md5(cmd.Password),
		Status:   entity.Disable,
	})
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	code := pkgstring.GenerateRandomCode()
	verification, err := pkgjwt.GenerateTokenWithExpiryAndKey(jwt.MapClaims{
		"user_id": create.ID,
	}, 30*time.Second, []byte(code))
	if err != nil {
		return nil, err
	}

	create.Verification = verification
	err = h.repo.Update(ctx, create)
	if err != nil {
		return nil, err
	}

	if err := h.emailClient.SendEmail(cmd.Email, "验证码", code); err != nil {
		fmt.Println("4", err)
		return nil, err
	}
	return create, nil
}
