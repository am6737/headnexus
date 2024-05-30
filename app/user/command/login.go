package command

import (
	"context"
	"errors"
	ctime "github.com/am6737/headnexus/common/time"
	"github.com/am6737/headnexus/domain/user/entity"
	pkgstring "github.com/am6737/headnexus/pkg/string"
	"github.com/dgrijalva/jwt-go"
)

func (h *UserHandler) Login(ctx context.Context, email string, password string) (string, *entity.User, error) {
	u, err := h.repo.Find(ctx, &entity.FindOptions{Email: email})
	if err != nil {
		return "", nil, err
	}

	if len(u) != 1 {
		return "", nil, errors.New("用户不存在或密码错误")
	}

	u1 := u[0]

	lastLoginAt := u1.LastLoginAt

	if pkgstring.Md5(password) != u1.Password {
		return "", nil, errors.New("用户不存在或密码错误")
	}

	if u1.Status != entity.Normal {
		return "", nil, errors.New("用户未激活")
	}

	token, err := h.jwtConfig.GenerateToken(jwt.MapClaims{
		"user_id": u1.ID,
	})
	if err != nil {
		return "", nil, err
	}

	u1.Token = token
	u1.LastLoginAt = ctime.CurrentTimestampMillis()

	err = h.repo.Update(ctx, u1)
	if err != nil {
		h.logger.Error(err)
		return "", nil, err
	}

	return token, &entity.User{
		Name:        u1.Name,
		Email:       u1.Email,
		LastLoginAt: lastLoginAt,
	}, nil
}
