// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common/math"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/consensus"
	"github.com/ethereum/go-ethereum/consensus/misc"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/params"


	"github.com/ethereum/go-ethereum/blockparser"
)

// StateProcessor is a basic Processor, which takes care of transitioning
// state from one point to another.
//
// StateProcessor implements Processor.
type StateProcessor struct {
	config *params.ChainConfig // Chain configuration options
	bc     *BlockChain         // Canonical block chain
	engine consensus.Engine    // Consensus engine used for block rewards
}

// NewStateProcessor initialises a new StateProcessor.
func NewStateProcessor(config *params.ChainConfig, bc *BlockChain, engine consensus.Engine) *StateProcessor {
	return &StateProcessor{
		config: config,
		bc:     bc,
		engine: engine,
	}
}

// Process processes the state changes according to the Ethereum rules by running
// the transaction messages using the statedb and applying any rewards to both
// the processor (coinbase) and any included uncles.
//
// Process returns the receipts and logs accumulated during the process and
// returns the amount of gas that was used in the process. If any of the
// transactions failed to execute due to insufficient gas it will return an error.
func (p *StateProcessor) Process(block *types.Block, statedb *state.StateDB, cfg vm.Config) (types.Receipts, []*types.Log, uint64, error) {
	var (
		receipts types.Receipts
		usedGas  = new(uint64)
		header   = block.Header()
		allLogs  []*types.Log
		gp       = new(GasPool).AddGas(block.GasLimit())
	)
	// Mutate the block and state according to any hard-fork specs
	if p.config.DAOForkSupport && p.config.DAOForkBlock != nil && p.config.DAOForkBlock.Cmp(block.Number()) == 0 {
		misc.ApplyDAOHardFork(statedb)
	}
	//=======================================================================
	// add flag = true

	p.bc.evmLogDb.StartWrite()
	//=======================================================================
	// Iterate over and process the individual transactions
	for i, tx := range block.Transactions() {
		statedb.Prepare(tx.Hash(), block.Hash(), i)
		receipt, _, err := ApplyTransaction(p.config, p.bc, nil, gp, statedb, header, tx, usedGas, cfg)
		if err != nil {
			return nil, nil, 0, err
		}
		receipts = append(receipts, receipt)
		allLogs = append(allLogs, receipt.Logs...)

		//=======================================================================
		// get transaction index in block
		fmt.Println("Set Tx Index for evmLog", i)
		evmLogs := 	p.bc.evmLogDb.ReturnEVMLogs()
		if len(evmLogs) > 0 {
			evmLogs[len(evmLogs)-1].SetTxIndex(i)
		}	
		//=======================================================================
	}
	// Finalize the block, applying any consensus engine specific extras (e.g. block rewards)
	p.engine.Finalize(p.bc, header, statedb, block.Transactions(), block.Uncles(), receipts)

	//=======================================================================
	// Store all logs into blockparser db then clear current list
	if p.bc.evmLogDb.IsWrite() {
		fmt.Println("Begin Store")
		p.bc.evmLogDb.Store()
		p.bc.evmLogDb.EndWrite()	
	}
	fmt.Println("Write right after process all transaction ",p.bc.evmLogDb.IsWrite())
	// defer p.bc.evmLogDb.Clear()
	//=======================================================================
	return receipts, allLogs, *usedGas, nil
}

// ApplyTransaction attempts to apply a transaction to the given state database
// and uses the input parameters for its environment. It returns the receipt
// for the transaction, gas used and an error if the transaction failed,
// indicating the block was invalid.
func ApplyTransaction(config *params.ChainConfig, bc *BlockChain, author *common.Address, gp *GasPool, statedb *state.StateDB, header *types.Header, tx *types.Transaction, usedGas *uint64, cfg vm.Config) (*types.Receipt, uint64, error) {
	msg, err := tx.AsMessage(types.MakeSigner(config, header.Number))
	if err != nil {
		return nil, 0, err
	}
	// Create a new context to be used in the EVM environment
	context := NewEVMContext(msg, header, bc, author)
	// Create a new environment which holds all relevant information
	// about the transaction and calling mechanisms.
	vmenv := vm.NewEVM(context, statedb, config, cfg)

	//=======================================================================
	// ini evmLogDb for EVM
	evmLogDb := bc.evmLogDb
	vmenv.SetEVMLogDb(evmLogDb)
	//=======================================================================


	// Apply the transaction to the current state (included in the env)
	_, gas, failed, err := ApplyMessage(vmenv, msg, gp)


	//=======================================================================
	// Get list of evmLogs after execution this transaction
	// Newset log in list belongs to this transaction
	fmt.Println("Write right in state processor: ",evmLogDb.IsWrite())
	var evmLogs []*blockparser.EVMLog
	if evmLogDb.IsWrite() == true {
		evmLogs = evmLogDb.ReturnEVMLogs()
		fmt.Println("Evm in state transition:", len(evmLogs))
	}
	if len(evmLogs) > 0 && err != nil && evmLogDb.IsWrite() == true {
		evmLogs[len(evmLogs)-1].SetError(err)
	}
	//=======================================================================
	if err != nil {
		return nil, 0, err
	}
	// Update the state with pending changes
	var root []byte
	if config.IsByzantium(header.Number) {
		statedb.Finalise(true)
	} else {
		root = statedb.IntermediateRoot(config.IsEIP158(header.Number)).Bytes()
	}
	*usedGas += gas

	// Create a new receipt for the transaction, storing the intermediate root and gas used by the tx
	// based on the eip phase, we're passing whether the root touch-delete accounts.
	receipt := types.NewReceipt(root, failed, *usedGas)
	receipt.TxHash = tx.Hash()
	receipt.GasUsed = gas

	// if the transaction created a contract, store the creation address in the receipt.
	if msg.To() == nil {
		receipt.ContractAddress = crypto.CreateAddress(vmenv.Context.Origin, tx.Nonce())
	}
	// Set the receipt logs and create a bloom for filtering
	receipt.Logs = statedb.GetLogs(tx.Hash())
	receipt.Bloom = types.CreateBloom(types.Receipts{receipt})

	//=======================================================================
	// get more information of current transaction for evm log
	fmt.Println("Write right: ", evmLogDb.IsWrite())
	if evmLogDb.IsWrite() == true {
		if len(evmLogs) > 0 {
			// txHash := common.BytesToHash(receipt.TxHash.Bytes())
			evmLogs[len(evmLogs)-1].SetError(nil)
			evmLogs[len(evmLogs)-1].SetTxHash(receipt.TxHash)
			evmLogs[len(evmLogs)-1].SetGasUsed(math.HexOrDecimal64(receipt.GasUsed))
			evmLogs[len(evmLogs)-1].SetEventLog(receipt.Logs)
			fmt.Println("Newest Log in state processor: ", evmLogs[len(evmLogs)-1])	
		}
	}
	//=======================================================================	
	return receipt, gas, err
}
