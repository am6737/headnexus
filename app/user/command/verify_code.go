package command

import (
	"context"
	"fmt"
	"github.com/am6737/headnexus/domain/user/entity"
	pkgjwt "github.com/am6737/headnexus/pkg/jwt"
	"github.com/dgrijalva/jwt-go"
)

func (h *UserHandler) Verify(ctx context.Context, email, code string) error {
	find, err := h.repo.Find(ctx, &entity.FindOptions{
		Email: email,
	})
	if err != nil {
		return err
	}

	if len(find) == 0 {
		return fmt.Errorf("user not found")
	}

	if find[0].Status == entity.Normal {
		return fmt.Errorf("请勿重复激活")
	}
	c, err := pkgjwt.ParseTokenWithKey(find[0].Verification, []byte(code))
	if err != nil {
		return fmt.Errorf("激活失败：%v", err.Error())
	}

	claims, ok := c.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("验证码错误")
	}

	if claims["user_id"] == "" {
		return fmt.Errorf("验证码错误")
	}

	find[0].Status = entity.Normal
	err = h.repo.Update(ctx, find[0])
	if err != nil {
		return err
	}

	return nil
}
