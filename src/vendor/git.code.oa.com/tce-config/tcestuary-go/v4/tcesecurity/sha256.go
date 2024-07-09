package tcesecurity

import (
	"crypto/sha256"
	"hash"
)

const (
	Sha256Algorithm = "sha256"
)

func init() {
	f := func() (Hash, error) {
		return &Sha256Hash{
			ctx: sha256.New(),
		}, nil
	}
	registerHashFunc(Sha256Algorithm, f)
}

type Sha256Hash struct {
	ctx hash.Hash
}

func (h *Sha256Hash) Update(data []byte) error {
	_, err := h.ctx.Write(data)
	return err
}

func (h *Sha256Hash) Digest() ([]byte, error) {
	return h.ctx.Sum(nil), nil
}
