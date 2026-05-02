package article

import (
	"blogs/internal/model/dto/request"
	"blogs/internal/model/entity"
	"blogs/internal/repository"
	"blogs/internal/service"
	appErrors "blogs/pkg/errors"
	"blogs/pkg/logger"
	"blogs/pkg/responses"
	"blogs/pkg/utils"
	"errors"
	"path/filepath"
	"strconv"
	"strings"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	article   service.ArticleService
	redisRepo repository.RedisRepository
}

func NewController(article service.ArticleService, redisRepo repository.RedisRepository) *Controller {
	return &Controller{article: article, redisRepo: redisRepo}
}

func (ctrl *Controller) List(c *gin.Context) {
	var req request.ArticleListQuery
	if err := c.ShouldBindQuery(&req); err != nil {
		responses.BadRequest(c, "参数错误")
		return
	}

	articles, total, err := ctrl.article.ListWithContent(req)
	if err != nil {
		responses.FromError(c, err, appErrors.CodeInternalError, "获取列表失败")
		return
	}

	responses.Success(c, responses.NewPageResponse(articles, total, req.Page, req.PageSize))
}

func (ctrl *Controller) Get(c *gin.Context) {
	id, ok := utils.ParamInt(c, "id")
	if !ok {
		return
	}

	article, err := ctrl.article.GetByID(id)
	if err != nil {
		responses.FromError(c, err, appErrors.CodeNotFound, "文章不存在")
		return
	}

	responses.Success(c, article)
}

func (ctrl *Controller) IncrementView(c *gin.Context) {
	id, ok := utils.ParamInt(c, "id")
	if !ok {
		return
	}

	if err := ctrl.article.IncrementViewCount(id); err != nil {
		logger.Warn("浏览量增加失败", zap.Int("文章id:", id), zap.Error(err))
	}

	responses.Success(c, nil)
}

func (ctrl *Controller) LikeArticle(c *gin.Context) {
	id, ok := utils.ParamInt(c, "id")
	if !ok {
		return
	}
	uid, ok := utils.CurrentUserID(c)
	if !ok {
		return
	}
	if err := ctrl.article.LikeArticle(id, uid); err != nil {
		responses.FromError(c, err, appErrors.CodeInternalError, "点赞失败")
		return
	}

	responses.Success(c, nil)
}

func (ctrl *Controller) UnlikeArticle(c *gin.Context) {
	id, ok := utils.ParamInt(c, "id")
	if !ok {
		return
	}
	uid, ok := utils.CurrentUserID(c)
	if !ok {
		return
	}

	if err := ctrl.article.UnlikeArticle(id, uid); err != nil {
		responses.FromError(c, err, appErrors.CodeInternalError, "取消点赞失败")
		return
	}

	responses.Success(c, nil)
}

func (ctrl *Controller) Create(c *gin.Context) {
	var req request.ArticleSubmit
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "参数错误")
		return
	}
	uid, ok := utils.CurrentUserID(c)
	if !ok {
		return
	}
	article := ctrl.newArticleFromSubmit(req, uid)

	if err := ctrl.article.Create(article, req.TagIDs, req.TagNames); err != nil {
		responses.FromError(c, err, appErrors.CodeInternalError, "提交文章失败")
		return
	}

	responses.Success(c, article)
}

func (ctrl *Controller) Update(c *gin.Context) {
	var req request.ArticleEdit
	if err := c.ShouldBindJSON(&req); err != nil {
		responses.BadRequest(c, "参数错误")
		return
	}

	articleInfo, err := ctrl.article.GetByID(req.ID)
	if err != nil {
		responses.FromError(c, err, appErrors.CodeNotFound, "文章不存在")
		return
	}
	article := ctrl.newArticleFromEdit(req, articleInfo.UserID)

	if err := ctrl.article.Update(article, req.TagIDs, req.TagNames); err != nil {
		responses.FromError(c, err, appErrors.CodeInternalError, "编辑文章失败")
		return
	}

	responses.Success(c, article)
}

