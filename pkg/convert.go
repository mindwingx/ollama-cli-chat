package pkg

const (
	megaByte string = "MB"
	gigaByte string = "GB"
)

type FormatBytes struct {
	bytes int64
	size  float64
	unit  string
}

func NewBytes() *FormatBytes {
	return &FormatBytes{}
}

func (b *FormatBytes) SetSize(bytes int64) *FormatBytes {
	size := float64(bytes) / (1024 * 1024)
	unit := megaByte

	if size > 1024 {
		size = size / 1024
		unit = gigaByte
	}

	{
		b.bytes = bytes
		b.size = size
		b.unit = unit
	}

	return b
}

func (b *FormatBytes) Bytes() int64 {
	return b.bytes
}

// Size auto-converted to MB or GB
func (b *FormatBytes) Size() float64 {
	return b.size
}

// Unit auto-set to "MB" or "GB"
func (b *FormatBytes) Unit() string {
	return b.unit
}

func (b *FormatBytes) MB() (size float64, unit string) {
	size = float64(b.bytes) / (1024 * 1024)
	unit = megaByte
	return
}

func (b *FormatBytes) GB() (size float64, unit string) {
	size = float64(b.bytes) / (1024 * 1024 * 1024)
	unit = gigaByte
	return
}
