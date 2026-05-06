package cnu

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func Unpack(inputPath, outDir string, quiet bool) error {
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
	} else {
		archiveOffset = 0
		f.Seek(0, io.SeekStart)
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

	// Ensure output directory exists
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("could not create output directory: %w", err)
	}

	// Track the current offset for the file payload block
	currentDataOffset := int64(header.ContentOffset) + archiveOffset
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
		filename = strings.TrimPrefix(filename, "/") // Remove leading slashes
		cleanPath := filepath.Clean(filename)
		targetPath := filepath.Join(outDir, cleanPath)

		// Create subdirectories if the target path includes them
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return fmt.Errorf("failed to create subdirectories for %s: %w", filename, err)
		}

		// Extract the file data using io.Copy and SectionReader (zero-copy memory)
		if err := extractFile(f, targetPath, currentDataOffset, int64(entry.Size)); err != nil {
			return fmt.Errorf("failed to write %s: %w", filename, err)
		}

		// Restore file timestamps
		modTime := time.Unix(int64(entry.Timestamp), 0)
		if err := os.Chtimes(targetPath, modTime, modTime); err != nil {
			if !quiet {
				fmt.Printf("Warning: Could not set timestamp for %s\n", filename)
			}
		}

		if !quiet {
			fmt.Printf("Extracted: %s (Size: %d bytes)\n", filename, entry.Size)
		}

		// Advance the data offset for the next file
		currentDataOffset += int64(entry.Size)
	}

	if !quiet {
		fmt.Printf("Successfully unpacked %d files to %s\n", header.NumFiles, outDir)
	}
	return nil
}

// readNullTerminatedString jumps to a specific offset and reads characters until a null byte is hit.
func readNullTerminatedString(r io.ReaderAt, offset int64) (string, error) {
	// Assuming filenames won't exceed 256 characters for efficiency
	var buf [256]byte
	n, err := r.ReadAt(buf[:], offset)
	if err != nil && err != io.EOF {
		return "", err
	}

	idx := bytes.IndexByte(buf[:n], 0)
	if idx == -1 {
		return "", fmt.Errorf("null terminator not found within bounds")
	}

	return string(buf[:idx]), nil
}

// extractFile streams a specific chunk of the archive to a new file on disk.
func extractFile(archive *os.File, targetPath string, offset int64, size int64) error {
	out, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Creates a reader that only sees 'size' bytes starting from 'offset'
	section := io.NewSectionReader(archive, offset, size)

	_, err = io.Copy(out, section)
	return err
}
