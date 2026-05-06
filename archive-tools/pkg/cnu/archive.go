package cnu

import "bytes"

func HasArchiveHeader(data []byte) bool {
	return bytes.Equal([]byte("CNU_File_Archive_3.0"), data[:20])
}
