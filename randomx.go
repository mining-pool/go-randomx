package randomx

//go:generate cmake -G "Unix Makefiles" RandomX/
//go:generate make

//#cgo CFLAGS: -I./randomx
//#cgo LDFLAGS: -L${SRCDIR}/lib/randomXL -lrandomx
//#cgo LDFLAGS: -L${SRCDIR}/lib -lstdc++
//#cgo LDFLAGS: -L${SRCDIR}/lib -lwinpthread
//#cgo LDFLAGS: -static-libgcc
/*
#include "randomx.h"
*/
import "C"
import (
	"unsafe"
)

type Flag int

var (
	DEFAULT     Flag = 0 // for all default
	LARGE_PAGES Flag = 1 // for dataset & rxCache
	HARD_AES    Flag = 2 // for vm
	FULL_MEM    Flag = 4 // for vm
	JIT         Flag = 8 // for vm
	SECURE      Flag = 16
)

func (f Flag) toC() C.randomx_flags {
	return (C.randomx_flags)(f)
}

func AllocCache(flags ...Flag) *C.randomx_cache {
	var SumFlag Flag = DEFAULT
	var cache *C.randomx_cache

	for _, flag := range flags {
		SumFlag = SumFlag | flag
	}

	cache = C.randomx_alloc_cache(SumFlag.toC())
	if cache == nil {
		panic("failed to alloc mem for rxCache")
	}

	return cache
}

func InitCache(cache *C.randomx_cache, key []byte) {
	if len(key) == 0 {
		panic("key cannot be NULL")
	}

	C.randomx_init_cache(cache, unsafe.Pointer(&key[0]), C.size_t(len(key)))
}

func ReleaseCache(cache *C.randomx_cache) {
	C.randomx_release_cache(cache)
}

func AllocDataset(flags ...Flag) *C.randomx_dataset {
	var SumFlag Flag = DEFAULT
	for _, flag := range flags {
		SumFlag = SumFlag | flag
	}

	var dataset *C.randomx_dataset
	dataset = C.randomx_alloc_dataset(SumFlag.toC())
	if dataset == nil {
		panic("failed to alloc mem for dataset")
	}

	return dataset
}

func DatasetItemCount() uint32 {
	var length C.ulong
	length = C.randomx_dataset_item_count()
	return uint32(length)
}

func InitDataset(dataset *C.randomx_dataset, cache *C.randomx_cache, startItem uint32, itemCount uint32) {
	if dataset == nil {
		panic("alloc dataset mem is required")
	}

	if cache == nil {
		panic("alloc cache mem is required")
	}

	C.randomx_init_dataset(dataset, cache, C.ulong(startItem), C.ulong(itemCount))
}

func GetDatasetMemory(dataset *C.randomx_dataset) unsafe.Pointer {
	return C.randomx_get_dataset_memory(dataset)
}

func ReleaseDataset(dataset *C.randomx_dataset) {
	C.randomx_release_dataset(dataset)
}

func CreateVM(cache *C.randomx_cache, dataset *C.randomx_dataset, flags ...Flag) *C.randomx_vm {
	var SumFlag Flag = DEFAULT
	for _, flag := range flags {
		SumFlag = SumFlag | flag
	}

	if dataset == nil {
		panic("failed creating vm: using empty dataset")
	}

	vm := C.randomx_create_vm(SumFlag.toC(), cache, dataset)

	if vm == nil {
		panic("failed to create vm")
	}

	return vm
}

func SetVMCache(vm *C.randomx_vm, cache *C.randomx_cache) {
	C.randomx_vm_set_cache(vm, cache)
}

func SetVMDataset(vm *C.randomx_vm, dataset *C.randomx_dataset) {
	C.randomx_vm_set_dataset(vm, dataset)
}

func DestroyVM(vm *C.randomx_vm) {
	C.randomx_destroy_vm(vm)
}

func CalculateHash(vm *C.randomx_vm, in []byte) []byte {
	out := make([]byte, C.RANDOMX_HASH_SIZE)
	if vm == nil {
		panic("failed hashing: using empty vm")
	}

	C.randomx_calculate_hash(vm, unsafe.Pointer(&in[0]), C.size_t(len(in)), unsafe.Pointer(&out[0]))
	return out
}

type RxCache struct {
	seed      []byte
	cache     *C.randomx_cache
	initCount uint64
}

type RxDataset struct {
	dataset *C.randomx_dataset
	rxCache *RxCache

	workerNum uint32
}

type RxVM struct {
	vm        *C.randomx_vm
	rxDataset *RxDataset
}
