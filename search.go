package randomx

/*
#include <stdint.h>
#include "randomx.h"

void search(randomx_vm* vm, const void* in, const uint64_t target, const uint64_t max_times, void* nonce, void* out, uint64_t* count)
{
    for ((*count)=0; (*count) < max_times; (*count)++)
	{
		*(uint32_t*)(nonce) = *(uint32_t*)(nonce) + 1;
		*(uint32_t*)(in+39) = *(uint32_t*)(nonce);
		randomx_calculate_hash(vm, in, 76, out);
		if (*(uint64_t*)(out+24) < target) {
			return;
		}
	}

	return;
}
*/
import "C"
import (
	"unsafe"
)

func Search(vm *C.randomx_vm, nonce []byte, maxTimes uint64, target uint64, in []byte) (hash []byte, count uint64) {
	out := make([]byte, C.RANDOMX_HASH_SIZE)
	if vm == nil {
		panic("failed hashing: using empty vm")
	}

	var cCount C.uint64_t
	C.search(vm, unsafe.Pointer(&in[0]), C.uint64_t(target), C.uint64_t(maxTimes), unsafe.Pointer(&nonce[0]), unsafe.Pointer(&out[0]), &cCount)
	return out, uint64(cCount)
}

func (vm *RxVM) Search(nonce []byte, maxTimes uint64, target uint64, in []byte) (hash []byte, count uint64) {
	hash, count = Search(vm.vm, nonce, maxTimes, target, in)
	return
}
