package core

import (
	"math/big"

	"github.com/blockcypher/go-ethereum/common"
	"github.com/blockcypher/go-ethereum/core/types"
	"github.com/blockcypher/go-ethereum/core/vm"
)

// Implementation of evm.InternalTxListener

type InternalTxWatcher struct {
	internals types.InternalTransactions
}

func NewInternalTxWatcher() *InternalTxWatcher {
	return &InternalTxWatcher{
		internals: make(types.InternalTransactions, 0),
	}
}

// Public API: for users
func (self *InternalTxWatcher) SetParentHash(ph common.Hash) {
	for i := range self.internals {
		self.internals[i].ParentHash = ph
	}
}

func (self *InternalTxWatcher) InternalTransactions() types.InternalTransactions {
	return self.internals
}

// Public API: For interfacing with EVM
func (self *InternalTxWatcher) RegisterCall(nonce uint64, gasPrice *big.Int, gas uint64, srcAddr, dstAddr common.Address, value *big.Int, data []byte, depth uint64) {
	self.internals = append(self.internals,
		types.NewInternalTransaction(nonce, gasPrice, gas,
			srcAddr, dstAddr, value, data, depth, self.index(), "call"))
}

func (self *InternalTxWatcher) RegisterStaticCall(nonce uint64, gasPrice *big.Int, gas uint64, srcAddr, dstAddr common.Address, data []byte, depth uint64) {

	self.internals = append(self.internals,
		types.NewInternalTransaction(nonce, gasPrice, gas,
			srcAddr, dstAddr, big.NewInt(0), data, depth, self.index(),
			"staticcall"))
}

func (self *InternalTxWatcher) RegisterCallCode(nonce uint64, gasPrice *big.Int, gas uint64, contractAddr common.Address, value *big.Int, data []byte, depth uint64) {
	self.internals = append(self.internals,
		types.NewInternalTransaction(nonce, gasPrice, gas,
			contractAddr, contractAddr, value, data, depth, self.index(),
			"call"))
}

func (self *InternalTxWatcher) RegisterCreate(nonce uint64, gasPrice *big.Int, gas uint64, srcAddr, newContractAddr common.Address, value *big.Int, data []byte, depth uint64) {
	self.internals = append(self.internals,
		types.NewInternalTransaction(nonce, gasPrice, gas,
			srcAddr, newContractAddr, value, data, depth, self.index(),
			"create"))
}

func (self *InternalTxWatcher) RegisterDelegateCall(nonce uint64, gasPrice *big.Int, gas uint64, callerAddr common.Address, value *big.Int, data []byte, depth uint64) {
	self.internals = append(self.internals,
		types.NewInternalTransaction(nonce, gasPrice, gas,
			callerAddr, callerAddr, value, data, depth, self.index(), "call"))
}

func (self *InternalTxWatcher) RegisterSuicide(nonce uint64, gasPrice *big.Int, gas uint64, contractAddr, creatorAddr common.Address, remainingValue *big.Int, depth uint64) {
	self.internals = append(self.internals,
		types.NewInternalTransaction(nonce, gasPrice, gas,
			contractAddr, creatorAddr, remainingValue,
			append([]byte{byte(vm.SELFDESTRUCT)}, creatorAddr[:]...),
			depth, self.index(), "suicide"))
}

// Utilities
func (self *InternalTxWatcher) index() uint64 {
	return uint64(len(self.internals))
}

func toBigInt(g uint64) *big.Int {
	return big.NewInt(0).SetUint64(g)
}
