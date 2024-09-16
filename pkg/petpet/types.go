package petpet

type Config struct {
	Width    int  // Recommend 128
	Height   int  // Recommend 128
	Delay    int  // Recommend 2-10
	Disposal byte // Recommend 0x02
}

var DefaultConfig = Config{
	Width:    128,
	Height:   128,
	Delay:    4,
	Disposal: 0x02,
}
