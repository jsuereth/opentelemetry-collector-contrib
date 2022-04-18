module github.com/open-telemetry/opentelemetry-collector-contrib/pkg/translator/jaeger

go 1.17

require (
	github.com/jaegertracing/jaeger v1.26.0
	github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal v0.34.0
	github.com/stretchr/testify v1.7.1
	go.opentelemetry.io/collector/model v0.49.0
)

require (
	github.com/apache/thrift v0.14.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.opentelemetry.io/collector/pdata v0.49.0 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
)

replace github.com/open-telemetry/opentelemetry-collector-contrib/internal/coreinternal => ../../../internal/coreinternal
