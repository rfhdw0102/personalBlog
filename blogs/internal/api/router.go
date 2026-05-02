package api

import (
	"blogs/internal/api/about"
	"blogs/internal/api/article"
	"blogs/internal/api/auth"
	"blogs/internal/api/category"
	"blogs/internal/api/comment"
	"blogs/internal/api/notification"
	"blogs/internal/api/tag"
	"blogs/internal/api/user"
	"blogs/internal/middleware"
	"blogs/internal/repository"
	"blogs/internal/service"

	"github.com/gin-gonic/gin"
)

type Router struct {
	authCtrl         *auth.Controller
	userCtrl         *user.Controller
	articleCtrl      *article.Controller
	notificationCtrl *notification.Controller
	commentCtrl      *comment.Controller
	categoryCtrl     *category.Controller
	tagCtrl          *tag.Controller
	aboutCtrl        *about.Controller
}

func NewRouter(
	authService service.AuthService,
	userService service.UserService,
	articleService service.ArticleService,
	notificationService service.NotificationService,
	commentService service.CommentService,
	tagService service.TagService,
	categoryService service.CategoryService,
	aboutService service.AboutService,
	redisRepo repository.RedisRepository,
) *Router {
	return &Router{
		authCtrl:         auth.NewController(authService, redisRepo),
		userCtrl:         user.NewController(userService, redisRepo),
		articleCtrl:      article.NewController(articleService, redisRepo),
		notificationCtrl: notification.NewController(notificationService, redisRepo),
		commentCtrl:      comment.NewController(redisRepo, commentService),
		categoryCtrl:     category.NewController(redisRepo, categoryService),
		tagCtrl:          tag.NewController(redisRepo, tagService),
		aboutCtrl:        about.NewController(redisRepo, aboutService),
	}
}

func (r *Router) Setup(engine *gin.Engine) {
	engine.Use(middleware.Recovery())
	engine.Use(middleware.Logger())
	engine.Use(middleware.CORS())

	// 静态资源
	engine.Static("/uploads", "./uploads")

	// 健康检查
	engine.GET("/api/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "blogs API is running",
		})
	})

	// API 路由组
	v := engine.Group("/api")
	{
		r.authCtrl.RegisterRoutes(v)
		r.userCtrl.RegisterRoutes(v)
		r.articleCtrl.RegisterRoutes(v)
		r.notificationCtrl.RegisterRoutes(v)
		r.commentCtrl.RegisterRoutes(v)
		r.categoryCtrl.RegisterRoutes(v)
		r.tagCtrl.RegisterRoutes(v)
		r.aboutCtrl.RegisterRoutes(v)
	}

}
