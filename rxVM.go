package randomx

func NewRxVM(rxDataset *RxDataset, flags ...Flag) *RxVM {
	vm := CreateVM(rxDataset.rxCache.cache, rxDataset.dataset, flags...)
	return &RxVM{
		vm:        vm,
		rxDataset: nil,
	}
}

func (vm *RxVM) Close() {
	if vm.vm != nil {
		DestroyVM(vm.vm)
	}
}

func (vm *RxVM) CalcHash(in []byte) []byte {
	return CalculateHash(vm.vm, in)
}

func (vm *RxVM) UpdateDataset(rxDataset *RxDataset) {
	SetVMCache(vm.vm, rxDataset.rxCache.cache)
	SetVMDataset(vm.vm, rxDataset.dataset)
}
