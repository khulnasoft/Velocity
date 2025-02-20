module github.com/khulnasoft/velocity/template/jet/v2

go 1.23

toolchain go1.23.4

require (
	github.com/CloudyKit/jet/v6 v6.2.0
	github.com/khulnasoft/velocity v0.0.0-20250220222828-6d269a9cedcc
	github.com/khulnasoft/velocity/template v0.0.0-00010101000000-000000000000
	github.com/khulnasoft/velocity/utils v0.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.10.0
)

require (
	github.com/CloudyKit/fastprinter v0.0.0-20200109182630-33d98a066a53 // indirect
	github.com/andybalholm/brotli v1.1.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/fxamacker/cbor/v2 v2.7.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/khulnasoft/schema v1.0.0 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/philhofer/fwd v1.1.3-0.20240916144458-20a13a1f6b7c // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/tinylib/msgp v1.2.5 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.59.0 // indirect
	github.com/x448/float16 v0.8.4 // indirect
	golang.org/x/crypto v0.33.0 // indirect
	golang.org/x/net v0.35.0 // indirect
	golang.org/x/sys v0.30.0 // indirect
	golang.org/x/text v0.22.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/khulnasoft/velocity/template => ../.

replace github.com/khulnasoft/velocity/utils => ../../utils

replace github.com/khulnasoft/velocity => ../../.
