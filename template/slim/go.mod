module github.com/khulnasoft/velocity/template/slim/v2

go 1.22

toolchain go1.23.4

require (
	github.com/khulnasoft/velocity/template v0.0.0-00010101000000-000000000000
	github.com/khulnasoft/velocity/utils v0.0.0-00010101000000-000000000000
	github.com/mattn/go-slim v0.0.4
	github.com/stretchr/testify v1.10.0
	github.com/valyala/bytebufferpool v1.0.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/khulnasoft/velocity/template => ../.
replace github.com/khulnasoft/velocity/utils => ../../utils
