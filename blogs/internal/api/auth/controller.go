package auth

import (
	"blogs/internal/model/dto/request"
	"blogs/internal/repository"
	"blogs/internal/service"
	"blogs/pkg/captcha"
	appErrors "blogs/pkg/errors"
	"blogs/pkg/logger"
	"blogs/pkg/responses"
	"blogs/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Controller struct {
	authService service.AuthService
	redisRepo   repository.RedisRepository
}

func NewController(authService service.AuthService, redisRepo repository.RedisRepository) *Controller {
	return &Controller{authService: authService, redisRepo: redisRepo}
}

func (ctrl *Controller) Register(c *gin.Context) {
	var registerRequest request.AuthRequest
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		responses.BadRequest(c, "请求参数错误")
		logger.Warn("请求参数错误", zap.Error(err))
		return
	}
	if !utils.QqEmailRegex.MatchString(registerRequest.Email) {
		responses.BadRequest(c, "目前仅支持qq邮箱")
		return
	}
	key, err := ctrl.authService.ConsumePrivateKey(registerRequest.KeyID)
	if err != nil {
		responses.FromError(c, err, appErrors.CodeInternalError, "服务器内部出错")
		return
	}
	if err = ctrl.authService.Register(registerRequest, key); err != nil {
		responses.FromError(c, err, appErrors.CodeBadRequest, "注册失败")
		return
	}
	responses.Success(c, "注册成功")
}

func (ctrl *Controller) Login(c *gin.Context) {
	var loginRequest request.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		responses.BadRequest(c, "请求参数错误")
		logger.Warn("请求参数错误", zap.Error(err))
		return
	}
	// 校验验证码
	if !captcha.Verify(loginRequest.CaptchaID, loginRequest.Captcha) {
		responses.BadRequest(c, "验证码错误")
		return
	}
	key, err := ctrl.authService.ConsumePrivateKey(loginRequest.KeyID)
	if err != nil {
		responses.FromError(c, err, appErrors.CodeInternalError, "服务器内部出错")
		return
	}
	loginUser, err := ctrl.authService.Login(loginRequest, key)
	if err != nil {
		responses.FromError(c, err, appErrors.CodeBadRequest, "登录失败")
		return
	}
	responses.Success(c, loginUser)
}

func (ctrl *Controller) Logout(c *gin.Context) {
	userID, ok := utils.CurrentUserID(c)
	if !ok {
		return
	}
	if err := ctrl.authService.Logout(c.Request.Context(), userID); err != nil {
		responses.InternalError(c, "登出失败")
		return
	}
	responses.Success(c, "登出成功")
}

func (ctrl *Controller) GetCaptcha(c *gin.Context) {
	id, b64s, err := captcha.Generate()
	if err != nil {
		responses.InternalError(c, "生成图形验证码失败，请刷新重试")
		return
	}
	responses.Success(c, map[string]string{
		"captcha_id":  id,
		"captcha_val": b64s,
	})
}

func (ctrl *Controller) SendCode(c *gin.Context) {
	var sendCodeRequest request.SendCodeRequest
	if err := c.ShouldBindJSON(&sendCodeRequest); err != nil {
		responses.BadRequest(c, "请求参数错误")
		return
	}
	if !utils.QqEmailRegex.MatchString(sendCodeRequest.Email) {
		responses.BadRequest(c, "目前仅支持qq邮箱")
		return
	}
	if err := ctrl.authService.SendCode(sendCodeRequest.Email); err != nil {
		responses.FromError(c, err, appErrors.CodeBadRequest, "发送验证码失败")
		return
	}
	responses.Success(c, "验证码已发送")
}

func (ctrl *Controller) ForgotPassword(c *gin.Context) {
	var forgotPasswordRequest request.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&forgotPasswordRequest); err != nil {
		responses.BadRequest(c, "请求参数错误")
		return
	}
	key, err := ctrl.authService.ConsumePrivateKey(forgotPasswordRequest.KeyID)
	if err != nil {
		responses.FromError(c, err, appErrors.CodeInternalError, "服务器内部出错")
		return
	}
	if err = ctrl.authService.ResetPassword(forgotPasswordRequest, key); err != nil {
		responses.FromError(c, err, appErrors.CodeBadRequest, "密码重置失败")
		return
	}
	responses.Success(c, "密码重置成功")
}

func (ctrl *Controller) GetPubKeyHandler(c *gin.Context) {
	key, err := ctrl.authService.GenerateRSAKey()
	if err != nil {
		responses.FromError(c, err, appErrors.CodeInternalError, "密钥生成失败")
		return
	}
	responses.Success(c, key)
}