func (ctrl *Controller) Delete(c *gin.Context) {
	id, ok := utils.ParamInt(c, "id")
	if !ok {
		return
	}
	_, err := ctrl.article.GetByID(id)

	if err != nil {
		responses.FromError(c, err, appErrors.CodeNotFound, "文章不存在")
		return
	}

	if err := ctrl.article.Delete(id); err != nil {
		responses.FromError(c, err, appErrors.CodeInternalError, "删除文章失败")
		return
	}

	responses.Success(c, nil)
}

func (ctrl *Controller) UploadCoverImage(c *gin.Context) {
	url, err := ctrl.uploadImage(c, "coverImage")
	if err != nil {
		logger.Error("上传封面图片出错", zap.Error(err))
		responses.InternalError(c, err.Error())
		return
	}
	responses.Success(c, url)
}

func (ctrl *Controller) UploadContentImage(c *gin.Context) {
	url, err := ctrl.uploadImage(c, "contentImage")
	if err != nil {
		logger.Error("上传内容图片出错", zap.Error(err))
		responses.InternalError(c, err.Error())
		return
	}
	responses.Success(c, url)

}

func (ctrl *Controller) uploadImage(c *gin.Context, imageType string) (string, error) {
	userID, ok := utils.CurrentUserID(c)
	if !ok {
		return "", errors.New("用户未登录")
	}

	file, err := c.FormFile(imageType)
	if err != nil {
		return "", errors.New("上传文件为空")
	}
	if err := utils.ValidateImage(file); err != nil {
		return "", err
	}

	uploadDir := filepath.Join("uploads", imageType)
	if err := utils.EnsureDir(uploadDir); err != nil {
		return "", err
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	filename := utils.BuildUploadFilename(userID, ext)
	dst := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, dst); err != nil {
		return "", err
	}

	url := "/uploads/" + imageType + "/" + filename
	return url, nil
}

func (ctrl *Controller) CheckLikeStatus(c *gin.Context) {
	id, ok := utils.ParamInt(c, "id")
	if !ok {
		return
	}
	userId := c.Query("uid")
	uid, err := strconv.Atoi(userId)
	if err != nil {
		responses.FromError(c, err, appErrors.CodeInternalError, "服务器内部出错")
		return
	}
	isLiked, err := ctrl.article.IsLiked(id, uid)
	if err != nil {
		responses.FromError(c, err, appErrors.CodeInternalError, "查询是否点赞失败")
		return
	}
	responses.Success(c, gin.H{"is_liked": isLiked})
}

func (ctrl *Controller) GetAdjacent(c *gin.Context) {
	id, ok := utils.ParamInt(c, "id")
	if !ok {
		return
	}

	sort, _ := strconv.Atoi(c.DefaultQuery("sort", "0"))

	result, err := ctrl.article.GetAdjacent(id, sort)
	if err != nil {
		responses.FromError(c, err, appErrors.CodeInternalError, "获取相邻文章失败")
		return
	}

	responses.Success(c, result)
}

func (ctrl *Controller) newArticleFromSubmit(req request.ArticleSubmit, userID int) *entity.Article {
	return &entity.Article{
		Title:      req.Title,
		Summary:    req.Summary,
		Content:    req.Content,
		CoverImage: req.CoverImage,
		CategoryID: req.CategoryID,
		Status:     req.Status,
		UserID:     userID,
	}
}

func (ctrl *Controller) newArticleFromEdit(req request.ArticleEdit, userID int) *entity.Article {
	return &entity.Article{
		Base: entity.Base{
			ID: req.ID,
		},
		Title:      req.Title,
		Summary:    req.Summary,
		Content:    req.Content,
		CoverImage: req.CoverImage,
		CategoryID: req.CategoryID,
		Status:     req.Status,
		UserID:     userID,
	}
}
