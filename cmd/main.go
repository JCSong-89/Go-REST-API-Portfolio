package main

import (
	"Go-REST-API-Portfolio/config"
	"Go-REST-API-Portfolio/internal/db/postrgres"
	"Go-REST-API-Portfolio/internal/db/redis"
	storageMinio "Go-REST-API-Portfolio/internal/db/storage-db"
	"Go-REST-API-Portfolio/internal/logger"
	"Go-REST-API-Portfolio/internal/server"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	_jaegerCfg "github.com/uber/jaeger-client-go/config"
	_jaegerLog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"
	"log"
	"os"
)

func main() {
	log.Println("서버시작")
	configPath := config.GetConfigPath(os.Getenv("config"))

	v, err := config.SetAppProfile(configPath)
	if err != nil {
		log.Fatalf("컨피그 에러: %v", err)
	}

	cfg, err := config.ParseViper(v)
	if err != nil {
		log.Fatalf("컨피그 바이퍼 에러: %v", err)
	}

	/*
		Log Set and Init
	*/
	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()
	appLogger.Infof("AppVersion: %s, LogLevel: %s, Mode: %s, SSL: %v", cfg.Server.AppVersion, cfg.Logger.Level, cfg.Server.Mode, cfg.Server.SSL)

	/*
		SQL DB Setup
	*/
	psqlDB, err := postrgres.NewPsqlDB(cfg)
	if err != nil {
		appLogger.Fatalf("Postgresql init Error: %s", err)
	} else {
		appLogger.Infof("PostgresDB 연결완료, 상태: %#v", psqlDB.Stats())
	}
	defer psqlDB.Close()

	/*
		Redis Setup
	*/
	redisClient := redis.NewRedisClient(cfg)
	defer redisClient.Close()

	/*
		Jaeger is Log tracer From MSA Apps, It is made uber since 2015
		JaegerSampler is making Sample. If Some haven't traced info and just null in Request to Service A then Making Sample info for trace.
		opt Const is all trace same decide about Sample.
		StdLogger is using pkg log. That why name is Standard Logger
		NullFactory is a metrics factory that returns NullCounter, NullTimer, NullHistogram, and NullGauge
		Metrics Function Create Metrics an option and init Metrics in the tracer
		Setup Jaeger
	*/
	jaegerCfgInstance := _jaegerCfg.Configuration{
		ServiceName: cfg.Jaeger.ServiceName,
		Sampler: &_jaegerCfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &_jaegerCfg.ReporterConfig{
			LogSpans:           cfg.Jaeger.LogSpans,
			LocalAgentHostPort: cfg.Jaeger.Host,
		},
	}

	storageDB, err := storageMinio.NewStorageClient(cfg.Minio.Endpoint, cfg.Minio.MinioAccessKey, cfg.Minio.MinioSecretKey, cfg.Minio.UseSSL)

	if err != nil {
		appLogger.Errorf("storageDB 연결에러: %s", err)
	}
	appLogger.Info("storageDB 연결")

	tracer, closer, err := jaegerCfgInstance.NewTracer(
		_jaegerCfg.Logger(_jaegerLog.StdLogger),
		_jaegerCfg.Metrics(metrics.NullFactory),
	)
	if err != nil {
		appLogger.Fatalf("Jaeger 초기화 에러: %v", err)
	}
	appLogger.Info("Jaeger 연결 완료")

	/*
		SetGlobalTracer return GlobalTracer and it is singleton
	*/
	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	appLogger.Info("opentracing 연결 완료")

	s := server.NewServer(cfg, psqlDB, redisClient, appLogger, storageDB)
	if err = s.Run(); err != nil {
		log.Fatal(err)
	}
}