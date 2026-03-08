package sdk

import (
	"encoding/binary"
)

// ===================================
// uint64 <-> []byte Helper
// ===================================

func u64ToBytes(val uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, val)

	// In Little Endian, leading zeros (most significant bytes) are at the end of the slice.
	// We iterate backwards to find the last non-zero byte.
	lastNonZeroIndex := len(b) - 1
	for lastNonZeroIndex >= 0 {
		if b[lastNonZeroIndex] != 0 {
			break
		}
		lastNonZeroIndex--
	}

	// If the value was 0, ensure we return at least one byte (0x00) instead of an empty slice.
	if lastNonZeroIndex < 0 {
		return []byte{0}
	}

	return b[:lastNonZeroIndex+1]
}

func bytesToU64(b []byte) uint64 {
	if len(b) > 8 {
		Abort("byte length less than or equal to 8")
	}

	// Create an 8-byte buffer initialized to zeros.
	buf := make([]byte, 8)

	// In Little Endian, the existing bytes are the least significant and go at the start.
	// Copy the input slice into the beginning of the buffer.
	copy(buf, b)

	return binary.LittleEndian.Uint64(buf)
}
