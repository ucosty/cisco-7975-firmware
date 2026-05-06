package cnu

import (
	"encoding/binary"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// PackFile holds our working metadata during the scanning phase
type PackFile struct {
	SourcePath  string
	ArchiveName string
	IsRaw       bool
	Size        uint32
	Timestamp   uint32
}

func Pack(inputDir, targetFile string, quiet bool) error {
	var packFiles []PackFile
	var namesBlock []byte

	// 1. Pass One: Scan directory and build metadata
	err := filepath.WalkDir(inputDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		// Calculate relative path for the archive
		relPath, err := filepath.Rel(inputDir, path)
		if err != nil {
			return err
		}

		// Firmware usually expects forward slashes, even if packed on Windows
		archiveName := filepath.ToSlash(relPath)

		// Identify the hidden kernel/bootloader files we extracted earlier
		isRaw := strings.HasPrefix(filepath.Base(archiveName), ".raw")

		packFiles = append(packFiles, PackFile{
			SourcePath:  path,
			ArchiveName: archiveName,
			IsRaw:       isRaw,
			Size:        uint32(info.Size()),
			Timestamp:   uint32(info.ModTime().Unix()),
		})
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to scan input directory: %w", err)
	}

	if len(packFiles) == 0 {
		return fmt.Errorf("input directory is empty")
	}

	// 2. Pre-calculate Offsets and build Data Blocks
	fileTableSize := uint32(len(packFiles) * 20) // 20 bytes per RawFileTableEntry
	currentNameOffset := uint32(HeaderSize) + fileTableSize

	var fileTable []RawFileTableEntry

	for _, pf := range packFiles {
		entry := RawFileTableEntry{
			F2:        2, // Standard flag per the Perl script
			Size:      pf.Size,
			Timestamp: pf.Timestamp,
		}

		if pf.IsRaw {
			entry.F1 = 6 // Special boot/kernel file
			entry.FilenameOffset = 0
		} else {
			entry.F1 = 7 // Normal file
			entry.FilenameOffset = currentNameOffset

			// Append null-terminated string to the names block
			nameBytes := append([]byte(pf.ArchiveName), 0x00)
			namesBlock = append(namesBlock, nameBytes...)

			currentNameOffset += uint32(len(nameBytes))
		}

		fileTable = append(fileTable, entry)
	}

	// The actual file contents begin immediately after the string block
	contentOffset := currentNameOffset

	// 3. Construct the Global Header
	var header RawHeader
	copy(header.Signature1[:], ExpectedSignature)
	header.NumFiles = uint32(len(packFiles))
	header.ContentOffset = contentOffset
	copy(header.Signature2[:], ExpectedSignature)

	// 4. Pass Two: Write Everything to Disk
	out, err := os.Create(targetFile)
	if err != nil {
		return fmt.Errorf("could not create output file: %w", err)
	}
	defer out.Close()

	// Write Header
	if err := binary.Write(out, binary.BigEndian, &header); err != nil {
		return fmt.Errorf("failed to write header: %w", err)
	}

	// Write File Table
	if err := binary.Write(out, binary.BigEndian, &fileTable); err != nil {
		return fmt.Errorf("failed to write file table: %w", err)
	}

	// Write Filenames Block
	if _, err := out.Write(namesBlock); err != nil {
		return fmt.Errorf("failed to write filenames block: %w", err)
	}

	if !quiet {
		fmt.Printf("Header generated. Content Offset: 0x%X\n", contentOffset)
	}

	// Write File Contents (Streaming)
	for _, pf := range packFiles {
		if err := appendFileContent(out, pf.SourcePath); err != nil {
			return fmt.Errorf("failed to write payload for %s: %w", pf.SourcePath, err)
		}
		if !quiet {
			fmt.Printf("Packed: %s (Size: %d)\n", pf.ArchiveName, pf.Size)
		}
	}

	if !quiet {
		fmt.Printf("\nSuccessfully packed %d files into %s\n", len(packFiles), targetFile)
	}
	return nil
}

// appendFileContent streams a file from disk into the open archive writer
func appendFileContent(archive io.Writer, sourcePath string) error {
	in, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer in.Close()

	_, err = io.Copy(archive, in)
	return err
}
