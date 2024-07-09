package tcestuary

import (
	"encoding/base64"
	"testing"
)

var hasher THasher

func init() {
	SetConfigDirectory("./_example")
	hasher, _ = NewTHasher()
}

func checkFun(b *testing.B, hasher THasher, msg, dig string) {
	c, _ := hasher.New()
	c.Update([]byte(msg))
	ret, _ := c.Digest()
	retB64 := base64.StdEncoding.EncodeToString(ret)
	if retB64 != dig {
		b.Errorf("Hash got = %v, want %v", retB64, dig)
	}
}

func Benchmark_TSM_HashXXXX_test(b *testing.B) {
	msg := "xxxx"
	dig := "zvIyGEz5o57JgUZvVen+MeUnPt29E/Vhqgw2zIjFkOY="
	for i := 1; i < b.N; i++ {
		go checkFun(b, hasher, msg, dig)
	}
}

func Benchmark_TSM_HashYYYY_test(b *testing.B) {
	msg := "yyyy"
	dig := "5+0D2OzZVHJpsce5OgXdnfV7LtggfA0u/vYByuMUNDw="
	for i := 1; i < b.N; i++ {
		go checkFun(b, hasher, msg, dig)
	}
}
