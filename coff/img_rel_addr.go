package coff

// Values reverse-engineered from windres output; names from teh Internets.
// Teh googlies Internets don't seem to have much to say about the AMD64 one,
// unfortunately :/ but it works...
const (
	_IMAGE_REL_AMD64_ADDR32NB = 0x03
	_IMAGE_REL_I386_DIR32NB   = 0x07
	_IMAGE_REL_ARM64_ADDR32NB = 0x02
	_IMAGE_REL_ARM_ADDR32NB   = 0x02
)