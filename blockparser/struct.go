package blockparser

import (
	"math/big"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
)
// define stuct
// ERROR

type EVMLog struct {
	sender 				common.Address
	receiver 			common.Address
	value				*big.Int
	tokenERC20			common.Address
	tokenInformation	[]string
	err 				error
	txHash				Hash *common.Hash
	blockNumber			*big.Int
	gasUsed				math.HexOrDecimal64
}
type EVMLogDb struct {
	customDb		ethdb.Database
	evmLog 			*[]EVMLog
}
func NewEVMLog(_sender common.Address, _receiver common.Address, _value *big.Int, _tokenERC20 common.Address, _tokenInformation []string, _err error) *EVMLog{
	return &EVMLog{
		sender: 			_sender,
		receiver:			_receiver,
		value:				_value,
		tokenERC20:			_tokenERC20,
		tokenInformation:	_tokenInformation,
		err:				_err,
	}
}

func NewEVMLogDb(_customDb ethdb.Database) *EVMLogDb{
	return &EVMLogDb{
		customDb:	_customDb,
	}
}

func (evmLogDb *EVMLogDb) Close(){
	evmLogDb.customDb.Close()
}

