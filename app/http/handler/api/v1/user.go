package v1

import (
	"github.com/gin-gonic/gin"
	"go-chat/app/entity"
	"go-chat/app/http/request"
	"go-chat/app/http/response"
	"go-chat/app/model"
	"go-chat/app/pkg/auth"
	"go-chat/app/service"
)

type User struct {
	service    *service.UserService
	smsService *service.SmsService
}

func NewUserHandler(
	userService *service.UserService,
	smsService *service.SmsService,
) *User {
	return &User{
		service:    userService,
		smsService: smsService,
	}
}

// Detail 个人用户信息
func (u *User) Detail(ctx *gin.Context) {
	user, _ := u.service.Dao().FindById(auth.GetAuthUserID(ctx))

	response.Success(ctx, gin.H{
		"detail": user,
	})
}

// ChangeDetail 修改个人用户信息
func (u *User) ChangeDetail(ctx *gin.Context) {
	params := &request.ChangeDetailRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	_, _ = u.service.Dao().Update(&model.User{ID: auth.GetAuthUserID(ctx)}, map[string]interface{}{
		"nickname": params.Nickname,
		"avatar":   params.Avatar,
		"gender":   params.Gender,
		"motto":    params.Profile,
	})

	response.Success(ctx, gin.H{}, "个人信息修改成功！")
}

// ChangePassword 修改密码接口
func (u *User) ChangePassword(ctx *gin.Context) {
	params := &request.ChangePasswordRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	if err := u.service.UpdatePassword(auth.GetAuthUserID(ctx), params.OldPassword, params.NewPassword); err != nil {
		response.BusinessError(ctx, "密码修改失败！")
		return
	}

	response.Success(ctx, gin.H{}, "密码修改成功！")
}

// ChangeMobile 修改手机号接口
func (u *User) ChangeMobile(ctx *gin.Context) {
	params := &request.ChangeMobileRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	if !u.smsService.CheckSmsCode(ctx.Request.Context(), entity.SmsChangeAccountChannel, params.Mobile, params.SmsCode) {
		response.BusinessError(ctx, "短信验证码填写错误！")
		return
	}

	user, _ := u.service.Dao().FindById(auth.GetAuthUserID(ctx))

	if user.Mobile != params.Mobile {
		response.BusinessError(ctx, "手机号与原手机号一致无需修改！")
		return
	}

	if !auth.Compare(user.Password, params.Password) {
		response.BusinessError(ctx, "账号密码填写错误！")
		return
	}

	_, err := u.service.Dao().Update(&model.User{ID: user.ID}, map[string]interface{}{
		"mobile": params.Mobile,
	})

	if err != nil {
		response.BusinessError(ctx, "手机号修改失败！")
		return
	}

	response.Success(ctx, gin.H{}, "手机号修改成功！")
}

// ChangeEmail 修改邮箱接口
func (u *User) ChangeEmail(ctx *gin.Context) {
	params := &request.ChangeEmailRequest{}
	if err := ctx.ShouldBind(params); err != nil {
		response.InvalidParams(ctx, err)
		return
	}

	// todo 1.验证邮件激活码是否正确
}
