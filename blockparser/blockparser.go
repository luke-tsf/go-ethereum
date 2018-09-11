package blockparser

import (
	"math/big"
	"fmt"
	"time"
)

var code = ""

func getCode(_code []byte) (success bool) {
	code = _code
	fmt.println(code)
	return true
}