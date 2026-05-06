package cnu

const (
	ExpectedSignature = "CNU_File_Archive_3.0"
	HeaderSize        = 0x248
)

type SignedHeader struct {
	Tag1        byte
	Length1     uint16
	VersionData uint16
	Tag2        byte
	Length2     uint16
	HeaderSize  uint16
}

// RawHeader exactly matches the binary layout of the archive's 584-byte header.
type RawHeader struct {
	Signature1    [20]byte
	Padding1      [28]byte
	NumFiles      uint32
	ContentOffset uint32
	Signature2    [20]byte
	Padding2      [508]byte
}

// RawFileTableEntry represents a single 20-byte record in the file table.
type RawFileTableEntry struct {
	F1             uint32 // 7 = normal file, 6 = raw/special
	F2             uint32
	Size           uint32
	Timestamp      uint32
	FilenameOffset uint32
}

func isValidSignatureHeader(header SignedHeader) (bool, int64) {
	validVersionTag := header.Tag1 == 1 && header.Length1 == 2
	validLengthTag := header.Tag2 == 2 && header.Length2 == 2
	return validVersionTag && validLengthTag, int64(header.HeaderSize)
}
