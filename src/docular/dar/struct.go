package dar

// FileHeader describes a file within a zip file.
// See the zip spec for details.
type FileHeader struct {
	// Name is the name of the file.
	// It must be a relative path: it must not start with a drive
	// letter (e.g. C:) or leading slash, and only forward slashes
	// are allowed.
	Name string

	CreatorVersion     uint16
	ReaderVersion      uint16
	Flags              uint16
	Method             uint16
	ModifiedTime       uint16 // MS-DOS time
	ModifiedDate       uint16 // MS-DOS date
	CRC32              uint32
	CompressedSize     uint32 // Deprecated: Use CompressedSize64 instead.
	UncompressedSize   uint32 // Deprecated: Use UncompressedSize64 instead.
	CompressedSize64   uint64
	UncompressedSize64 uint64
	Extra              []byte
	ExternalAttrs      uint32 // Meaning depends on CreatorVersion
	Comment            string
}
