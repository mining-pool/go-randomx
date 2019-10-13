package randomx

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
	"math/rand"
	//"math/rand"
	"runtime"
	"strconv"
	//"strconv"
	"sync"
	"testing"
	"time"
)

var testPairs = [][][]byte{
	// randomX
	{
		[]byte("test key 000"),
		[]byte("This is a test"),
		[]byte("639183aae1bf4c9a35884cb46b09cad9175f04efd7684e7262a0ac1c2f0b4e3f"),
	},
	// randomXL
	{
		[]byte("test key 000"),
		[]byte("This is a test"),
		[]byte("b291ec8a532bc4f78bd75b43d211e1169bb65b1a8f66d4250376ba1d6fcff1bd"),
	},
}

func TestAllocCache(t *testing.T) {
	cache := AllocCache(FlagDefault)
	InitCache(cache, []byte("123"))
	ReleaseCache(cache)
}

func TestAllocDataset(t *testing.T) {
	ds := AllocDataset(FlagDefault)
	cache := AllocCache(FlagDefault)

	seed := make([]byte, 32)
	InitCache(cache, seed)
	log.Println("rxCache initialization finished")

	count := DatasetItemCount()
	log.Println("dataset count:", count/1024/1024, "mb")
	InitDataset(ds, cache, 0, count)
	log.Println(GetDatasetMemory(ds))

	ReleaseDataset(ds)
	ReleaseCache(cache)
}

func TestCreateVM(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	var tp = testPairs[1]
	cache := AllocCache(FlagDefault)
	log.Println("alloc cache mem finished")
	seed := tp[0]
	InitCache(cache, seed)
	log.Println("cache initialization finished")

	ds := AllocDataset(FlagDefault)
	log.Println("alloc dataset mem finished")
	count := DatasetItemCount()
	log.Println("dataset count:", count)
	var wg sync.WaitGroup
	var workerNum = uint32(runtime.NumCPU())
	for i := uint32(0); i < workerNum; i++ {
		wg.Add(1)
		a := (count * i) / workerNum
		b := (count * (i + 1)) / workerNum
		go func() {
			defer wg.Done()
			InitDataset(ds, cache, a, b-a)
		}()
	}
	wg.Wait()
	log.Println("dataset initialization finished") // too slow when one thread
	vm := CreateVM(cache, ds, FlagJIT, FlagHardAES, FlagFullMEM)

	var hashCorrect = make([]byte, hex.DecodedLen(len(tp[2])))
	_, err := hex.Decode(hashCorrect, tp[2])
	if err != nil {
		log.Println(err)
	}

	if bytes.Compare(CalculateHash(vm, tp[1]), hashCorrect) != 0 {
		t.Fail()
	}
}

func TestNewRxVM(t *testing.T) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	start := time.Now()
	pair := testPairs[1]
	workerNum := uint32(runtime.NumCPU())

	seed := pair[0]
	dataset := NewRxDataset(FlagJIT)
	if dataset.GoInit(seed, workerNum) == false {
		log.Fatal("failed to init dataset")
	}
	//defer dataset.Close()
	fmt.Println("Finished generating dataset in", time.Since(start).Seconds(), "sec")

	vm := NewRxVM(dataset, FlagFullMEM, FlagHardAES, FlagJIT, FlagSecure)
	//defer vm.Close()

	blob := pair[1]
	hash := vm.CalcHash(blob)

	var hashCorrect = make([]byte, hex.DecodedLen(len(pair[2])))
	_, err := hex.Decode(hashCorrect, pair[2])
	if err != nil {
		log.Println(err)
	}

	if bytes.Compare(hash, hashCorrect) != 0 {
		log.Println(hash)
		t.Fail()
	}
}

// go test -v -bench "." -benchtime=30m
func BenchmarkCalculateHash(b *testing.B) {
	cache := AllocCache(FlagDefault)
	ds := AllocDataset(FlagDefault)
	InitCache(cache, []byte("123"))
	FastInitFullDataset(ds, cache, uint32(runtime.NumCPU()))
	vm := CreateVM(cache, ds, FlagDefault)
	for i := 0; i < b.N; i++ {
		nonce := strconv.FormatInt(rand.Int63(), 10) // just test
		CalculateHash(vm, []byte("123"+nonce))
	}

	DestroyVM(vm)
}
