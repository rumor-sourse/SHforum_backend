package routes

import (
	_ "SHforum_backend/docs"
	"SHforum_backend/internal/settings"
	"SHforum_backend/middlewares"
	"SHforum_backend/pkg/logger"
	"SHforum_backend/routes/admin"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"time"
)

func SetUp(mode string) *gin.Engine {
	// 设置成发布模式
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	r.Use(logger.GinLogger(),
		logger.GinRecovery(true),
		middlewares.RateLimitMiddleware(time.Microsecond*time.Duration(200), 20000),
		otelgin.Middleware(settings.Conf.Name),
		func(c *gin.Context) {
			c.Header("Trace-Id", trace.SpanFromContext(c.Request.Context()).SpanContext().TraceID().String())
		})

	// 注册swagger路由
	r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))
	api := r.Group("/api/v1")

	admin.RegisterUserRouter(api)
	admin.RegisterCommunityRouter(api)
	admin.RegisterPostRouter(api)
	admin.RegisterCommentController(api)

	pprof.Register(r)
	// 测试
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	// 404
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "404 not found",
		})
	})
	return r
}
