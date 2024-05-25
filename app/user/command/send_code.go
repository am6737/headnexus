package command

import (
	"context"
	"fmt"
	"github.com/am6737/headnexus/domain/user/entity"
	pkgjwt "github.com/am6737/headnexus/pkg/jwt"
	pkgstring "github.com/am6737/headnexus/pkg/string"
	"github.com/dgrijalva/jwt-go"
	"time"
)

func (h *UserHandler) SendCode(ctx context.Context, email string) error {
	find, err := h.repo.Find(ctx, &entity.FindOptions{
		Email: email,
	})
	if err != nil {
		return err
	}
	if len(find) == 0 {
		return fmt.Errorf("用户不存在")
	}

	code := pkgstring.GenerateRandomCode()
	verification, err := pkgjwt.GenerateTokenWithExpiryAndKey(jwt.MapClaims{
		"user_id": find[0].ID,
	}, 30*time.Minute, []byte(code))
	if err != nil {
		return err
	}

	find[0].Verification = verification
	err = h.repo.Update(ctx, find[0])
	if err != nil {
		return err
	}

	if err := h.emailClient.SendEmail(email, "验证码", code); err != nil {
		fmt.Println("4", err)
		return err
	}
	return nil
}
