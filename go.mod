module github.com/StephanHCB/go-backend-service-common

go 1.17

// enforce some security updates so our scanner doesn't complain about code that the
// go dependency resolution system will not actually include in the end result anyway
exclude github.com/prometheus/client_golang v1.11.0

require (
	github.com/StephanHCB/go-autumn-acorn-registry v0.2.0
	github.com/StephanHCB/go-autumn-config-api v0.2.1
	github.com/StephanHCB/go-autumn-config-env v0.1.0
	github.com/StephanHCB/go-autumn-logging v0.3.0
	github.com/StephanHCB/go-autumn-logging-zerolog v0.3.1
	github.com/StephanHCB/go-autumn-web-swagger-ui v0.2.3
	github.com/go-chi/chi/v5 v5.0.7
	github.com/go-http-utils/headers v0.0.0-20181008091004-fed159eddc2a
	github.com/prometheus/client_golang v1.13.0
	github.com/rs/zerolog v1.28.0
	github.com/stretchr/testify v1.8.0
)

require (
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/mattn/go-colorable v0.1.12 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	github.com/shurcooL/httpfs v0.0.0-20190707220628-8d4bc4ba7749 // indirect
	github.com/shurcooL/vfsgen v0.0.0-20200824052919-0d455de96546 // indirect
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
