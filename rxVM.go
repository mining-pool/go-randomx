package randomx

import "C"

func NewRxVM(dataset *RxDataset, scratchpad []byte, flags ...Flag) *RxVM {
	vm := CreateVM(dataset.cache.cache, dataset.dataset, flags...)
	return &RxVM{vm: vm}
}

func (vm *RxVM) Close() {
	if vm.vm != nil {
		DestroyVM(vm.vm)
	}
}
