---
id: faster-velocity
title: ⚡ Make Velocity Faster
sidebar_position: 7
---

## Custom JSON Encoder/Decoder

Since Velocity v2.32.0, we have adopted `encoding/json` as the default JSON library for its stability and reliability. However, the standard library can be slower than some third-party alternatives. If you find the performance of `encoding/json` unsatisfactory, we suggest considering these libraries:

- [goccy/go-json](https://github.com/goccy/go-json)
- [bytedance/sonic](https://github.com/bytedance/sonic)
- [segmentio/encoding](https://github.com/segmentio/encoding)
- [minio/simdjson-go](https://github.com/minio/simdjson-go)

```go title="Example"
package main

import "go.khulnasoft.com/velocity"
import "github.com/goccy/go-json"

func main() {
    app := velocity.New(velocity.Config{
        JSONEncoder: json.Marshal,
        JSONDecoder: json.Unmarshal,
    })

    # ...
}
```

### References

- [Set custom JSON encoder for client](../client/rest.md#setjsonmarshal)
- [Set custom JSON decoder for client](../client/rest.md#setjsonunmarshal)
- [Set custom JSON encoder for application](../api/velocity.md#jsonencoder)
- [Set custom JSON decoder for application](../api/velocity.md#jsondecoder)
