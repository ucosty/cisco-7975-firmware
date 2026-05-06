package cnu

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
)

// TLV represents a parsed Type-Length-Value element.
type TLV struct {
	Tag      byte
	Length   uint16
	Value    []byte
	Children []TLV
}

// isContainer checks if a tag is known to contain nested TLV elements.
// Based on the Cisco SBN structure:
// * Tag 3: Certificate info container
// * Tag 7: Cryptographic info container
// * Tag 9: Signature parameters container
func isContainer(tag byte) bool {
	return tag == 3 || tag == 7 || tag == 9
}

// ParseSBN parses a byte slice into a slice of TLV structures.
func ParseSBN(data []byte) ([]TLV, error) {
	var tlvs []TLV
	buf := bytes.NewReader(data)

	for buf.Len() > 0 {
		// Read the tag
		tag, err := buf.ReadByte()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to read tag: %w", err)
		}

		// Read the length (2 bytes, Big-Endian)
		var length uint16
		if err := binary.Read(buf, binary.BigEndian, &length); err != nil {
			return nil, fmt.Errorf("failed to read length for tag %02x: %w", tag, err)
		}

		// Read the value
		val := make([]byte, length)
		if _, err := io.ReadFull(buf, val); err != nil {
			return nil, fmt.Errorf("failed to read value for tag %02x (expected %d bytes): %w", tag, length, err)
		}

		tlv := TLV{
			Tag:    tag,
			Length: length,
			Value:  val,
		}

		if isContainer(tag) {
			children, err := ParseSBN(val)
			if err != nil {
				return nil, fmt.Errorf("failed to parse children of tag %02x: %w", tag, err)
			}
			tlv.Children = children
		}

		tlvs = append(tlvs, tlv)

		if tag == 0x0e {
			break
		}
	}

	return tlvs, nil
}

// getTagName converts a tag id into a human-readable name
func getTagName(tag byte) string {
	switch tag {
	case 0x01:
		return "version"
	case 0x02:
		return "length"
	case 0x03:
		return "certificate"
	case 0x04:
		return "signer"
	case 0x06:
		return "ca"
	case 0x0c:
		return "signature"
	case 0x0e:
		return "filename"
	default:
		return "unknown"
	}
}

// PrintSBN is a helper to pretty-print the nested structure
func PrintSBN(tlvs []TLV, indentLevel int) {
	indent := strings.Repeat("  ", indentLevel)
	for _, t := range tlvs {
		if len(t.Children) > 0 {
			fmt.Printf("%s[Tag %02x] Container (Length: %d)\n", indent, t.Tag, t.Length)
			PrintSBN(t.Children, indentLevel+1)
		} else {
			displayVal := formatValue(t.Value)
			fmt.Printf("%s[Tag 0x%02x : %s] Length: %-4d | Value: %s\n", indent, t.Tag, getTagName(t.Tag), t.Length, displayVal)
		}
	}
}

// GetSignatureLength is a helper which returns the size of the signature block in bytes
func GetSignatureLength(tlvs []TLV) (int, error) {
	for _, t := range tlvs {
		if t.Tag == 0x02 {
			length := binary.BigEndian.Uint16(t.Value)
			return int(length), nil
		}
	}
	return 0, fmt.Errorf("could not find signature length entry")
}

// formatValue does a rudimentary check to print strings nicely or hex dump binary data
func formatValue(val []byte) string {
	if len(val) == 0 {
		return "<empty>"
	}

	if val[len(val)-1] == 0x00 && len(val) > 1 {
		isPrintable := true
		for i := 0; i < len(val)-1; i++ {
			if val[i] < 32 || val[i] > 126 {
				isPrintable = false
				break
			}
		}

		if isPrintable {
			return fmt.Sprintf("\"%s\"", string(val[:len(val)-1]))
		}
	}

	// Truncate long hex outputs (like the 256-byte signature) for terminal readability
	hexStr := hex.EncodeToString(val)
	if len(hexStr) > 64 {
		return hexStr[:64] + "... (truncated)"
	}
	return hexStr
}
