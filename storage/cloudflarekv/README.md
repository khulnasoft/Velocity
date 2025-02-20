---
id: cloudflarekv
title: Cloudflare KV
---

![Release](https://img.shields.io/github/v/tag/khulnasoft/storage?filter=cloudflarekv*)
[![Discord](https://img.shields.io/discord/704680098577514527?style=flat&label=%F0%9F%92%AC%20discord&color=00ACD7)](https://khulnasoft.com/discord)
![Test](https://img.shields.io/github/actions/workflow/status/khulnasoft/storage/test-cloudflarekv.yml?label=Tests)
![Security](https://img.shields.io/github/actions/workflow/status/khulnasoft/storage/gosec.yml?label=Security)
![Linter](https://img.shields.io/github/actions/workflow/status/khulnasoft/storage/linter.yml?label=Linter)

A Cloudflare KV storage driver using [cloudflare/cloudflare-go](https://github.com/cloudflare/cloudflare-go).

**Note: Requires Go 1.21 and above**

### Table of Contents

- [Signatures](#signatures)
- [Installation](#installation)
- [Examples](#examples)
- [Config](#config)
- [Default Config](#default-config)

### Signatures

```go
func New(config ...Config) Storage
func (s *Storage) Get(key string) ([]byte, error)
func (s *Storage) Set(key string, val []byte, exp time.Duration) error
func (s *Storage) Delete(key string) error
func (s *Storage) Reset() error
func (s *Storage) Close() error
func (s *Storage) Conn() *cloudflare.API
```

### Installation

```bash
go mod init github.com/<user>/<repo>
```

And then install the Cloudflare KV implementation:

```bash
go get go.khulnasoft.com/velocity/storage/cloudflarekv
```

### Examples

Import the storage package.

```go
import "go.khulnasoft.com/velocity/storage/cloudflarekv"
```

You can use the following methods to create storage. The Key must be an API Token generated with at least `Account.Workers KV Storage` permission. Check the [Create API Token](https://developers.cloudflare.com/fundamentals/api/get-started/create-token/) documentation to generate one.

```go
// Initialize default config
store := cloudflarekv.New()

store := cloudflarekv.New(cloudflarekv.Config{
    Key: "",
    Email: "",
    AccountID: "velocity",
    NamespaceID: "velocity",
    Reset: false,
})

```

### Config

```go
type Config struct {

	// Cloudflare Auth Token
	//
	// Optional. Default is ""
	Key string

	// Cloudflare Email
	//
	// Optional. Default is ""
	Email string

	// Account id
	//
	// Optional. Default is "velocity"
	AccountID string

	// Namespace id
	//
	// Optional. Default is "velocity"
	NamespaceID string

	// Reset clears any existing keys in existing Table
	//
	// Optional. Default is false
	Reset bool
}
```

### Default Config

```go
var ConfigDefault = Config{
	Key:         "",
	Email:       "",
	AccountID:   "velocity",
	NamespaceID: "velocity",
	Reset:       false,
}
```
