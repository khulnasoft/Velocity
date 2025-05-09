// ⚡️ Velocity is an Express inspired web framework written in Go with ☕️
// 🤖 Github Repository: https://github.com/khulnasoft/velocity
// 📌 API Documentation: https://docs.khulnasoft.com

package utils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_ToUpper(t *testing.T) {
	t.Parallel()
	require.Equal(t, "/MY/NAME/IS/:PARAM/*", ToUpper("/my/name/is/:param/*"))
}

const (
	largeStr = "/RePos/KhulnaSoft/VeloCity/iSsues/187643/CoMmEnts/RePos/KhulnaSoft/VeloCity/iSsues/CoMmEnts"
	upperStr = "/REPOS/KHULNASOFT/VELOCITY/ISSUES/187643/COMMENTS/REPOS/KHULNASOFT/VELOCITY/ISSUES/COMMENTS"
	lowerStr = "/repos/khulnasoft/velocity/issues/187643/comments/repos/khulnasoft/velocity/issues/comments"
)

func Benchmark_ToUpper(b *testing.B) {
	var res string
	b.Run("velocity", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = ToUpper(largeStr)
		}
		require.Equal(b, upperStr, res)
	})
	b.Run("IfToUpper-Upper", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = IfToUpper(upperStr)
		}
		require.Equal(b, upperStr, res)
	})
	b.Run("IfToUpper-Mixed", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = IfToUpper(largeStr)
		}
		require.Equal(b, upperStr, res)
	})
	b.Run("default", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = strings.ToUpper(largeStr)
		}
		require.Equal(b, upperStr, res)
	})
}

func Test_ToLower(t *testing.T) {
	t.Parallel()
	require.Equal(t, "/my/name/is/:param/*", ToLower("/MY/NAME/IS/:PARAM/*"))
	require.Equal(t, "/my1/name/is/:param/*", ToLower("/MY1/NAME/IS/:PARAM/*"))
	require.Equal(t, "/my2/name/is/:param/*", ToLower("/MY2/NAME/IS/:PARAM/*"))
	require.Equal(t, "/my3/name/is/:param/*", ToLower("/MY3/NAME/IS/:PARAM/*"))
	require.Equal(t, "/my4/name/is/:param/*", ToLower("/MY4/NAME/IS/:PARAM/*"))
}

func Benchmark_ToLower(b *testing.B) {
	var res string
	b.Run("velocity", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = ToLower(largeStr)
		}
		require.Equal(b, lowerStr, res)
	})
	b.Run("IfToLower-Lower", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = IfToLower(lowerStr)
		}
		require.Equal(b, lowerStr, res)
	})
	b.Run("IfToLower-Mixed", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = IfToLower(largeStr)
		}
		require.Equal(b, lowerStr, res)
	})
	b.Run("default", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = strings.ToLower(largeStr)
		}
		require.Equal(b, lowerStr, res)
	})
}

func Test_IfToUpper(t *testing.T) {
	t.Parallel()
	require.Equal(t, "MYNAMEISPARAM", IfToUpper("MYNAMEISPARAM")) // already uppercase
	require.Equal(t, "MYNAMEISPARAM", IfToUpper("mynameisparam")) // lowercase to uppercase
	require.Equal(t, "MYNAMEISPARAM", IfToUpper("MyNameIsParam")) // mixed case
}

func Test_IfToLower(t *testing.T) {
	t.Parallel()
	require.Equal(t, "mynameisparam", IfToLower("mynameisparam"))                   // already lowercase
	require.Equal(t, "mynameisparam", IfToLower("myNameIsParam"))                   // mixed case
	require.Equal(t, "https://khulnasoft.com", IfToLower("https://khulnasoft.com")) // Origin Header Type URL
	require.Equal(t, "mynameisparam", IfToLower("MYNAMEISPARAM"))                   // uppercase to lowercase
}

// Benchmark_IfToLower_HeadersOrigin benchmarks the IfToLower function with an origin header type URL.
// These headers are typically lowercase, so the function should return the input string without modification.
func Benchmark_IfToToLower_HeadersOrigin(b *testing.B) {
	var res string
	b.Run("velocity", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = ToLower("https://khulnasoft.com")
		}
		require.Equal(b, "https://khulnasoft.com", res)
	})
	b.Run("IfToLower-Lower", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = IfToLower("https://khulnasoft.com")
		}
		require.Equal(b, "https://khulnasoft.com", res)
	})
	b.Run("default", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			res = strings.ToLower("https://khulnasoft.com")
		}
		require.Equal(b, "https://khulnasoft.com", res)
	})
}
