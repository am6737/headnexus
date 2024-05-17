package command

import (
	"context"
	"errors"
	"github.com/am6737/headnexus/domain/user/entity"
	pkgstring "github.com/am6737/headnexus/pkg/string"
	"github.com/dgrijalva/jwt-go"
)

func (h *UserHandler) Login(ctx context.Context, email string, password string) (string, error) {
	u, err := h.repo.Find(ctx, &entity.FindOptions{Email: email})
	if err != nil {
		return "", err
	}

	if len(u) != 1 {
		return "", errors.New("用户不存在或密码错误")
	}

	u1 := u[0]

	if pkgstring.Md5(password) != u1.Password {
		return "", errors.New("用户不存在或密码错误")
	}

	if u1.Status != entity.Normal {
		return "", errors.New("用户未激活")
	}

	token, err := h.jwtConfig.GenerateToken(jwt.MapClaims{
		"user_id": u1.ID,
	})
	if err != nil {
		return "", err
	}

	u1.Token = token

	err = h.repo.Update(ctx, u1)
	if err != nil {
		h.logger.Error(err)
		return "", err
	}

	return token, nil
}
