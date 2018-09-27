package blockparser

// call value from database
import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	// "github.com/syndtr/goleveldb/leveldb/iterator"
	// "github.com/syndtr/goleveldb/leveldb"
	// "github.com/ethereum/go-ethereum/ethdb"
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
	// fmt.Println("Enter get account History")
	// evmLogDb := tsfBackendAPI.evmLogDb
	// addressString := address.String()
	// fmt.Println("Address to get history: ", addressString)
	// fmt.Println("Type of db", reflect.TypeOf(evmLogDb.customDb)
	// LD, err := evmLogDb.customDb.NewLDBDatabase([]byte(addressString))
	// fmt.Println("iterator:", value)
	// tokenValue, err := evmLogDb.customDb.Get([]byte(addressString))
	// fmt.Println("Result in get account History", tokenValue, err)
	// if err != nil{
	// 	return "error"
	// }
	// fmt.Println("Account History", addressString, string(tokenValue))
	// return string(tokenValue)
	return "Nothing"
}