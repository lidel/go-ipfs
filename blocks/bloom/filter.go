// package bloom implements a simple bloom filter.
package bloom

import (
	"encoding/binary"
	"errors"
	// Non crypto hash, because speed
	"hash"
	"hash/fnv"
)

type Filter interface {
	Add([]byte)
	Find([]byte) bool
	Merge(Filter) (Filter, error)
}

func NewFilter(size int) Filter {
	return &filter{
		hash:   fnv.New32a(),
		filter: make([]byte, size),
		k:      3,
	}
}

type filter struct {
	filter []byte
	hash   hash.Hash32
	k      int
}

func BasicFilter() Filter {
	return NewFilter(2048)
}

func (f *filter) Add(bytes []byte) {
	for _, bit := range getBitIndicies(f, bytes) {
		f.setBit(bit)
	}
}

func getBitIndicies(f *filter, bytes []byte) []uint32 {
	indicies := make([]uint32, f.k)

	f.hash.Write(bytes)
	b := make([]byte, 4)

	for i := 0; i < f.k; i++ {
		res := f.hash.Sum32()
		indicies[i] = res % (uint32(len(f.filter)) * 8)

		binary.LittleEndian.PutUint32(b, res)
		f.hash.Reset()
		f.hash.Write(b)
	}
	f.hash.Reset()

	return indicies
}

func (f *filter) Find(bytes []byte) bool {
	for _, bit := range getBitIndicies(f, bytes) {
		if !f.getBit(bit) {
			return false
		}
	}
	return true
}

func (f *filter) setBit(i uint32) {
	f.filter[i/8] |= (1 << byte(i%8))
}

func (f *filter) getBit(i uint32) bool {
	return f.filter[i/8]&(1<<byte(i%8)) != 0
}

func (f *filter) Merge(o Filter) (Filter, error) {
	casfil, ok := o.(*filter)
	if !ok {
		return nil, errors.New("Unsupported filter type")
	}

	if len(casfil.filter) != len(f.filter) {
		return nil, errors.New("filter lengths must match!")
	}

	nfilt := new(filter)

	// this bit is sketchy, need a way of comparing hash functions
	//nfilt.hashes = f.hashes

	nfilt.filter = make([]byte, len(f.filter))
	for i, v := range f.filter {
		nfilt.filter[i] = v | casfil.filter[i]
	}

	return nfilt, nil
}
