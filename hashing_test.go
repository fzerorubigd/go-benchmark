package benchmark

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"hash"
	"testing"

	"github.com/cespare/xxhash"
	"github.com/dchest/siphash"
	blake2bsimd "github.com/minio/blake2b-simd"
	"github.com/minio/highwayhash"
	xxhash32pier "github.com/pierrec/xxHash/xxHash32"
	xxhash64pier "github.com/pierrec/xxHash/xxHash64"
	"github.com/spaolacci/murmur3"
	xxhash32vova "github.com/vova616/xxhash"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"
)

const hashBufferSize = 8

func benchmarkHash(b *testing.B, hash func() hash.Hash, length int64) {
	data := make([]byte, length)
	b.SetBytes(length)
	var c uint64

	for i := 0; i < b.N; i++ {
		h := hash()
		_, err := h.Write(data[:])
		if err != nil {
			panic(err)
		}
		b := h.Sum(nil)
		c += uint64(b[0])
	}
	if c == 0 { // unlikely that this will ever fail but it prevents optimizing away the Sum() call
		b.Fail()
	}
}

func benchmarkHash64(b *testing.B, hash func() hash.Hash64, length int64) {
	data := make([]byte, length)
	b.SetBytes(length)
	var c uint64

	for i := 0; i < b.N; i++ {
		h := hash()
		_, err := h.Write(data[:])
		if err != nil {
			panic(err)
		}
		c += h.Sum64()
	}
	if c == 0 { // unlikely that this will ever fail but it prevents optimizing away the Sum() call
		b.Fail()
	}
}

func benchmarkHash64seed(b *testing.B, hash func(uint64) hash.Hash64, length int64) {
	data := make([]byte, length)
	b.SetBytes(length)
	var c uint64

	for i := 0; i < b.N; i++ {
		h := hash(1471)
		_, err := h.Write(data[:])
		if err != nil {
			panic(err)
		}
		c += h.Sum64()
	}
	if c == 0 { // unlikely that this will ever fail but it prevents optimizing away the Sum() call
		b.Fail()
	}
}

func benchmarkHash32seed(b *testing.B, hash func(uint32) hash.Hash32, length int64) {
	data := make([]byte, length)
	b.SetBytes(length)
	for i := 0; i < b.N; i++ {
		h := hash(1471)
		_, err := h.Write(data[:])
		if err != nil {
			panic(err)
		}
		_ = h.Sum32()
	}
}

func benchmarkHash64to32(b *testing.B, hash func() hash.Hash64, length int64) {
	data := make([]byte, length)
	b.SetBytes(length)
	for i := 0; i < b.N; i++ {
		h := hash()
		_, err := h.Write(data[:])
		if err != nil {
			panic(err)
		}
		var b []byte
		s := h.Sum64()
		_ = append(
			b,
			byte(s>>56),
			byte(s>>48),
			byte(s>>40),
			byte(s>>32),
		)
	}
}

func benchmarkHash64to16(b *testing.B, hash func() hash.Hash64, length int64) {
	data := make([]byte, length)
	b.SetBytes(length)
	for i := 0; i < b.N; i++ {
		h := hash()
		_, err := h.Write(data[:])
		if err != nil {
			panic(err)
		}
		var b []byte
		s := h.Sum64()
		_ = append(
			b,
			byte(s>>56),
			byte(s>>48),
		)
	}
}

func benchmarkHash64to8(b *testing.B, hash func() hash.Hash64, length int64) {
	data := make([]byte, length)
	b.SetBytes(length)
	for i := 0; i < b.N; i++ {
		h := hash()
		_, err := h.Write(data[:])
		if err != nil {
			panic(err)
		}
		s := h.Sum64()
		_ = byte(s >> 56)
	}
}

func benchmarkHashKeyError(b *testing.B, hash func([]byte) (hash.Hash, error), length int64) {
	data := make([]byte, length)
	b.SetBytes(length)
	key := make([]byte, 16)
	var c uint64

	for i := 0; i < b.N; i++ {
		h, _ := hash(key)
		_, err := h.Write(data[:])
		if err != nil {
			panic(err)
		}
		b := h.Sum(nil)
		c += uint64(b[0])
	}
	if c == 0 { // unlikely that this will ever fail but it prevents optimizing away the Sum() call
		b.Fail()
	}
}

