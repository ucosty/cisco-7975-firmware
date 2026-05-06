# CNU Firmware Archive

This document describes the format of the Cisco CNU firmware archive files. It covers both signed and unsigned variants of this format. The file structure was determined by reverse engineering examples of the archive files use to update the firmware on Cisco IP phones from the 7945, 7965, and 7975 series.

## Data Types & Conventions

* Endianness: All integers (offsets, sizes, timestamps, flags) are strictly Big─Endian

* Strings: All text strings are ASCII and must be null─terminated (0x00)

* Path Separators: Directory structures must use forward slashes (/), regardless of the host operating system

* Alignment: There is no explicit byte─alignment padding between the file payloads

## Archive Structure

A CNU archive is a file archive format. It is an uncompressed archive with an optional digital signature.

```
┌────────────────────────────────┐
│  [Optional digital signature]  │
│                                │
├────────────────────────────────┤ 
│  Archive Header                │ 
│                                │
├────────────────────────────────┤
│  File Table Block              │
│                                │
├────────────────────────────────┤
│  Filename Block                │
│                                │
├────────────────────────────────┤
│  Content Payload Block         │
│                                │
└────────────────────────────────┘
```

### Archive Header

A CNU archive begins with a 584 byte header, describing the number of files in the archive as well as the position in the file where the file data can be located.

| Offset (Hex) | Length (Bytes) | Type   | Description                                                                                     |
|--------------|----------------|--------|-------------------------------------------------------------------------------------------------|
| 0x00         | 20             | String | Primary Signature: Must be exactly "CNU_File_Archive_3.0"                                       |
| 0x14         | 28             | Binary | Padding 1: Zero-filled (0x00)                                                                   |
| 0x30         | 4              | uint32 | File Count: Total number of files contained in the archive.                                     |
| 0x34         | 4              | uint32 | Content Offset: The absolute byte offset in the archive where the Content Payload Block begins. |
| 0x38         | 20             | String | Secondary Signature: Copy of "CNU_File_Archive_3.0"                                             |
| 0x4C         | 508            | Binary | Padding 2: Zero-filled (0x00)                                                                   |

### File Table Block

Immediately after the archive header is an array of file entries. Each file entry describes the attributes of the file and the offset to its filename.

| Offset (Hex) | Length (Bytes) | Type   | Description                  |
|--------------|----------------|--------|------------------------------|
| 0x00         | 4              | uint32 | File type                    |
| 0x04         | 4              | uint32 | File flags (unknown purpose) |
| 0x08         | 4              | uint32 | File size                    |
| 0x0c         | 4              | uint32 | Timestamp                    |
| 0x10         | 4              | uint32 | Absolute offset to filename  |


### Filename Block

The Filename Block is block of null-terminated ASCII filenames. These filenames are referenced by entries in the file table block.

### Content Payload Block

The Content Payload Block contains all of the archived files contiguously. There are offsets into this block provided, instead the position of any specific file needs to be caluclated. The position of the content payload block is described in the archive header, in the content offset field.

To correctly parse the content payload block, in order to get the contents of the files in the archive, you  must iterate over the file table block. Increment an offset into the content payload block by the size of the files iterated over.

```c++
auto file_offset = header.content_offset;

for(auto file: file_entries) {
    // Do something with the current file
    extract_file(file, file_offset)
    
    // Get the next file offset
    auto file_offset += file.size;
}
```

## ImHex Pattern

The unsigned version of the archive can be parsed with this ImHex pattern

```c
#pragma endian big

struct ArchiveHeader {
    char magic[20];
    char padding1[28];
    u32 file_count;
    u32 content_offset;
    char magic2[20];
    char padding2[508];
};

u32 current_data_offset;

struct FileEntry {
    u32 file_type;
    u32 flag;
    u32 size;
    u32 timestamp;
    u32 filename_offset;

    // Not all files have a filename, but most do
    if (filename_offset != 0) {
        char filename[] @ filename_offset;
    }

    u8 file_data[size] @ current_data_offset;
    current_data_offset += size;
};

ArchiveHeader header @ 0x00;
current_data_offset = header.content_offset;
FileEntry files[header.file_count] @ $;
```
