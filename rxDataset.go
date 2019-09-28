package randomx

import "C"
import (
	"log"
	"sync"
)

func NewRxDataset(flags ...Flag) *RxDataset {
	cache := NewRxCache(flags...)
	dataset := AllocDataset(flags...)

	return &RxDataset{
		dataset: dataset,
		rxCache: cache,

		workerNum: 1,
	}
}

func (ds *RxDataset) Close() {
	if ds.dataset != nil {
		ReleaseDataset(ds.dataset)
	}

	ds.rxCache.Close()
}

func (ds *RxDataset) Init(seed []byte, workerNum uint32) bool {
	if ds.rxCache.Init(seed) == false {
		log.Println("WARN: rxCache has already been initialized by the same seed")
	}

	if ds.rxCache == nil || ds.rxCache.cache == nil {
		return false
	}

	datasetItemCount := DatasetItemCount()
	var wg sync.WaitGroup

	for i := uint32(0); i < workerNum; i++ {
		a := (datasetItemCount * i) / workerNum
		b := (datasetItemCount * (i + 1)) / workerNum
		wg.Add(1)
		go func() {
			InitDataset(ds.dataset, ds.rxCache.cache, a, b-a)
			wg.Done()
		}()
	}
	wg.Wait()
	return true
}

// unfinished
func (ds *RxDataset) HugePages() (uint32, uint32) {
	return 0, 0
}
