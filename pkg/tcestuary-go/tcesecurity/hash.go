package tcesecurity

// Hash
type Hash interface {
	Update([]byte) error
	Digest() ([]byte, error)
}

// 散列算法
type HashFunc func() (Hash, error)

// 支持的散列算法
var SupportHashFunc = make(map[string]HashFunc, 0)

// 注册支持的散列算法
func registerHashFunc(k string, f HashFunc) {
	SupportHashFunc[k] = f
}
