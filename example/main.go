package main

import (
	"github.com/maoxs2/go-randomx"
	"log"
	"math/rand"
	"runtime"
	"strconv"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	start := time.Now()
	var count = 0

	for c := 0; c < runtime.NumCPU(); c++ {

		log.Println("cpu", c, "running")
		go func() {
			cache := randomx.AllocCache(randomx.JIT, randomx.HARD_AES)
			randomx.InitCache(cache, []byte("123"))
			vm := randomx.CreateVM(cache, nil, randomx.JIT, randomx.HARD_AES)

			for {
				nonce := strconv.FormatInt(rand.Int63(), 10)
				randomx.CalculateHash(vm, []byte("123"+nonce))
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
