package blockparser

// new blockparser
// ~= main
import(
	"fmt"
	// "github.com/ethereum/go-ethereum/common"
	"encoding/hex"
)

func (evmLogDb *EVMLogDb) GetNewEVMLog(evmLog *EVMLog) (bool){
	sender := evmLog.sender.String()
	senderHex, err := hex.DecodeString(sender)
	if err != nil{
		return false
	}
	receiver := evmLog.receiver.String()
	receiverHex, err := hex.DecodeString(receiver)
	if err != nil{
		return false
	}
	// value := evmLog.value.String()
	evmLogDb.customDb.Put(senderHex,receiverHex)
	result,err := evmLogDb.customDb.Get(senderHex)
	if err != nil{
		return false
	}
	fmt.Println(result)
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