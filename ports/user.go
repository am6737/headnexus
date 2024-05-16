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
	userInfo, err := h.app.User.Queries.Handler.Get(c, &user.GetUser{
		ID: userID,
	})
	if err != nil {
		return
	}

	if userInfo.Password != password {
		http.FailedResponse(c, "旧密码错误")
		return
	}

	userInfo.Password = pkgstring.Md5(*req.NewPassword)

	_, err = h.app.User.Commands.Handler.Update(c, &user.UpdateUser{
		ID:           userInfo.ID,
		Password:     userInfo.Password,
		Token:        userInfo.Token,
		Name:         userInfo.Name,
		Email:        userInfo.Email,
		Status:       uint(userInfo.Status),
		Verification: userInfo.Verification,
	})
	if err != nil {
		http.FailedResponse(c, "修改失败")
		return
	}

	http.SuccessResponse(c, "修改成功", gin.H{})
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

	email, err := req.Email.MarshalJSON()
	if err != nil {
		http.FailedResponse(c, "参数错误")
		return
	}

	// 去除双引号
	newemail := strings.Trim(string(email), "\"")

	if err := h.app.User.Commands.Handler.Verify(c, newemail, *req.Code); err != nil {
		http.FailedResponse(c, err.Error())
		return
	}

	http.SuccessResponse(c, "激活成功", nil)
}

func (h *HttpHandler) SendCode(c *gin.Context) {
	req := &v1.VerifyCodeRequest{}
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
	if err := h.app.User.Commands.Handler.SendCode(c, newemail); err != nil {
		http.FailedResponse(c, err.Error())
		return
	}

	http.SuccessResponse(c, "发送成功", nil)
}