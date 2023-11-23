module github.com/StephanHCB/go-backend-service-common

go 1.20

// enforce some security updates so our scanner doesn't complain about code that the
// go dependency resolution system will not actually include in the end result anyway
exclude github.com/prometheus/client_golang v1.11.0

require (
	github.com/StephanHCB/go-autumn-acorn-registry v0.3.1
	github.com/StephanHCB/go-autumn-config-api v0.2.1
	github.com/StephanHCB/go-autumn-config-env v0.2.2
	github.com/StephanHCB/go-autumn-logging v0.3.0
	github.com/StephanHCB/go-autumn-logging-zerolog v0.5.0
	github.com/StephanHCB/go-autumn-restclient v0.7.3
	github.com/StephanHCB/go-autumn-restclient-apm v0.1.0
	github.com/StephanHCB/go-autumn-restclient-prometheus v0.1.2
	github.com/StephanHCB/go-autumn-web-swagger-ui v0.2.3
	github.com/go-chi/chi/v5 v5.0.10
	github.com/go-http-utils/headers v0.0.0-20181008091004-fed159eddc2a
	github.com/golang-jwt/jwt/v4 v4.5.0
	github.com/prometheus/client_golang v1.17.0
	github.com/rs/zerolog v1.31.0
	github.com/stretchr/testify v1.8.4
	go.elastic.co/apm/module/apmchiv5/v2 v2.4.5
	go.elastic.co/apm/v2 v2.4.7
)

require (
	github.com/armon/go-radix v1.0.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/elastic/go-sysinfo v1.7.1 // indirect
	github.com/elastic/go-windows v1.0.1 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/joeshaw/multierror v0.0.0-20140124173710-69b34d4ec901 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.4.1-0.20230718164431-9a2bf3000d16 // indirect
	github.com/prometheus/common v0.44.0 // indirect
	github.com/prometheus/procfs v0.11.1 // indirect
	github.com/shurcooL/httpfs v0.0.0-20190707220628-8d4bc4ba7749 // indirect
	github.com/shurcooL/vfsgen v0.0.0-20200824052919-0d455de96546 // indirect
	github.com/tidwall/tinylru v1.2.1 // indirect
	go.elastic.co/apm/module/apmhttp/v2 v2.4.5 // indirect
	go.elastic.co/fastjson v1.1.0 // indirect
	golang.org/x/sys v0.12.0 // indirect
	google.golang.org/protobuf v1.31.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	howett.net/plist v1.0.0 // indirect
)
