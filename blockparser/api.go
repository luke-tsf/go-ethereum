package blockparser

// call value from database
import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
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
		return ""
	}
	fmt.Println("Value of Token", tokenAddressString, string(tokenValue))
	return string(tokenValue)
}