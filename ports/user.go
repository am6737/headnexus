package ports

import (
	"fmt"
	"strings"

	v1 "github.com/am6737/headnexus/api/http/v1"
	"github.com/am6737/headnexus/app/user"
	ctime "github.com/am6737/headnexus/common/time"
	"github.com/am6737/headnexus/domain/user/entity"
	"github.com/am6737/headnexus/pkg/http"
	pkgstring "github.com/am6737/headnexus/pkg/string"
	"github.com/gin-gonic/gin"
)

func (h *HttpHandler) GetUserInfo(c *gin.Context) {
	uid := c.Value("user_id").(string)
	if uid == "" {
		http.FailedResponse(c, "user not found")
		return
	}
	userInfo, err := h.app.User.Queries.Handler.Get(c, &user.GetUser{ID: uid})
	if err != nil {
		http.FailedResponse(c, err.Error())
		return
	}

	http.SuccessResponse(c, "获取用户信息成功", getUserInfoToResponse(userInfo))
}

func getUserInfoToResponse(info *entity.User) *v1.UserInfo {
	return &v1.UserInfo{
		Email:       info.Email,
		Id:          info.ID,
		LastLoginAt: ctime.FormatTimeSince(info.LastLoginAt),
	}
}

func (h *HttpHandler) ChangePassword(c *gin.Context) {
	req := &v1.ChangePasswordRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		http.FailedResponse(c, "参数错误")
		return
	}

	fmt.Println("req.NewPassword => ", req.NewPassword)
	fmt.Println("req.ConfirmPassword => ", req.ConfirmPassword)

	if req.NewPassword != req.ConfirmPassword {
		http.FailedResponse(c, "密码不一致")
		return
	}

	// 从上下文中获取用户信息
	id, ok := c.Get("user_id")
	if !ok {
		http.FailedResponse(c, "User ID not found in context")
		return
	}

	userID := id.(string)

	password := pkgstring.Md5(req.OldPassword)
	userInfo, err := h.app.User.Queries.Handler.Get(c, &user.GetUser{
		ID: userID,
	})
	if err != nil {
		return
	}

	fmt.Println(" userInfo.Password => ", userInfo.Password)
	fmt.Println("password => ", password)

	if userInfo.Password != password {
		http.FailedResponse(c, "旧密码错误")
		return
	}

	userInfo.Password = pkgstring.Md5(req.NewPassword)

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

	token, info, err := h.app.User.Commands.Handler.Login(c, newemail, *req.Password)
	if err != nil {
		http.FailedResponse(c, err.Error())
		return
	}

	http.SuccessResponse(c, "登录成功", gin.H{
		"token": token,
		"user_info": &v1.UserInfo{
			Email:       info.Email,
			Id:          info.ID,
			LastLoginAt: ctime.FormatTimeSince(info.LastLoginAt),
		},
	})
}

func (h *HttpHandler) RegisterUser(c *gin.Context) {
	req := &v1.RegisterUserJSONRequestBody{}
	if err := c.ShouldBindJSON(req); err != nil {
		http.FailedResponse(c, "参数错误")
		return
	}

	if req.Password != req.ConfirmPassword {
		http.FailedResponse(c, "密码不一致")
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
		Password: req.Password,
	})
	if err != nil {
		http.FailedResponse(c, err.Error())
		return
	}

	http.SuccessResponse(c, "注册成功", nil)
}

func (h *HttpHandler) VerifyCode(c *gin.Context, params v1.VerifyCodeParams) {

	email, err := params.Email.MarshalJSON()
	if err != nil {
		http.FailedResponse(c, "参数错误")
		return
	}

	// 去除双引号
	newemail := strings.Trim(string(email), "\"")

	if err := h.app.User.Commands.Handler.Verify(c, newemail, params.Code); err != nil {
		http.FailedResponse(c, err.Error())
		return
	}

	http.SuccessResponse(c, "激活成功", nil)
}

func (h *HttpHandler) SendCode(c *gin.Context) {
	req := &v1.SendCodeRequest{}
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
