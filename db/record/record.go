package record

import (
	"encoding/binary"
	"hash/crc32"
)

const (
	valueKind = iota
	tombstoneKind
)

const (
	kindByteSize   = 1
	crcLen         = 4
	keyLenByteSize = 4
	valLenByteSize = 4
	metaLength     = kindByteSize + crcLen + keyLenByteSize + valLenByteSize
)

// Record - a database record
type Record struct {
	kind  byte
	key   string
	value []byte
}

// NewValue - returns a new record of value type.
func NewValue(key string, value []byte) *Record {
	return &Record{
		kind:  valueKind,
		key:   key,
		value: value,
	}
}

// NewTomstone - returns a new record of tombstone type.
func NewTombstone(key string) *Record {
	return &Record{
		kind:  tombstoneKind,
		key:   key,
		value: []byte{},
	}
}

// Key - returns the record key.
func (r *Record) Key() string {
	return r.key
}

// Value - returns the record value.
func (r *Record) Value() []byte {
	return r.value
}

// IsTombstone - check if the record is tombstone type
func (r *Record) IsTombstone() bool {
	return r.kind == tombstoneKind
}

func (r *Record) ToBytes() []byte {
	keyBytes := []byte(r.key)
	keyLen := make([]byte, keyLenByteSize)
	binary.BigEndian.PutUint32(keyLen, uint32(len(keyBytes)))

	valLen := make([]byte, valLenByteSize)
	binary.BigEndian.PutUint32(valLen, uint32(len(r.value)))

	data := []byte{}
	crc := crc32.NewIEEE()

	for _, v := range [][]byte{
		{r.kind}, keyLen, valLen, []byte(r.key), r.value,
	} {
		data = append(data, v...)
		crc.Write(v)
	}

	crcData := make([]byte, crcLen)
	binary.BigEndian.PutUint32(crcData, crc.Sum32())

	return append(crcData, data...)
}
