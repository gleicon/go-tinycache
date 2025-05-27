package datautils

import (
	"math"

	"github.com/spaolacci/murmur3"
)

// BloomFilter represents a Bloom filter data structure
type BloomFilter struct {
	bitset []uint64 // Using uint64 for efficient bit operations
	size   uint     // Size of the bitset in bits
	k      uint     // Number of hash functions
}

// New creates a new Bloom filter optimized for expectedElements with falsePositiveRate
func NewBloomFilter(expectedElements int, falsePositiveRate float64) *BloomFilter {
	// Calculate optimal size and number of hash functions
	size := optimalBitSize(expectedElements, falsePositiveRate)
	k := optimalHashCount(size, expectedElements)

	// Create a bitset with enough uint64 elements
	bitsetSize := (size + 63) / 64 // Round up to nearest uint64
	return &BloomFilter{
		bitset: make([]uint64, bitsetSize),
		size:   size,
		k:      k,
	}
}

// optimalBitSize calculates the optimal size of the bitset
func optimalBitSize(n int, p float64) uint {
	return uint(math.Ceil(-float64(n) * math.Log(p) / math.Pow(math.Log(2), 2)))
}

// optimalHashCount calculates the optimal number of hash functions
func optimalHashCount(size uint, n int) uint {
	return uint(math.Max(1, math.Round(float64(size)/float64(n)*math.Log(2))))
}

// Add adds an element to the Bloom filter
func (bf *BloomFilter) Add(data []byte) {
	for i := uint(0); i < bf.k; i++ {
		position := bf.getPosition(data, i)
		index, bit := position/64, position%64
		bf.bitset[index] |= 1 << bit
	}
}

// Contains checks if an element might be in the Bloom filter
func (bf *BloomFilter) Contains(data []byte) bool {
	for i := uint(0); i < bf.k; i++ {
		position := bf.getPosition(data, i)
		index, bit := position/64, position%64
		if bf.bitset[index]&(1<<bit) == 0 {
			return false
		}
	}
	return true
}

// getPosition calculates the bit position for a given element and hash function
func (bf *BloomFilter) getPosition(data []byte, hashNum uint) uint {
	// Create different hash functions using the seed value
	hash := murmur3.Sum64WithSeed(data, uint32(hashNum))
	return uint(hash % uint64(bf.size))
}

// Clear resets the Bloom filter
func (bf *BloomFilter) Clear() {
	for i := range bf.bitset {
		bf.bitset[i] = 0
	}
}

// Size returns the size of the Bloom filter in bits
func (bf *BloomFilter) Size() uint {
	return bf.size
}

// Remove removes an element from the Bloom filter. Bloom filters do not support true removal, but this function clears the bits for the element.
// Note: This does not guarantee that the element is no longer considered present, as other elements may hash to the same bits.
func (bf *BloomFilter) Remove(data []byte) {
	for i := uint(0); i < bf.k; i++ {
		position := bf.getPosition(data, i)
		index, bit := position/64, position%64
		bf.bitset[index] &^= 1 << bit // Clear the bit
	}
}