func benchmarkHashKey64(b *testing.B, hash func([]byte) hash.Hash64, length int64) {
	data := make([]byte, length)
	b.SetBytes(length)
	key := make([]byte, 16)
	var c uint64

	for i := 0; i < b.N; i++ {
		h := hash(key)
		_, err := h.Write(data[:])
		if err != nil {
			panic(err)
		}
		c += h.Sum64()
	}
	if c == 0 { // unlikely that this will ever fail but it prevents optimizing away the Sum() call
		b.Fail()
	}
}

func benchmarkHashKey64Error(b *testing.B, hash func([]byte) (hash.Hash64, error), length int64) {
	data := make([]byte, length)
	b.SetBytes(length)
	key := make([]byte, 32)
	var c uint64

	for i := 0; i < b.N; i++ {
		h, _ := hash(key)
		_, err := h.Write(data[:])
		if err != nil {
			panic(err)
		}
		c += h.Sum64()
	}
	if c == 0 { // unlikely that this will ever fail but it prevents optimizing away the Sum() call
		b.Fail()
	}
}

func benchmarkHashKey32(b *testing.B, hash func([]byte) hash.Hash32, length int64) {
	data := make([]byte, length)
	b.SetBytes(length)
	key := make([]byte, 8)

	for i := 0; i < b.N; i++ {
		h := hash(key)
		_, err := h.Write(data[:])
		if err != nil {
			panic(err)
		}
		_ = h.Sum32()
	}
}

func BenchmarkHashing64MD5(b *testing.B) {
	benchmarkHash(b, md5.New, hashBufferSize)
}

func BenchmarkHashing64SHA1(b *testing.B) {
	benchmarkHash(b, sha1.New, hashBufferSize)
}

func BenchmarkHashing64SHA256(b *testing.B) {
	benchmarkHash(b, sha256.New, hashBufferSize)
}

func BenchmarkHashing64SHA3B224(b *testing.B) {
	benchmarkHash(b, sha3.New224, hashBufferSize)
}

func BenchmarkHashing64SHA3B256(b *testing.B) {
	benchmarkHash(b, sha3.New256, hashBufferSize)
}

func BenchmarkHashing64RIPEMD160(b *testing.B) {
	benchmarkHash(b, ripemd160.New, hashBufferSize)
}

func BenchmarkHashing64Blake2B(b *testing.B) {
	benchmarkHashKeyError(b, blake2b.New256, hashBufferSize)
}

func BenchmarkHashing64Blake2BSimd(b *testing.B) {
	benchmarkHash(b, blake2bsimd.New256, hashBufferSize)
}

func BenchmarkHashing64Murmur3(b *testing.B) {
	benchmarkHash64(b, murmur3.New64, hashBufferSize)
}
func BenchmarkHashing64SipHash(b *testing.B) {
	benchmarkHashKey64(b, siphash.New, hashBufferSize)
}
func BenchmarkHashing64XXHash(b *testing.B) {
	benchmarkHash64(b, xxhash.New, hashBufferSize)
}
func BenchmarkHashing64XXHashpier(b *testing.B) {
	benchmarkHash64seed(b, xxhash64pier.New, hashBufferSize)
}
func BenchmarkHashing32XXHashvova(b *testing.B) {
	benchmarkHash32seed(b, xxhash32vova.New, hashBufferSize)
}
func BenchmarkHashing32XXHashpier(b *testing.B) {
	benchmarkHash32seed(b, xxhash32pier.New, hashBufferSize)
}
func BenchmarkHashing32XXHash(b *testing.B) {
	benchmarkHash64to32(b, xxhash.New, hashBufferSize)
}
func BenchmarkHashing16XXHash(b *testing.B) {
	benchmarkHash64to16(b, xxhash.New, hashBufferSize)
}
func BenchmarkHashing8XXHash(b *testing.B) {
	benchmarkHash64to8(b, xxhash.New, hashBufferSize)
}

func BenchmarkHashing64HighwayHash(b *testing.B) {
	benchmarkHashKey64Error(b, highwayhash.New64, hashBufferSize)
}
