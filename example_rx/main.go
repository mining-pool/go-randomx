package main

import (
	"bytes"
	"encoding/hex"
	"github.com/maoxs2/go-randomx"
	"log"
	"math/big"
	"math/rand"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	start := time.Now()
	var count = 0
	var flags = []randomx.Flag{
		randomx.JIT, randomx.HARD_AES, randomx.LARGE_PAGES, randomx.SECURE,
	}

	workerNum := uint32(runtime.NumCPU())

	for c := uint32(0); c < workerNum; c++ {

		log.Println("cpu", c, "running")
		go func() {
			cache := randomx.AllocCache(flags...)
			dataset := randomx.AllocDataset(flags...)
			strSeed := "2c19a83a045fbadc935b0e6967fa6614caa8de96a731113bcb9f5c4428ee7598"
			seed, _ := hex.DecodeString(strSeed)
			randomx.InitCache(cache, seed)
			strBlob := "0c0ca08deeea05c8a1aa685cc5335b882e99c0902249643a06b94395999c1642eb04a0dc8e3d1a00000000b9c5cf14074e4538fd89d9a33f0e150a3f258567bf2c031b95dd661f0e2cb22604"
			blob, _ := hex.DecodeString(strBlob)

			datasetItemCount := randomx.DatasetItemCount()
			for i := uint32(0); 1 < workerNum; i++ {
				a := (datasetItemCount * i) / workerNum
				b := (datasetItemCount * (i + 1)) / workerNum
				go randomx.InitDataset(dataset, cache, a, b-a)
			}

			vm := randomx.CreateVM(cache, dataset, flags...)

			for {
				nonce := big.NewInt(rand.Int63())

				randomx.CalculateHash(vm, bytes.Join([][]byte{blob, nonce.Bytes()}, nil))
				count++
				if count%100 == 0 {
					e := int(time.Since(start).Seconds())
					log.Println("speed", count/e, "h/s")
				}
			}
		}()
	}

	select {}
}
