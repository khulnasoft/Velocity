module github.com/khulnasoft/velocity/template/handlebars/v2

go 1.20

require (
	github.com/khulnasoft/velocity v2.52.6
	github.com/khulnasoft/velocity/template v0.0.0-00010101000000-000000000000
	github.com/khulnasoft/velocity/utils v0.0.0-00010101000000-000000000000
	github.com/mailgun/raymond/v2 v2.0.48
	github.com/stretchr/testify v1.10.0
)

require (
	github.com/andybalholm/brotli v1.1.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.51.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/khulnasoft/velocity/template => ../.
replace github.com/khulnasoft/velocity/utils => ../../utils
