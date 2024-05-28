package command

import (
	"context"
	"fmt"
	"time"

	"github.com/am6737/headnexus/domain/user/entity"
	pkgjwt "github.com/am6737/headnexus/pkg/jwt"
	pkgstring "github.com/am6737/headnexus/pkg/string"
	"github.com/dgrijalva/jwt-go"
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

	//todo 优化激活邮件链接的前缀
	if h.emailClient != nil {
		if err := h.emailClient.SendEmail(email, "激活账号", getEmailTemplate(fmt.Sprintf("%s://%s", "http", h.httpConfig.Addr), email, code)); err != nil {
			fmt.Println("4", err)
			return err
		}
	}

	return nil
}

func getEmailTemplate(host, email, code string) string {
	return fmt.Sprintf(`
   <p>感谢您注册我们的服务。请点击下面的链接激活您的账号：</p>
   <a href="%s/api/v1/users/verify-code?email=%s&code=%s" style="background-color: #007bff; color: #ffffff; padding: 10px 20px; border-radius: 5px; text-decoration: none;">激活账号</a>
   <p>如果链接无法点击，请复制并粘贴以下链接到您的浏览器地址栏：</p>
   <p>%s/api/v1/users/verify-code?email=%s&code=%s</p>
   <p>如果您没有尝试注册，请忽略此邮件。</p>
   <p>感谢您的配合。</p>
   <p>祝好</p>
`, host, email, code, host, email, code)
}
