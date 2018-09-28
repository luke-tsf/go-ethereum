package blockparser

// call value from database
import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	// "github.com/syndtr/goleveldb/leveldb/iterator"
	// "github.com/syndtr/goleveldb/leveldb"
	"github.com/ethereum/go-ethereum/ethdb"
)
type TSFBackendAPI struct {
	evmLogDb *EVMLogDb
}

func NewTSFBackendAPI(_evmLogDb *EVMLogDb) *TSFBackendAPI{
	return &TSFBackendAPI{_evmLogDb}
}

func (tsfBackendAPI *TSFBackendAPI) GetTokenInfo(tokenAddress common.Address) (string){
	evmLogDb := tsfBackendAPI.evmLogDb
	tokenAddressString := tokenAddress.String()
	tokenValue, err := evmLogDb.customDb.Get([]byte(tokenAddressString))
	if err != nil{
		return "error"
	}
	fmt.Println("Value of Token", tokenAddressString, string(tokenValue))
	return string(tokenValue)
}

func (tsfBackendAPI *TSFBackendAPI) GetAccountHistory(address common.Address) (string){
	result := ""

	fmt.Println("Enter get account History")
	evmLogDb := tsfBackendAPI.evmLogDb
	addressString := address.String()
	prefix := addressString

	ldbDatabase := evmLogDb.customDb.(*ethdb.LDBDatabase)
	fmt.Println("Address to get history: ", addressString)
	iter := ldbDatabase.NewIterator()
	fmt.Println("iterator:", iter)

	for ok := iter.Seek([]byte(prefix)); ok && strings.HasPrefix(string(iter.Key()), prefix); ok = iter.Next(){
		key := iter.Key()
		value := iter.Value()
		fmt.Println("History", addressString, string(key), string(value))
		result += string(key) + "||" + string(value) + "***"
	}

	if result != "" {
		result = result[:len(result)-3]
		return result
	} else {
		return "Nothing"
	}
}