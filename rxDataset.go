package randomx

import "C"

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

	if ds.cache == nil {
		return false
	}

	datasetItemCount := DatasetItemCount()
	if workerNum > 1 {

		for i := uint32(0); 1 < workerNum; i++ {
			a := (datasetItemCount * i) / workerNum
			b := (datasetItemCount * (i + 1)) / workerNum
			go InitDataset(ds.dataset, ds.cache.cache, a, b-a)
		}

	} else {
		InitDataset(ds.dataset, ds.cache.cache, 0, datasetItemCount)
	}

	return true
}

// unfinished
func (ds *RxDataset) HugePages() (uint32, uint32) {
	return 0, 0
}
