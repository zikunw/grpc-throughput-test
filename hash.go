package main

import (
	"crypto/md5"
	"crypto/sha256"
	"crypto/sha512"
	"hash/crc32"
	"hash/crc64"
	"time"

	"github.com/cespare/xxhash"
	"github.com/spaolacci/murmur3"
)

func SHA_256() time.Duration {
	start := time.Now()
	for i := 0; i < WORDS; i++ {
		sha_256 := sha256.New()
		sha_256.Write([]byte("Hello World"))
	}
	// log time
	//log.Printf("Took %v to send %d messages", time.Since(start), WORDS)
	return time.Since(start)
}

func SHA_512() time.Duration {
	start := time.Now()
	for i := 0; i < WORDS; i++ {
		sha_512 := sha512.New()
		sha_512.Write([]byte("Hello World"))
	}
	// log time
	//log.Printf("Took %v to send %d messages", time.Since(start), WORDS)
	return time.Since(start)
}

func xxHash64() time.Duration {
	start := time.Now()
	for i := 0; i < WORDS; i++ {
		xxHash := xxhash.New()
		xxHash.Write([]byte("Hello World"))
	}
	// log time
	//log.Printf("Took %v to send %d messages", time.Since(start), WORDS)
	return time.Since(start)
}

func MD5() time.Duration {
	start := time.Now()
	for i := 0; i < WORDS; i++ {
		md5 := md5.New()
		md5.Write([]byte("Hello World"))
	}
	// log time
	//log.Printf("Took %v to send %d messages", time.Since(start), WORDS)
	return time.Since(start)
}

func Murmur32() time.Duration {
	start := time.Now()
	for i := 0; i < WORDS; i++ {
		murmur32 := murmur3.New32()
		murmur32.Write([]byte("Hello World"))
	}
	// log time
	//log.Printf("Took %v to send %d messages", time.Since(start), WORDS)
	return time.Since(start)
}

func Murmur64() time.Duration {
	start := time.Now()
	for i := 0; i < WORDS; i++ {
		murmur64 := murmur3.New64()
		murmur64.Write([]byte("Hello World"))
	}
	// log time
	//log.Printf("Took %v to send %d messages", time.Since(start), WORDS)
	return time.Since(start)
}

func CRC32() time.Duration {
	start := time.Now()
	for i := 0; i < WORDS; i++ {
		crc32 := crc32.NewIEEE()
		crc32.Write([]byte("Hello World"))
	}
	// log time
	//log.Printf("Took %v to send %d messages", time.Since(start), WORDS)
	return time.Since(start)
}

func CRC64() time.Duration {
	start := time.Now()
	for i := 0; i < WORDS; i++ {
		crc64 := crc64.New(crc64.MakeTable(crc64.ECMA))
		crc64.Write([]byte("Hello World"))
	}
	// log time
	//log.Printf("Took %v to send %d messages", time.Since(start), WORDS)
	return time.Since(start)
}
