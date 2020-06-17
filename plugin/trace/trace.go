package trace

import (
	"io"
	"time"

	log "github.com/micro/go-micro/v2/logger"
	span "github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
)

var close io.Closer

func initTrace(name string) {
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:            true,
			BufferFlushInterval: 1 * time.Second,
		},
	}
	var (
		tracer span.Tracer
		err    error
	)
	tracer, close, err = cfg.New(
		name,
		config.Logger(jaeger.StdLogger),
	)
	if err != nil {
		log.Error(err)
		return
	}

	span.SetGlobalTracer(tracer)
}

func Close() {
	close.Close()
}
