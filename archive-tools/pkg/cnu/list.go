package cnu

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"strings"
)

func List(inputPath string) error {
	// Open the archive file
	f, err := os.Open(inputPath)
	if err != nil {
		return fmt.Errorf("could not open input file: %w", err)
	}
	defer f.Close()

	// Peek at the header bytes and determine if the file is signed
	var signedHeader SignedHeader
	if err := binary.Read(f, binary.BigEndian, &signedHeader); err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}

	isSigned, archiveOffset := isValidSignatureHeader(signedHeader)
	if isSigned {
		f.Seek(archiveOffset, io.SeekStart)
	}

	// Parse the main header
	var header RawHeader
	if err := binary.Read(f, binary.BigEndian, &header); err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}

	// Verify the signature
	sig1 := string(bytes.TrimRight(header.Signature1[:], "\x00"))
	if sig1 != ExpectedSignature {
		return fmt.Errorf("invalid archive signature: expected '%s', got '%s'", ExpectedSignature, sig1)
	}

	// Track the current offset for the file payload block
	rawFilesCount := 0

	// Parse the file table and extract
	for i := uint32(0); i < header.NumFiles; i++ {
		var entry RawFileTableEntry
		if err := binary.Read(f, binary.BigEndian, &entry); err != nil {
			return fmt.Errorf("failed to read file table entry %d: %w", i, err)
		}

		// Determine the filename
		var filename string
		if entry.F1 == 7 {
			filename, err = readNullTerminatedString(f, int64(entry.FilenameOffset)+archiveOffset)
			if err != nil {
				return fmt.Errorf("failed to read filename for entry %d: %w", i, err)
			}
		} else {
			rawFilesCount++
			filename = fmt.Sprintf("raw%d", rawFilesCount)
		}

		// Sanitize filename to prevent directory traversal escapes
		filename = strings.TrimPrefix(filename, "/")
		fmt.Printf("%s (Size: %d bytes)\n", filename, entry.Size)
	}
	return nil
}
