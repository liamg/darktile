package fontinfo

import (
	"fmt"
	"io"
)

func read(r io.Reader, length int) ([]byte, error) {
	buf := make([]byte, length)
	if n, err := r.Read(buf); err != nil {
		return nil, err
	} else if n < length {
		return nil, fmt.Errorf("invalid length")
	}
	return buf, nil
}

func u16(buf []byte) uint16 {
	return (uint16(buf[0]) << 8) + uint16(buf[1])
}

func u32(buf []byte) uint32 {
	return (uint32(buf[0]) << 24) + (uint32(buf[1]) << 16) + (uint32(buf[2]) << 8) + uint32(buf[3])
}

type fontMetadata struct {
	FontFamily string
	FontStyle  string
}

func readMetadata(r io.ReadSeeker) (*fontMetadata, error) {

	buf, err := read(r, 12)
	if err != nil {
		return nil, err
	}

	tableCount := u16(buf[4:6])

	for i := 0; i < int(tableCount); i++ {

		if _, err := r.Seek(12+(int64(i)*16), 0); err != nil {
			return nil, err
		}

		table, err := read(r, 16)
		if err != nil {
			return nil, err
		}

		if string(table[0:4]) != "name" {
			continue
		}
		offset := u32(table[8:12])
		return readNameTable(r, offset)
	}

	return nil, fmt.Errorf("name table not found")
}

func readNameTable(r io.ReadSeeker, offset uint32) (*fontMetadata, error) {

	if _, err := r.Seek(int64(offset), 0); err != nil {
		return nil, fmt.Errorf("invalid font file")
	}

	nameTable, err := read(r, 6)
	if err != nil {
		return nil, err
	}

	nameCount := u16(nameTable[2:4])
	stringOffset := int64(u16(nameTable[4:6])) + int64(offset)

	var done uint8
	var metadata fontMetadata

	nameRecordStart := offset + 6

	for j := 0; j < int(nameCount); j++ {
		recordOffset := nameRecordStart + uint32(12*j)
		if _, err := r.Seek(int64(recordOffset), 0); err != nil {
			return nil, err
		}
		buf, err := read(r, 12)
		if err != nil {
			return nil, err
		}
		language := u16(buf[4:6])
		if language != 0 && language != 1033 { // not english or english us
			continue
		}
		nameID := u16(buf[6:8])

		switch nameID {
		case 1, 2: //family, style

			if _, err := r.Seek(int64(stringOffset)+int64(u16(buf[10:12])), 0); err != nil {
				return nil, err
			}
			raw, err := read(r, int(u16(buf[8:10])))
			if err != nil {
				return nil, err
			}
			if nameID == 1 {
				done |= 1
				metadata.FontFamily = string(raw)
			} else {
				done |= 2
				metadata.FontStyle = string(raw)
			}
			if done == 3 { // bail early if we have what we need
				return &metadata, nil
			}

		}
	}

	return &metadata, nil
}
