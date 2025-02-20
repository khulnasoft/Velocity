module github.com/khulnasoft/velocity/template/ace/v2

go 1.17

require github.com/yosssi/ace v0.0.5

require (
	github.com/khulnasoft/velocity/template v0.0.0-00010101000000-000000000000
	github.com/khulnasoft/velocity/utils v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.10.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/khulnasoft/velocity/template => ../.
replace github.com/khulnasoft/velocity/utils => ../../utils
