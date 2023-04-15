package main

import (
	"fmt"
	"time"
)

const WORDS = 1000000

func main() {
	// Connection benchmarks

	// createServer()
	// createClient()
	// fmt.Scanln()

	// createServer()
	// createStreamClient()
	// fmt.Scanln()

	// createOneWriteTCPserver()
	// createOneWriteTCPclient()
	// fmt.Scanln()

	createServer()
	createRepeatedClient()
	fmt.Scanln()

	// createMultiWriteTCPserver()
	// createMultiWriteTCPclient()
	// fmt.Scanln()

	// go createAllWordsOneWriteTCPserver()
	// createAllWordsOneWriteTCPclient()
	// fmt.Scanln()

	// hash benchmarks
	// SHA_256()
	// SHA_512()
	// xxHash64()
	// fmt.Printf("SHA_256: %v\n", averageTime(SHA_256))
	// fmt.Printf("SHA_512: %v\n", averageTime(SHA_512))
	// fmt.Printf("xxHash64: %v\n", averageTime(xxHash64))
	// fmt.Printf("MD5: %v\n", averageTime(MD5))
	// fmt.Printf("Murmur32: %v\n", averageTime(Murmur32))
	// fmt.Printf("Murmur64: %v\n", averageTime(Murmur64))
	// fmt.Printf("CRC32: %v\n", averageTime(CRC32))
	// fmt.Printf("CRC64: %v\n", averageTime(CRC64))
}

// Run a function 100 times and return the average time
func averageTime(f func() time.Duration) time.Duration {
	var total time.Duration
	for i := 0; i < 100; i++ {
		total += f()
	}
	return total / 100
}
