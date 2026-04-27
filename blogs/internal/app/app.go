package app

import (
	"blogs/internal/api"
	"blogs/internal/model/entity"
	"blogs/internal/repository"
	"blogs/internal/service"
	"blogs/pkg/config"
	"blogs/pkg/database"
	"blogs/pkg/job"
	"blogs/pkg/logger"
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type App struct {
	cfg                  *config.Config
	mysqlDB              *gorm.DB
	redisDB              *redis.Client
	router               *api.Router
	server               *http.Server
	notificationConsumer *job.NotificationConsumer
	syncJob              *job.ArticleSync
}

// NewApp 创建应用实例
func NewApp() *App {
	return &App{}
}

// Initialize 初始化应用
func (a *App) Initialize() error {
	// 加载配置
	if err := a.initConfig(); err != nil {
		return err
	}

	// 初始化日志
	if err := a.initLogger(); err != nil {
		return err
	}

	// 初始化数据库
	if err := a.initDatabase(); err != nil {
		return err
	}

	// 初始化依赖
	a.initDependencies()

	// 初始化路由
	a.initRouter()

	// 初始化服务器
	a.initServer()

	return nil
}

// initConfig 记载配置
func (a *App) initConfig() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("加载配置失败：%w", err)
	}
	a.cfg = cfg
	return err
}

// initLogger 初始化日志
func (a *App) initLogger() error {
	if err := logger.Init(&a.cfg.Log); err != nil {
		return fmt.Errorf("日志初始化失败: %w", err)
	}
	logger.Info("配置加载成功")
	logger.Info("=========================================")
	return nil
}

// initDatabase 初始化数据库
func (a *App) initDatabase() error {
	mysqlDB, err := database.InitMysql(&a.cfg.Database.MySQL)
	if err != nil {
		return fmt.Errorf("MySQL 初始化失败: %w", err)
	}
	a.mysqlDB = mysqlDB
	// 自动迁移数据库表
	logger.Info("开始数据库迁移...")
	if err := a.mysqlDB.AutoMigrate(
		entity.User{},
		entity.Article{},
		entity.Tag{},
		entity.Comment{},
		entity.Category{},
		entity.Like{},
		entity.Notification{},
	); err != nil {
		logger.Warn("数据库迁移警告", zap.Error(err))
	} else {
		logger.Info("数据库迁移成功")
	}
	redisDB, err := database.InitRedis(&a.cfg.Database.Redis)
	if err != nil {
		return fmt.Errorf("Redis 初始化失败: %w", err)
	}
	a.redisDB = redisDB
	return nil
}

// initDependencies 初始化依赖
func (a *App) initDependencies() {
	// 创建Repo
	redisRepo := repository.NewRedisRepository(a.redisDB)
	userRepo := repository.NewUserRepository(a.mysqlDB)
	articleRepo := repository.NewArticleRepository(a.mysqlDB, redisRepo)
	notificationRepo := repository.NewNotificationRepository(a.mysqlDB)
	commentRepo := repository.NewCommentRepository(a.mysqlDB, redisRepo)
	categoryRepo := repository.NewCategoryRepository(a.mysqlDB)
	tagRepo := repository.NewTagRepository(a.mysqlDB)
	// 创建Service
	authService := service.NewAuthService(userRepo, redisRepo)
	userService := service.NewUserService(userRepo, redisRepo)
	articleService := service.NewArticleService(articleRepo, categoryRepo, tagRepo, commentRepo, redisRepo)
	notificationService := service.NewNotificationService(notificationRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	commentService := service.NewCommentService(commentRepo, articleRepo)
	tagService := service.NewTagService(tagRepo)
	// 创建Consumer
	a.notificationConsumer = job.NewNotificationConsumer(redisRepo, notificationRepo)
	// 创建定时同步器
	a.syncJob = job.NewArticleSync(redisRepo, articleRepo)

	// 创建router
	a.router = api.NewRouter(authService, userService, articleService, notificationService, commentService, tagService, categoryService, redisRepo)
}

// initRouter 初始化路由
func (a *App) initRouter() {
	// 设置 Gin 模式
	gin.SetMode(a.cfg.App.Mode)
}

// initServer 初始化 HTTP 服务器
func (a *App) initServer() {
	engine := gin.New()

	// 注册路由
	a.router.Setup(engine)

	// 创建 HTTP 服务器
	a.server = &http.Server{
		Addr:              fmt.Sprintf(":%d", a.cfg.App.Port),
		Handler:           engine,
		ReadHeaderTimeout: 60 * time.Second,  // 增加读取头超时，防范慢连接攻击
		IdleTimeout:       120 * time.Second, // 增加空闲连接超时，释放资源
		MaxHeaderBytes:    1 << 20,           // 1 MB
	}
}

// Run 运行应用
func (a *App) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 启动通知消费者
	go a.notificationConsumer.Start(ctx)

	// 启动同步器
	go a.startSyncTask(ctx)

	// 启动 HTTP 服务器
	go func() {
		logger.Info("HTTP 服务器启动",
			zap.String("addr", a.server.Addr),
			zap.String("mode", a.cfg.App.Mode),
		)
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("HTTP 服务器启动失败", zap.Error(err))
		}
	}()

	// 优雅关闭
	a.gracefulShutdown(cancel)
}

// gracefulShutdown 优雅关闭
func (a *App) gracefulShutdown(cancel context.CancelFunc) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("正在关闭服务器...")

	// 停止消费者
	cancel()

	ctx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	// 关闭 HTTP 服务器
	if err := a.server.Shutdown(ctx); err != nil {
		logger.Error("服务器关闭失败", zap.Error(err))
	}

	// 关闭数据库连接
	_ = database.CloseMySQL()
	_ = database.CloseRedis()

	// 同步日志
	_ = logger.Sync()

	logger.Info("服务器已关闭")
	logger.Info("=========================================")
}

func (a *App) startSyncTask(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				a.syncJob.SyncChangedToMySQL(ctx)
			case <-ctx.Done():
				return
			}
		}
	}()
}
