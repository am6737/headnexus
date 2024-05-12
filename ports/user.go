package ports

import (
	v1 "github.com/am6737/headnexus/api/http/v1"
	"github.com/am6737/headnexus/app/user"
	"github.com/am6737/headnexus/pkg/http"
	pkgstring "github.com/am6737/headnexus/pkg/string"
	"github.com/gin-gonic/gin"
	"strings"
)

func (h *HttpHandler) ChangePassword(c *gin.Context) {
	req := &v1.ChangePasswordRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		http.FailedResponse(c, "参数错误")
		return
	}

	// 从上下文中获取用户信息
	id, ok := c.Get("user_id")
	if !ok {
		http.FailedResponse(c, "User ID not found in context")
		return
	}

	userID := id.(string)

	password := pkgstring.Md5(*req.OldPassword)
	u, err := h.app.User.Queries.Handler.Get(c, &user.GetUser{
		ID: userID,
	})
	if err != nil {
		http.FailedResponse(c, "用户不存在")
		return
	}

	if u.Password != password {
		http.FailedResponse(c, "修改失败")
		return
	}
	_, err = h.app.User.Commands.Handler.Update(c, &user.UpdateUser{
		ID:       userID,
		Password: *req.NewPassword,
	})
	if err != nil {
		http.FailedResponse(c, "修改失败")
		return
	}

}

func (h *HttpHandler) LoginUser(c *gin.Context) {
	req := &v1.LoginRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		http.FailedResponse(c, "参数错误")
		return
	}

	email, err := req.Email.MarshalJSON()
	if err != nil {
		http.FailedResponse(c, "参数错误")
		return
	}

	// 去除双引号
	newemail := strings.Trim(string(email), "\"")

	token, err := h.app.User.Commands.Handler.Login(c, newemail, *req.Password)
	if err != nil {
		http.FailedResponse(c, err.Error())
		return
	}

	http.SuccessResponse(c, "登录成功", gin.H{
		"token": token,
	})
}

func (h *HttpHandler) RegisterUser(c *gin.Context) {
	req := &v1.RegisterRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		http.FailedResponse(c, "参数错误")
		return
	}

	email, err := req.Email.MarshalJSON()
	if err != nil {
		http.FailedResponse(c, "参数错误")
		return
	}

	// 去除双引号
	newemail := strings.Trim(string(email), "\"")

	name := newemail
	if req.Name != nil {
		name = *req.Name
	}

	_, err = h.app.User.Commands.Handler.Register(c, &user.CreateUser{
		Name:     name,
		Email:    newemail,
		Password: *req.Password,
	})
	if err != nil {
		http.FailedResponse(c, err.Error())
		return
	}

	http.SuccessResponse(c, "注册成功", nil)
}

func (h *HttpHandler) VerifyCode(c *gin.Context) {
	req := &v1.VerifyCodeRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		http.FailedResponse(c, "参数错误")
		return
	}
}
