package user

import (
	"blogs/internal/model/dto/request"
	"blogs/internal/repository"
	"blogs/internal/service"
	appErrors "blogs/pkg/errors"
	"blogs/pkg/logger"
	"blogs/pkg/responses"
	"blogs/pkg/utils"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Controller struct {
	userService service.UserService
	redisRepo   repository.RedisRepository
}

func NewController(userService service.UserService, redisRepo repository.RedisRepository) *Controller {
	return &Controller{userService: userService, redisRepo: redisRepo}
}

func (ctrl *Controller) GetInfo(c *gin.Context) {
	userID, ok := utils.CurrentUserID(c)
	if !ok {
		return
	}
	user, err := ctrl.userService.GetByID(userID)
	if err != nil {
		responses.FromError(c, err, appErrors.CodeNotFound, "用户不存在")
		return
	}
	responses.Success(c, user)
}

func (ctrl *Controller) GetAuthor(c *gin.Context) {
	user, err := ctrl.userService.GetAuthor()
	if err != nil {
		responses.FromError(c, err, appErrors.CodeNotFound, "用户不存在")
		return
	}
	responses.Success(c, user)
}

func (ctrl *Controller) Update(c *gin.Context) {
	userID, ok := utils.CurrentUserID(c)
	if !ok {
		return
	}

	var updateReq request.UpdateUserRequest
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		responses.BadRequest(c, "请求参数错误")
		return
	}

	key, err := ctrl.resolvePrivateKey(updateReq.Password, updateReq.KeyID)
	if err != nil {
		responses.FromError(c, err, appErrors.CodeInternalError, "获取私钥失败")
		return
	}

	if err := ctrl.userService.Update(userID, updateReq, key); err != nil {
		responses.FromError(c, err, appErrors.CodeBadRequest, "更新失败")
		return
	}
	responses.Success(c, "更新成功")
}

func (ctrl *Controller) UploadAvatar(c *gin.Context) {
	userID, ok := utils.CurrentUserID(c)
	if !ok {
		return
	}

	file, err := c.FormFile("avatar")
	if err != nil {
		responses.BadRequest(c, "请选择上传文件")
		return
	}
	if err := utils.ValidateImage(file); err != nil {
		responses.BadRequest(c, err.Error())
		return
	}

	uploadDir := "./uploads/avatars"
	if err = utils.EnsureDir(uploadDir); err != nil {
		logger.Error("创建头像目录失败", zap.Error(err))
		responses.InternalError(c, "保存文件失败")
		return
	}

	user, err := ctrl.userService.GetByID(userID)
	if err != nil {
		responses.FromError(c, err, appErrors.CodeNotFound, "用户不存在")
		return
	}
	ctrl.removeOldAvatar(user.Avatar)

	ext := strings.ToLower(filepath.Ext(file.Filename))
	filename := fmt.Sprintf("%s%s", user.Username, ext)
	dst := filepath.Join(uploadDir, filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		logger.Error("保存头像失败", zap.Error(err))
		responses.InternalError(c, "保存文件失败")
		return
	}

	avatarPath := "/uploads/avatars/" + filename
	if err := ctrl.userService.UpdateAvatar(userID, avatarPath); err != nil {
		responses.FromError(c, err, appErrors.CodeInternalError, "更新头像失败")
		return
	}
	responses.Success(c, avatarPath)
}

func (ctrl *Controller) List(c *gin.Context) {
	var req request.UserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		responses.BadRequest(c, "参数错误")
		return
	}
	users, total, err := ctrl.userService.List(req)
	if err != nil {
		responses.FromError(c, err, appErrors.CodeInternalError, "获取列表失败")
		return
	}
	responses.Success(c, responses.NewPageResponse(users, total, req.Page, req.PageSize))
}

func (ctrl *Controller) Create(c *gin.Context) {
	var req request.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "请求参数错误")
		return
	}

	key, err := ctrl.userService.ConsumePrivateKey(req.KeyID)
	if err != nil {
		responses.FromError(c, err, appErrors.CodeInternalError, "获取私钥失败")
		return
	}

	if err := ctrl.userService.Create(req, key); err != nil {
		responses.FromError(c, err, appErrors.CodeBadRequest, "创建失败")
		return
	}
	responses.Success(c, "创建成功")
}

func (ctrl *Controller) AdminUpdate(c *gin.Context) {
	id, ok := utils.ParamInt(c, "id")
	if !ok {
		return
	}

	var updateReq request.UpdateUserRequest
	if err := c.ShouldBindJSON(&updateReq); err != nil {
		responses.BadRequest(c, "请求参数错误")
		return
	}
	key, err := ctrl.resolvePrivateKey(updateReq.Password, updateReq.KeyID)
	if err != nil {
		responses.FromError(c, err, appErrors.CodeInternalError, "获取私钥失败")
		return
	}
	if err := ctrl.userService.Update(id, updateReq, key); err != nil {
		responses.FromError(c, err, appErrors.CodeBadRequest, "更新失败")
		return
	}
	responses.Success(c, "更新成功")
}

func (ctrl *Controller) Delete(c *gin.Context) {
	id, ok := utils.ParamInt(c, "id")
	if !ok {
		return
	}
	if err := ctrl.userService.Delete(id); err != nil {
		responses.FromError(c, err, appErrors.CodeInternalError, "删除失败")
		return
	}
	responses.Success(c, "删除成功")
}

func (ctrl *Controller) GetStats(c *gin.Context) {
	res, err := ctrl.userService.GetStats()
	if err != nil {
		responses.FromError(c, err, appErrors.CodeInternalError, "获取数据失败")
	}
	responses.Success(c, res)
}

func (ctrl *Controller) PostQr(c *gin.Context) {
	userID, ok := utils.CurrentUserID(c)
	if !ok {
		return
	}

	file, err := c.FormFile("qr")
	if err != nil {
		responses.BadRequest(c, "上传文件为空")
		return
	}
	if err := utils.ValidateImage(file); err != nil {
		responses.BadRequest(c, err.Error())
		return
	}

	uploadDir := filepath.Join("uploads", "qr")
	if err := utils.EnsureDir(uploadDir); err != nil {
		responses.InternalError(c, "上传失败")
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	filename := utils.BuildUploadFilename(userID, ext)
	dst := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, dst); err != nil {
		responses.InternalError(c, "上传失败")
		return
	}

	url := "/uploads/" + "qr" + "/" + filename
	responses.Success(c, url)
}

func (ctrl *Controller) resolvePrivateKey(password, keyID string) (string, error) {
	if password == "" {
		return "", nil
	}
	return ctrl.userService.ConsumePrivateKey(keyID)
}

func (ctrl *Controller) removeOldAvatar(oldAvatar string) {
	if oldAvatar == "" || strings.Contains(oldAvatar, "default.png") {
		return
	}

	oldAvatarFile := "." + oldAvatar
	if err := os.Remove(oldAvatarFile); err != nil {
		logger.Warn("删除旧头像失败", zap.String("path", oldAvatarFile), zap.Error(err))
	}
}
