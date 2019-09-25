package randomx

import (
	"math/rand"
	"strconv"
	"testing"
)

func TestAllocCache(t *testing.T) {
	cache := AllocCache(DEFAULT)
	InitCache(cache, []byte("123"))
	ReleaseCache(cache)
}

func TestAllocDataset(t *testing.T) {
	ds := AllocDataset(DEFAULT)
	t.Log(DatasetItemCount())
	cache := AllocCache(DEFAULT)
	InitCache(cache, []byte("123"))

	InitDataset(ds, cache, 0, 2)
	t.Log(GetDatasetMemory(ds))

	ReleaseDataset(ds)
	ReleaseCache(cache)
}

func TestCreateVM(t *testing.T) {
	cache := AllocCache(DEFAULT)
	ds := AllocDataset(DEFAULT)
	InitCache(cache, []byte("123"))
	InitDataset(ds, cache, 0, 200)

	vm := CreateVM(cache, ds, DEFAULT)
	t.Log(CalculateHash(vm, []byte("123")))

	DestroyVM(vm)
}

// go test -v -bench "." -benchtime=30s
func BenchmarkCalculateHash(b *testing.B) {
	cache := AllocCache(DEFAULT)
	ds := AllocDataset(DEFAULT)
	InitCache(cache, []byte("123"))
	InitDataset(ds, cache, 0, 200)

	vm := CreateVM(cache, ds, DEFAULT)
	for i := 0; i < b.N; i++ {
		nonce := strconv.FormatInt(rand.Int63(), 10)
		CalculateHash(vm, []byte("123"+nonce))
	}

	DestroyVM(vm)
}
