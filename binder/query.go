package binder

import (
	"github.com/valyala/fasthttp"
	"go.khulnasoft.com/velocity/utils"
)

// QueryBinding is the query binder for query request body.
type QueryBinding struct {
	EnableSplitting bool
}

// Name returns the binding name.
func (*QueryBinding) Name() string {
	return "query"
}

// Bind parses the request query and returns the result.
func (b *QueryBinding) Bind(reqCtx *fasthttp.Request, out any) error {
	data := make(map[string][]string)
	var err error

	reqCtx.URI().QueryArgs().VisitAll(func(key, val []byte) {
		if err != nil {
			return
		}

		k := utils.UnsafeString(key)
		v := utils.UnsafeString(val)
		err = formatBindData(out, data, k, v, b.EnableSplitting, true)
	})

	if err != nil {
		return err
	}

	return parse(b.Name(), out, data)
}

// Reset resets the QueryBinding binder.
func (b *QueryBinding) Reset() {
	b.EnableSplitting = false
}
