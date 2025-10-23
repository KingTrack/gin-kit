package engine

import (
	"context"

	runtimedatacenter "github.com/KingTrack/gin-kit/kit/runtime/datacenter"

	"github.com/KingTrack/gin-kit/kit/conf"
	tlscontext "github.com/KingTrack/gin-kit/kit/internal/tls/context"
	runtimelogger "github.com/KingTrack/gin-kit/kit/runtime/logger"
	"github.com/pkg/errors"
	"go.uber.org/multierr"
)

type ResourceOption func(ctx context.Context, config *conf.Config) error

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

func initMetric(e *Engine) ResourceOption {
	return func(ctx context.Context, config *conf.Config) error {
		config.Metric.ServiceName = config.Server.ServiceName
		config.Metric.Endpoint = config.Hostname
		if err := e.metricRegistry.Init(ctx, &config.Metric); err != nil {
			return errors.WithMessage(err, "metric 初始化失败")
		}
		return nil
	}
}

func initTracer(e *Engine) ResourceOption {
	return func(ctx context.Context, config *conf.Config) error {
		config.Tracer.ServiceName = config.Server.ServiceName
		if err := e.tracerRegistry.Init(ctx, &config.Tracer); err != nil {
			return errors.WithMessage(err, "tracer 初始化失败")
		}
		return nil
	}
}

func initLogger(e *Engine) ResourceOption {
	return func(ctx context.Context, config *conf.Config) error {
		if err := e.loggerRegistry.Init(ctx, &config.Logger); err != nil {
			return errors.WithMessage(err, "logger 初始化失败")
		}
		runtimelogger.Set(e.loggerRegistry)
		return nil
	}
}

func initDateCenter(e *Engine) ResourceOption {
	return func(ctx context.Context, config *conf.Config) error {
		if err := e.datacenterRegistry.Init(ctx, &config.Datacenter); err != nil {
			return errors.WithMessage(err, "datacenter 初始化失败")
		}
		runtimedatacenter.Set(e.datacenterRegistry)
		return nil
	}
}

func initHTTPClient(e *Engine) ResourceOption {
	return func(ctx context.Context, config *conf.Config) error {
		if len(config.HTTPClient) == 0 {
			return nil
		}

		var errs error
		for _, v := range config.HTTPClient {
			clientConfig := v
			err := e.httpClientRegistry.Add(ctx, &clientConfig)
			if err != nil {
				errs = multierr.Append(errs, err)
				continue
			}

			err = e.datacenterRegistry.AddHTTPClient(ctx, &clientConfig)
			if err != nil {
				errs = multierr.Append(errs, err)
				continue
			}
		}
		if errs != nil {
			return errors.WithMessage(errs, "http client 初始化失败")
		}
		return nil
	}
}

func initMySQL(e *Engine) ResourceOption {
	return func(ctx context.Context, config *conf.Config) error {
		if len(config.MySQL) == 0 {
			return nil
		}

		if err := e.mysqlRegistry.Init(ctx, config.MySQL); err != nil {
			return errors.WithMessage(err, "mysql 初始化失败")
		}
		return nil
	}
}

func initRedis(e *Engine) ResourceOption {
	return func(ctx context.Context, config *conf.Config) error {
		if len(config.Redis) == 0 {
			return nil
		}

		if err := e.redisRegistry.Init(ctx, config.Redis); err != nil {
			return errors.WithMessage(err, "redis 初始化失败")
		}
		return nil
	}
}

func initKafkaProducer(e *Engine) ResourceOption {
	return func(ctx context.Context, config *conf.Config) error {
		if len(config.KafkaProducer) == 0 {
			return nil
		}

		if err := e.kafkaProducerRegistry.Init(ctx, config.KafkaProducer); err != nil {
			return errors.WithMessage(err, "kafka producer 初始化失败")
		}
		return nil
	}
}

func initKafkaConsumer(e *Engine) ResourceOption {
	return func(ctx context.Context, config *conf.Config) error {
		if len(config.KafkaConsumer) == 0 {
			return nil
		}

		if err := e.kafkaConsumerRegistry.Init(ctx, config.KafkaConsumer); err != nil {
			return errors.WithMessage(err, "kafka consumer 初始化失败")
		}
		return nil
	}
}
