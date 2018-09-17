package blockparser

// new blockparser
// ~= main
import(
	"fmt"
	"reflect"

	// "github.com/ethereum/go-ethereum/common"
	// "encoding/hex"
)

func (evmLogDb *EVMLogDb) GetNewEVMLog(evmLog *EVMLog) (bool){
	batch := evmLogDb.customDb.NewBatch()
	var key = evmLog.sender.String()
	var value = string(evmLog.receiver.String())
	batch.Put([]byte(key), []byte(value))
	batch.Write()
	getValue, err := evmLogDb.customDb.Get([]byte(key))
	if err != nil{
		return false
	}
	fmt.Println(string(getValue))
	return true
}

func (evmLogDb *EVMLogDb) GetNewEVMLogToken(evmLog *EVMLog) (bool){
	batch := evmLogDb.customDb.NewBatch()
	var key = evmLog.tokenERC20.String()
	fmt.Println("Type of Key", reflect.TypeOf(key))
	var value = string(evmLog.tokenInformation)
	batch.Put([]byte(key), []byte(value))
	batch.Write()
	getValue, err := evmLogDb.customDb.Get([]byte(key))
	if err != nil{
		return false
	}
	// var address =  string(common.BytesToAddress(key))
	fmt.Println("Value of Token", string(key), string(getValue))
	return true
}
func (evmLogDb *EVMLogDb) TestDb() (bool){
	key := []byte{'k','e','y'}
	value := []byte{'v','a','l','u','e'}
	evmLogDb.customDb.Put(key,value)
	result, err := evmLogDb.customDb.Get(key)
	if err != nil{
		return false
	}
	fmt.Println(result)
	return true
}