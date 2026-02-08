package logger

var (
	sliceCap = 4
)

type OutDump struct {
	Dump []byte
}

func (d *OutDump) Write(p []byte) (n int, err error) {
	d.Dump = p
	return len(p), nil
}

func NewOutDump() *OutDump {
	return &OutDump{}
}

type OutMultiDump struct {
	Dumps [][]byte
}

func (d *OutMultiDump) Write(p []byte) (n int, err error) {
	d.Dumps = append(d.Dumps, p)
	return len(p), nil
}

func NewOutMultiDump() *OutMultiDump {
	return &OutMultiDump{
		Dumps: make([][]byte, 0, sliceCap),
	}
}
