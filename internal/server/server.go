package server

import (
	"Go-REST-API-Portfolio/config"
	"Go-REST-API-Portfolio/internal/logger"
	"Go-REST-API-Portfolio/internal/middlewares"
	"Go-REST-API-Portfolio/internal/prometheus"
	"context"
	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/minio/minio-go"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	CERT_FILE        = "ssl/Server.crt"
	KEY_FILE         = "ss/Server.pem"
	MAX_HEADER_BYTES = 1 << 20 // Header Size Limit
	CTX_TIMEOUT      = 5
)

type Server struct {
	echo        *echo.Echo
	cfg         *config.Config
	db          *sqlx.DB
	redisClient *redis.Client
	logger      logger.Logger
	minioClient *minio.Client
}

func (s *Server) Run() error {
	/*
		HTTTPS
	*/
	if s.cfg.Server.SSL {
		if err := s.Init(s.echo); err != nil {
			return err
		}

		s.echo.Server.ReadTimeout = time.Second * s.cfg.Server.ReadTimeout
		s.echo.Server.WriteTimeout = time.Second * s.cfg.Server.WriteTimeout

		go func() {
			s.logger.Infof("서버 listening PORT: %v", s.cfg.Server.Port)
			s.echo.Server.MaxHeaderBytes = MAX_HEADER_BYTES
			if err := s.echo.StartTLS(s.cfg.Server.Port, CERT_FILE, KEY_FILE); err != nil {
				s.logger.Fatalf("TLS Server 시작 에러: %v", err)
			}
		}()

		/*
			os.Signal get Unix SIignal and Notifiy function get notify when signal from os signal
		*/
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

		<-quit

		/*Life time Set up */
		ctx, shutdown := context.WithTimeout(context.Background(), CTX_TIMEOUT*time.Second)
		defer shutdown()

		s.logger.Info("서버 셧다운")
		return s.echo.Server.Shutdown(ctx)
	}

	/*
		HTTP
	*/
	server := &http.Server{
		Addr:           s.cfg.Server.Port,
		ReadTimeout:    time.Second * s.cfg.Server.ReadTimeout,
		WriteTimeout:   time.Second * s.cfg.Server.WriteTimeout,
		MaxHeaderBytes: MAX_HEADER_BYTES,
	}

	go func() {
		s.logger.Infof("서버 listening PORT: %v", s.cfg.Server.Port)
		if err := s.echo.StartServer(server); err != nil {
			s.logger.Fatalf("Error starting Server: ", err)
		}
	}()

	if err := s.Init(s.echo); err != nil {
		return err
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), CTX_TIMEOUT*time.Second)
	defer shutdown()

	s.logger.Info("서버 셧다운")
	return s.echo.Server.Shutdown(ctx)
}

func (s *Server) Init(e *echo.Echo) error {
	mt, err := prometheus.CreateMetrics(s.cfg.Metrics.URL, s.cfg.Metrics.ServiceName)
	if err != nil {
		s.logger.Errorf("CreateMetrics Error: %s", err)
	}
	s.logger.Infof("메트릭스 URL: %s, 서비스Name: %s", s.cfg.Metrics.URL, s.cfg.Metrics.ServiceName)

	/*
		HTTP -> HTTPS Redirect
	*/
	if s.cfg.Server.SSL {
		e.Pre(middleware.HTTPSRedirect())
	}

	md := middlewares.NewMiddleware(s.cfg, []string{"*"}, s.logger)
	e.Use(md.ReqLogMiddlewareInit)

	// CORS SetUp
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderXRequestID, "X-CSRF-Token"},
	}))

	/*
		DisablePrintStack is trace stack and printing
		RecoverWithConfig return Recover that is make recove from painc
	*/
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         1 << 10, // 1kb
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))
	/*
		return X-request-ID
	*/
	e.Use(middleware.RequestID())

	/*
		metrics Info Injection middleware when reqeust come
	*/
	e.Use(md.MetricsMiddlewareInit(mt))

	/*
		압축방식 설정
	*/
	e.Use(middleware.Gzip())

	/*
		Body Size limit and Secure middleware in echo
	*/
	e.Use(middleware.Secure())
	e.Use(middleware.BodyLimit("2M"))

	/*
		URI Set up
	*/
	appVersion := e.Group("api/v1")
	health := appVersion.Group("/health") // health check

	health.GET("", func(ctx echo.Context) error {
		s.logger.Info("서버가동상태조회")
		return ctx.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})

	return nil
}

func NewServer(c *config.Config, db *sqlx.DB, redisClient *redis.Client, logger logger.Logger, minioClient *minio.Client) *Server {
	return &Server{echo.New(), c, db, redisClient, logger, minioClient}
}