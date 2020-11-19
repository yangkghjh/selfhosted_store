package project

// Decoder from template to template applications
type Decoder func([]byte) (*Application, error)

// Encoder from applications to template
type Encoder func(*Application) ([]byte, error)

var (
	decoders map[string]Decoder
	encoders map[string]Encoder
)

func init() {
	decoders = make(map[string]Decoder)
	encoders = make(map[string]Encoder)
}

// RegisterDecoder register decoder
func RegisterDecoder(name string, f Decoder) {
	decoders[name] = f
}

// RegisterEncoder register encoder
func RegisterEncoder(name string, f Encoder) {
	encoders[name] = f
}

// LoadDecoder load decoder
func LoadDecoder(name string) Decoder {
	if f, ok := decoders[name]; ok {
		return f
	}

	return nil
}

// LoadEncoder load encoder
func LoadEncoder(name string) Encoder {
	if f, ok := encoders[name]; ok {
		return f
	}

	return nil
}
