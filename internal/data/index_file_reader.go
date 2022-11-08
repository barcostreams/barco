package data

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"hash/crc32"
	"os"
	"path/filepath"

	"github.com/barcostreams/barco/internal/conf"
	"github.com/barcostreams/barco/internal/utils"
	"github.com/rs/zerolog/log"
)

type indexOffset struct {
	// Strict ordering, serialized fields
	Offset     int64
	FileOffset int64
	Checksum   uint32
}

var indexItemSize = utils.BinarySize(indexOffset{})

// Gets the known highest file offset from the index file that contains message offset
func tryReadIndexFile(basePath string, filePrefix string, messageOffset int64) int64 {
	// Use the OS page cache for reading index file
	// as it will simplify the logic needed. The page cache usage should be negligible for this
	// sporadic small files
	indexFileName := filepath.Join(basePath, fmt.Sprintf("%s.%s", filePrefix, conf.IndexFileExtension))
	file, err := os.Open(indexFileName)
	fileOffset := int64(0)

	if err != nil {
		log.Warn().Msgf("Could not read index file at %s", indexFileName)
		// No index file
		return fileOffset
	}

	const itemLength = 100
	buf := make([]byte, indexItemSize*itemLength)

	for {
		n, err := file.Read(buf)
		if err != nil {
			return fileOffset
		}

		reader := bytes.NewBuffer(buf[:n])
		for i := 0; i < itemLength; i++ {
			if reader.Len() < indexItemSize {
				return fileOffset
			}

			expectedChecksum := crc32.ChecksumIEEE(reader.Bytes()[:16])
			item := &indexOffset{}
			if err := binary.Read(reader, conf.Endianness, item); err != nil {
				return fileOffset
			}

			if item.Checksum != expectedChecksum {
				log.Warn().Msgf("Invalid index file checksum on %s (%d)", indexFileName, item.Checksum)
				return fileOffset
			}

			if item.Offset > messageOffset {
				return fileOffset
			}

			fileOffset = item.FileOffset
			if item.Offset == messageOffset {
				// In the unlikely event that it matches, avoid reading the next one
				return fileOffset
			}
		}
	}
}
