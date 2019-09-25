package randomx

import "C"
import (
	"sync"
)

func NewRxDataset(flags ...Flag) *RxDataset {
	dataset := AllocDataset(flags...)
	cache := NewRxCache(flags...)

	return &RxDataset{
		dataset: dataset,
		cache:   cache,
	}
}

func (ds *RxDataset) Close() {
	if ds.dataset != nil {
		ReleaseDataset(ds.dataset)
	}

	ds.cache.Close()
}

func (ds *RxDataset) Init(seed []byte, workerNum uint32) bool {
	ds.cache.Init(seed)

	if ds.cache == nil || ds.cache.cache == nil {
		return false
	}

	datasetItemCount := DatasetItemCount()
	var wg sync.WaitGroup

	for i := uint32(0); i < workerNum; i++ {
		a := (datasetItemCount * i) / workerNum
		b := (datasetItemCount * (i + 1)) / workerNum
		wg.Add(1)
		go func() {
			InitDataset(ds.dataset, ds.cache.cache, a, b-a)
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
