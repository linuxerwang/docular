/*
Package dar implements the Doc Archive file format.

========== Doc Archive Format (version 01 / 02 / 03) ==========

A doc archive file stores multiple document files (normally html files) into a single binary file,
compressively.

In general, it has the following format:

    01010...101010101010......10101010101010......10101010101010......10101010
    |---header---|---doc 1 content---|---doc 2 content---|----doc indices----|


[Header Section]

    The first 40 bytes of the archive file compose a header. A header records important information
    about where and how documents are stored in the archive.

    darxx..ssssnnnniiii.....................

    0x00-0x02 string "dar", meaning Doc ARchive
    0x03-0x04 string, version number, 01 | 02 | 03, current version is 03
    0x05-0x08 ssss: string for 4 bytes integer (little endian), address of document indices
    0x09-0x12 nnnn: string for 4 bytes integer (little endian), number of documents
    0x13-0x16 iiii: string for 4 bytes integer (little endian), index of default document
              ....: unused, bytes of zero


[Doc Contents Section]

    Contents of a document is divided by a constant into small chunks. For version 01 and 02, each
    chunk is 512*1024 bytes in length (before compressed); for version 03, each chunk is 32*1024
    bytes in length (before compressed). Note that the last chunk may be not exactly in same length.
    Each chunk will be compressed with zlib and then saved into the doc archive. To decompress each
    chunk, it has to record the length of compressed chunk at the same time. This is done by storing
    the 4 bytes length of integer followed by recording the compressed chunk. In all, a document is
    stored as follows:

    LLLL0101010101010101010LLLL01010101010101010101010101LLLL01010101010

    Here LLLL means a string for 4 bytes integer (little endian) - the length of compressed chunk.
    The chunk bytes follows right after the length string. But how can we know how many chunks a
    specific document has? The total length of the uncompressed document will be stored in the 3rd
    section and with this length (together with the constent chunk size) we can calculate how many
    chunks this document has:

    chunks = length//MAX_CHUNK_SIZE
    if length%MAX_CHUNK_SIZE != 0:
        chunks = chunks + 1


[Doc Indices Section]

    Bytes 0x05-0x08 in header records the address of document indices. The indices is as follows:

    AAAALLLLNNNNxxxxxxxxxxxxxxxxxxxxxAAAALLLLNNNNxxxxxxxxAAAALLLLNNNNxxxxxxxxxxxxxxxxx

    AAAA: string for 4 bytes integer (little endian), address of the document contents
    LLLL: string for 4 bytes integer (little endian), length of the whole uncompressed document
    NNNN: string for 4 bytes integer (little endian), length of the document name
    xxxx: bytes for the document name


[Default Document]

    Each doc archive may (or not) have a default document. When requested without a specific document
    name, the default document may be returned. Bytes 0x13-0x16 in header records the index of the
    default document, in the order of document indices, starting from 1. So if 0x13-0x16 is number
    15, then the 15th document in indices represents the default document. 0 means there's no default
    document.

*/

package dar
