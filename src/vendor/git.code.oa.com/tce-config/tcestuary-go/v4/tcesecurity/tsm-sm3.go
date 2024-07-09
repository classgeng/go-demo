package tcesecurity

import (
	"fmt"

	sm "git.code.oa.com/tce-config/tcestuary-go/v4/tcesecurity/tencentsm"
)

const (
	Tsm3Algorithm = "tsm-sm3"
)

func init() {
	f := func() (Hash, error) {
		var ctx sm.SM3_ctx_t
		sm.SM3Init(&ctx)
		return &TSM3Hash{
			ctx: &ctx,
		}, nil
	}
	registerHashFunc(Tsm3Algorithm, f)
}

type TSM3Hash struct {
	ctx *sm.SM3_ctx_t
}

func (h *TSM3Hash) Update(data []byte) error {
	if code := sm.SM3Update(h.ctx, data, len(data)); code != 0 {
		return fmt.Errorf("sm3 update failed, code: %d", code)
	}
	return nil
}

func (h *TSM3Hash) Digest() ([]byte, error) {
	out := make([]byte, 32)
	if code := sm.SM3Final(h.ctx, out); code != 0 {
		return out, fmt.Errorf("sm3 final failed, code: %d", code)
	}
	return out, nil
}
