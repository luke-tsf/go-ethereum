package blockparser

import (
	"math/big"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethdb"
)
// define stuct
// ERROR

type EVMLog struct {
	sender 			common.Address
	receiver 		common.Address
	value			*big.Int
	tokenERC20		[]byte
	err 			string
}

type EVMLogDb struct {
	customDb		ethdb.Database
}
func NewEVMLog(_sender common.Address, _receiver common.Address, _value *big.Int, _tokenERC20 []byte, _err string) *EVMLog{
	return &EVMLog{
		sender: 	_sender,
		receiver:	_receiver,
		value:		_value,
		tokenERC20:	_tokenERC20,
		err:		_err,
	}
}

func NewEVMLogDb(_customDb ethdb.Database) *EVMLogDb{
	return &EVMLogDb{
		customDb:	_customDb,
	}
}

