package stringx

type Numeric interface {
	uint8 |
		uint16 |
		uint32 |
		uint64 |
		int8 |
		int16 |
		int32 |
		int64 |
		float32 |
		float64 |
		int |
		uint
}

type convertTo int

const (
	Uint8 = convertTo(iota)
	Uint16
	Uint32
	Uint64
	Int8
	Int16
	Int32
	Int64
	Float32
	Float64
	Int
	Uint
)

type charset int

const (
	UTF8 = charset(iota)
	GB18030
	HZGB2312
	GBK
)
