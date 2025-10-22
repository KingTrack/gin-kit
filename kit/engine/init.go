package engine

import (
	"os"

	"github.com/KingTrack/gin-kit/kit/conf"
	contextregistry "github.com/KingTrack/gin-kit/kit/internal/context/registry"
	datacenterregistry "github.com/KingTrack/gin-kit/kit/internal/datacenter/registry"
	httpclientregistry "github.com/KingTrack/gin-kit/kit/internal/httpclient/registry"
	kafkaconsumerregistry "github.com/KingTrack/gin-kit/kit/internal/kafka/consumer/registry"
	kafkaproducerregistry "github.com/KingTrack/gin-kit/kit/internal/kafka/producer/registry"
	loggerregistry "github.com/KingTrack/gin-kit/kit/internal/logger/registry"
	metricregistry "github.com/KingTrack/gin-kit/kit/internal/metric/registry"
	mysqlregistry "github.com/KingTrack/gin-kit/kit/internal/mysql/registry"
	redisregistry "github.com/KingTrack/gin-kit/kit/internal/redis/registry"
	tracerregistry "github.com/KingTrack/gin-kit/kit/internal/tracer/registry"
	"github.com/KingTrack/gin-kit/kit/plugin/decoder"
	"github.com/pkg/errors"
)

type Engine struct {
	path                   string
	config                 *conf.Config
	httpClientRegistry     *httpclientregistry.Registry
	mysqlRegistry          *mysqlregistry.Registry
	redisRegistry          *redisregistry.Registry
	kafkaProducerRegistry  *kafkaproducerregistry.Registry
	kafkaConsumerRegistry  *kafkaconsumerregistry.Registry
	globalResourceFuncs    []ResourceOption
	namespaceResourceFuncs []ResourceOption
	tracerRegistry         *tracerregistry.Registry
	loggerRegistry         *loggerregistry.Registry
	metricRegistry         *metricregistry.Registry
	contextRegistry        *contextregistry.Registry
	datacenterRegistry     *datacenterregistry.Registry
}

func NewDefault() *Engine {
	return &Engine{
		tracerRegistry:     tracerregistry.New(),
		loggerRegistry:     loggerregistry.New(),
		metricRegistry:     metricregistry.New(),
		contextRegistry:    contextregistry.New(),
		datacenterRegistry: datacenterregistry.New(),
	}
}

func New(path string) *Engine {

	engine := &Engine{
		path:                  path,
		httpClientRegistry:    httpclientregistry.New(),
		mysqlRegistry:         mysqlregistry.New(),
		redisRegistry:         redisregistry.New(),
		kafkaProducerRegistry: kafkaproducerregistry.New(),
		kafkaConsumerRegistry: kafkaconsumerregistry.New(),
		tracerRegistry:        tracerregistry.New(),
		loggerRegistry:        loggerregistry.New(),
		metricRegistry:        metricregistry.New(),
		contextRegistry:       contextregistry.New(),
		datacenterRegistry:    datacenterregistry.New(),
	}

	engine.globalResourceFuncs = []ResourceOption{
		initLogger(engine),
		initMetric(engine),
		initTracer(engine),
		initDateCenter(engine),
	}

	engine.namespaceResourceFuncs = []ResourceOption{
		initHTTPClient(engine),
		initMySQL(engine),
		initRedis(engine),
		initKafkaProducer(engine),
	}

	return engine
}

func (e *Engine) Init(opts ...Option) error {
	hostname, err := os.Hostname()
	if err != nil {
		return errors.WithMessage(err, "hostname 获取失败")
	}
	e.config.Hostname = hostname

	if err := e.initResource(&conf.Namespace{
		RootPath: e.path,
		Source:   conf.NewSource(),
		Decoder:  decoder.TOMLDecoder{},
	}); err != nil {
		return err
	}

	for _, opt := range opts {
		if err := opt(e); err != nil {
			return err
		}
	}

	return nil
}
