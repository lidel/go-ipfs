package bloom

import (
	"fmt"
	"testing"
)

func TestFilter(t *testing.T) {
	f := NewFilter(128)

	keys := [][]byte{
		[]byte("hello"),
		[]byte("fish"),
		[]byte("ipfsrocks"),
		[]byte("i want ipfs socks"),
	}

	// fmt.Println(f)

	f.Add(keys[0])
	if !f.Find(keys[0]) {
		t.Fatal("Failed to find single inserted key!")
	}

	f.Add(keys[1])
	if !f.Find(keys[1]) {
		t.Fatal("Failed to find key!")
	}

	f.Add(keys[2])
	f.Add(keys[3])

	for _, k := range keys {
		if !f.Find(k) {
			t.Fatal("Couldnt find one of three keys")
		}
	}

	if f.Find([]byte("beep boop")) {
		t.Fatal("Got false positive! Super unlikely!")
	}

	fmt.Println(f)

}
