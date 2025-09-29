package engine

import (
	"context"

	"github.com/KingTrack/gin-kit/kit/conf"
	datacenterregistry "github.com/KingTrack/gin-kit/kit/internal/datacenter/registry"
	loggerregistry "github.com/KingTrack/gin-kit/kit/internal/logger/registry"
	metricregistry "github.com/KingTrack/gin-kit/kit/internal/metric/registry"
	mysqlregistry "github.com/KingTrack/gin-kit/kit/internal/mysql/registry"
	redisregistry "github.com/KingTrack/gin-kit/kit/internal/redis/registry"
	tlscontext "github.com/KingTrack/gin-kit/kit/internal/tls/context"
	tracerregistry "github.com/KingTrack/gin-kit/kit/internal/tracer/registry"
	"github.com/pkg/errors"
)

type ResourceFunc func(ctx context.Context, config *conf.Config) error

func (e *Engine) initResource(config *conf.Namespace) error {

	if len(config.RootPath) == 0 {
		return errors.New("启动路径path配置为空")
	}

	data, err := config.Source.Load(config.RootPath)
	if err != nil {
		return err
	}

	for k, v := range data {
		var namespaceConfig conf.Config
		if err := config.Decoder.Decode(v, &namespaceConfig); err != nil {
			return err
		}

		ctx := tlscontext.NewBackground(tlscontext.Value(k))

		if k == e.config.Namespace.ToString() {
			for _, v := range e.globalResourceFuncs {
				fn := v
				if err := fn(ctx, &namespaceConfig); err != nil {
					return err
				}
			}
		}

		for _, v := range e.namespaceResourceFuncs {
			fn := v
			if err := fn(ctx, &namespaceConfig); err != nil {
				return err
			}
		}
	}

	return nil
}

func initMySQL(registry *mysqlregistry.Registry) ResourceFunc {
	return func(ctx context.Context, config *conf.Config) error {
		if len(config.MySQL) == 0 {
			return nil
		}

		if err := registry.Init(ctx, config.MySQL); err != nil {
			return errors.WithMessage(err, "mysql 初始化失败")
		}
		return nil
	}
}

func initRedis(registry *redisregistry.Registry) ResourceFunc {
	return func(ctx context.Context, config *conf.Config) error {
		if len(config.Redis) == 0 {
			return nil
		}

		if err := registry.Init(ctx, config.Redis); err != nil {
			return errors.WithMessage(err, "redis 初始化失败")
		}
		return nil
	}
}

func initMetric(registry *metricregistry.Registry) ResourceFunc {
	return func(ctx context.Context, config *conf.Config) error {
		config.Metric.ServiceName = config.Server.ServiceName
		config.Metric.Endpoint = config.Hostname
		if err := registry.Init(ctx, &config.Metric); err != nil {
			return errors.WithMessage(err, "metric 初始化失败")
		}
		return nil
	}
}

func initTracer(registry *tracerregistry.Registry) ResourceFunc {
	return func(ctx context.Context, config *conf.Config) error {
		config.Tracer.ServiceName = config.Server.ServiceName
		if err := registry.Init(ctx, &config.Tracer); err != nil {
			return errors.WithMessage(err, "tracer 初始化失败")
		}
		return nil
	}
}

func initLogger(registry *loggerregistry.Registry) ResourceFunc {
	return func(ctx context.Context, config *conf.Config) error {
		config.Tracer.ServiceName = config.Server.ServiceName
		if err := registry.Init(ctx, &config.Logger); err != nil {
			return errors.WithMessage(err, "logger 初始化失败")
		}
		return nil
	}
}

func initDateCenter(registry *datacenterregistry.Registry) ResourceFunc {
	return func(ctx context.Context, config *conf.Config) error {
		if err := registry.Init(ctx, &config.Datacenter); err != nil {
			return errors.WithMessage(err, "datacenter 初始化失败")
		}
		return nil
	}
}
