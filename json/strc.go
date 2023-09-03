package json

type Ret struct {
	b []byte
}

func (r *Ret) String() string {
	return string(r.b)
}

func (r *Ret) Bytes() []byte {
	return r.b
}
