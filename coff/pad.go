package coff

type PaddedData struct {
	Data    Sizer
	Padding []byte
}

func pad(data Sizer) PaddedData {
	return PaddedData{
		Data:    data,
		Padding: make([]byte, -data.Size()&7),
	}
}
