package ico

import (
	"image"
)

const BI_RGB = 0

type Dir struct {
	Reserved uint16
	Type     uint16
	Count    uint16
}

type DirEntryCommon struct {
	Width      byte   // Width, in pixels, of the image
	Height     byte   // Height, in pixels, of the image
	ColorCount byte   // Number of colors in image (0 if >=8bpp)
	Reserved   byte   // Reserved (must be 0)
	Planes     uint16 // Color Planes
	BitCount   uint16 // Bits per pixel
	BytesInRes uint32 // How many bytes in this resource?
}

type DirEntry struct {
	DirEntryCommon
	ImageOffset uint32
}

type BitMapInfoHeader struct {
	Size          uint32
	Width         int32
	Height        int32
	Planes        uint16
	BitCount      uint16
	Compression   uint32
	SizeImage     uint32
	XPelsPerMeter int32
	YPelsPerMeter int32
	ClrUsed       uint32
	ClrImportant  uint32
}

type RGBQuad struct {
	Blue     byte
	Green    byte
	Red      byte
	Reserved byte
}

type ICO struct {
	image.Image
}
