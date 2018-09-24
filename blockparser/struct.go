package blockparser

import (
	"math/big"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/core/types"
)
// define stuct
// ERROR

type EVMLog struct {
	sender 				common.Address
	receiver 			common.Address
	value				*big.Int
	tokenERC20			common.Address
	tokenInformation	[]string
	blockNumber			*big.Int
	txIndex				int
	txHash				common.Hash
	gasUsed				math.HexOrDecimal64
	eventLog			[]*types.Log
	err 				error
}
type EVMLogDb struct {
	customDb		ethdb.Database
	evmLogs 		[]*EVMLog
}
func NewEVMLog(_sender common.Address, _receiver common.Address, _value *big.Int, 
	_tokenERC20 common.Address, _tokenInformation []string, _blockNumber *big.Int, 
	_txIndex int, _txHash common.Hash, _gasUsed math.HexOrDecimal64,  _err error) *EVMLog{
	return &EVMLog{
		sender: 			_sender,
		receiver:			_receiver,
		value:				_value,
		tokenERC20:			_tokenERC20,
		tokenInformation:	_tokenInformation,
		err:				_err,
		txHash:				_txHash,
		blockNumber:		_blockNumber,
		gasUsed:			_gasUsed,
		txIndex:			_txIndex,
	}
}

func NewEVMLogNewToken(_sender common.Address, _receiver common.Address, _value *big.Int, 
	_tokenERC20 common.Address, _tokenInformation []string, _err error) *EVMLog{
	return &EVMLog{
		sender: 			_sender,
		receiver:			_receiver,
		value:				_value,
		tokenERC20:			_tokenERC20,
		tokenInformation:	_tokenInformation,
		err:				_err,
	}
}
func (evmLog *EVMLog) SetError(err error) {
	evmLog.err = err
}
func (evmLog *EVMLog) SetTxHash(txHash common.Hash) {
	evmLog.txHash = txHash
}
func (evmLog *EVMLog) SetGasUsed(gasUsed math.HexOrDecimal64) {
	evmLog.gasUsed = gasUsed
}
func (evmLog *EVMLog) SetEventLog(eventLog []*types.Log) {
	evmLog.eventLog = eventLog
}

func (evmLog *EVMLog) SetTxIndex(txIndex int) {
	evmLog.txIndex = txIndex
}

func NewEVMLogDb(_customDb ethdb.Database) *EVMLogDb{
	return &EVMLogDb{
		customDb:	_customDb,
	}
}

func (evmLogDb *EVMLogDb) Close(){
	evmLogDb.customDb.Close()
}

