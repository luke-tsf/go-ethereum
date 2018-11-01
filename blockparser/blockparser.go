package blockparser

// new blockparser
// ~= main
import(
	"fmt"
	"reflect"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	// "github.com/ethereum/go-ethereum/common/hexutil"
	// "encoding/hex"
)

// Transfer: 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef
// Add new evm Log to current array of evm log db
var (
	addressZero = "0x0000000000000000000000000000000000000000"
	transferSigString = "0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef"
	transferSig = common.HexToHash(transferSigString)
	blockBigNumber = big.NewInt(100000000)
	txBigNumber = 99999
	eventBigNumber = uint(99)
)

func (evmLogDb *EVMLogDb) AddNewEVMLog(evmLog *EVMLog) (bool){
	evmLogDb.evmLogs = append(evmLogDb.evmLogs, evmLog) 
	fmt.Println("Get New EVM Log", evmLog)
	fmt.Println("Total EVM Log", len(evmLogDb.evmLogs))
	if reflect.DeepEqual(evmLog, evmLogDb.evmLogs[len(evmLogDb.evmLogs) - 1]){
		fmt.Println("Equal")
	} else {
		fmt.Println("Not equal")
	}
	return true				
}

// get a list of current evmlogs
func (evmLogDb *EVMLogDb) ReturnEVMLogs() ([]*EVMLog) {
	return evmLogDb.evmLogs
}

// Test if blockparser receive new evm log or not
func (evmLogDb *EVMLogDb) GetNewEVMLog(evmLog *EVMLog) (bool){
	fmt.Println("Get New EVM Log", evmLog)
	return true				
}
// Execute this function to store all current evmLogs in evmLogDb
func (evmLogDb *EVMLogDb) Store() {
	evmLogs := evmLogDb.evmLogs
	for i, evmLog := range evmLogs {
		fmt.Println("evm Log number : ", i, "in total: ", len(evmLogs))
		// fmt.Println("EvmLog Token ERC20", evmLog.tokenERC20.String())
		// fmt.Println("EvmLog Token Info", reflect.DeepEqual(evmLog.tokenInformation,[]string{}))
		if evmLog.err == nil {
			if evmLog.receiver.String() == addressZero && evmLog.tokenERC20.String() != addressZero{
				evmLogDb.storeEVMLogTokenInfo(evmLog)
			} else {
				evmLogDb.storeEVMLogERC20OrTransfer(evmLog)
			} 
		}	
	}
	evmLogDb.clear()
}
// Clear current list of evmlogs in evmlogDb
func (evmLogDb *EVMLogDb) clear() {
	evmLogDb.evmLogs = evmLogDb.evmLogs[:0]
	fmt.Println("Empty list of evmlog: ", len(evmLogDb.evmLogs) == 0)
}



func (evmLogDb *EVMLogDb) storeEVMLogERC20OrTransfer(evmLog *EVMLog) (bool){
	var value = evmLog.value.String()
	// call transfer
	fmt.Println("This is value of transaction ",value)
	evmLogDb.storeEVMLogTransfer(evmLog)

	if reflect.DeepEqual(evmLog.eventLog,[]*types.Log{}) == false {
		evmLogDb.storeEVMLogERC20(evmLog)
	}
	return true
}



// Store all evmLog to db 
/*
	Db for token information
						tokenAddress 					||	 name *|* symbol *|* decimal *|* totalSupply
	====================================================================================================================================
		0x2462fe786b651f19e43ba6c287da50c1790805a9		||   		Luke Coin*|*LUK*|*0*|*100000
*/
func (evmLogDb *EVMLogDb) storeEVMLogTokenInfo(evmLog *EVMLog) (bool){
	batch := evmLogDb.customDb.NewBatch()
	var key = evmLog.tokenERC20.String()
	// fmt.Println("Type of Key", reflect.TypeOf(key))
	var valueList = []string(evmLog.tokenInformation)
	value := valueList[0] + "*|*" + valueList[1] + "*|*" + valueList[2] + "*|*" + valueList[3]
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
// Store DB format
/*
	History DB for Ethereum Transfer Transaction
		address - 99million minor blockNumber - 99999 minor txIndex - transactionHash	||	 from-to-value-flag(receiver=0 or sender=1, 2 self)-flagIsRemoved(0:mainchian,1:oldchain,2:uncle) 	
	===========================================================================================================================================================================================
			0x123456-99999989-9998-89-0x789012											||   		0x123456-0xabcdef-10--0
	address: 0x123456
	blNumber: 10 => 99999989
	txIndex: 1 => 99998
	txHash: 0x789012

	from: 0x123456
	to: 0xabcdef
	value (amount of ETH): 10
	sender: 0x123456 => flag = 1
	flagIsRemoved: 0 - mainchain

	Transaction history format:
	Key: prefix is transactionHash
	value: key of this transaction
		transactionHash - address - 99million minor blockNumber - 99999 minor txIndex	||	 address - 99million minor blockNumber - 99999 minor txIndex - transactionHash 	
	===========================================================================================================================================================================================
			0x789012-0x123456-99999989-9998-89											||   		0x123456-99999989-9998-89-0x789012	
*/
func (evmLogDb *EVMLogDb) storeEVMLogTransfer(evmLog *EVMLog) (bool){
	batch := evmLogDb.customDb.NewBatch()
	var sender = evmLog.sender.String()
	var receiver = evmLog.receiver.String()
	var value = evmLog.value.String()

	// var blockNumberBigInt *big.Int
	// fmt.Println("Here 1")
	var bigNumber = big.NewInt(100000000)
	// fmt.Println("Init big int number", reflect.TypeOf(bigNumber), bigNumber)
	// fmt.Println("Get block number", reflect.TypeOf(evmLog.blockNumber), evmLog.blockNumber)
	bigNumber.Sub(bigNumber,evmLog.blockNumber)
	// fmt.Println("Here 2")
	blockNumber := bigNumber.String()

	var txIndex = txBigNumber - evmLog.txIndex
	var txHash = evmLog.txHash.String()

	if sender == receiver {
		// Transfer history
		var keySender = string(sender + "-" + blockNumber + "-" + strconv.Itoa(txIndex) + "-" + txHash)
		var valueSender = string(sender + "-" + receiver + "-" +  value + "-" + "2" + "-" + "0")
		// Transaction history
		var keyTransaction = string(txHash + "-" + sender + "-" + blockNumber + "-" + strconv.Itoa(txIndex))

		// Store history
		batch.Put([]byte(keyTransaction), []byte(keySender))
		batch.Put([]byte(keySender), []byte(valueSender))

		batch.Write()
		getValueSender, err := evmLogDb.customDb.Get([]byte(keySender))
		if err != nil{
			return false
		}
		fmt.Println("Transfer History of Sender", string(keySender), string(getValueSender))
	} else {
		// generate key and value for sender
		var keySender = string(sender + "-" + blockNumber + "-" + strconv.Itoa(txIndex) + "-" + txHash)
		var valueSender = string(sender + "-" + receiver + "-" +  value + "-" + "1" + "-" + "0")
		// Transaction history
		var keyTransactionSender = string(txHash + "-" + sender + "-" + blockNumber + "-" + strconv.Itoa(txIndex))
		
		batch.Put([]byte(keyTransactionSender), []byte(keySender))
		batch.Put([]byte(keySender), []byte(valueSender))
		
		// generate key and value for receiver
		var keyReceiver = string(receiver + "-" + blockNumber + "-" + strconv.Itoa(txIndex) + "-" + txHash)
		var valueReceiver = string(sender + "-" + receiver + "-" + value + "-" + "0" + "-" + "0")
		// Transaction history
		var keyTransactionReceiver = string(txHash + "-" + receiver + "-" + blockNumber + "-" + strconv.Itoa(txIndex))
		
		batch.Put([]byte(keyTransactionReceiver), []byte(keyReceiver))
		batch.Put([]byte(keyReceiver), []byte(valueReceiver))

		batch.Write()
		getValueSender, err := evmLogDb.customDb.Get([]byte(keySender))
		getValueReceiver, err := evmLogDb.customDb.Get([]byte(keyReceiver))
		getTransactionSender, err := evmLogDb.customDb.Get([]byte(keyTransactionSender))
		getTransactionReceiver, err := evmLogDb.customDb.Get([]byte(keyTransactionReceiver))
		if err != nil{
			fmt.Println("err in get transfer log", err)
			return false
		}
		// var address =  string(common.BytesToAddress(key))
		fmt.Println("Transfer History of Sender", string(keySender), string(getValueSender))
		fmt.Println("Transfer History of Receiver", string(keyReceiver), string(getValueReceiver))
		fmt.Println("Transfer History of Transaction Sender", string(keyTransactionSender), string(getTransactionSender))
		fmt.Println("Transfer History of Transaction Receiver", string(keyTransactionReceiver), string(getTransactionReceiver))
	}
	
	return true
}
/*
	History DB for Ethereum ERC20 Transaction
		address - token - 99million minor blockNumber - 99999 minor txIndex - 99 minor eventIndex - transactionHash	||	 from-to-value-flag(receiver=0 or sender=1 or self = 2)-flagIsRemoved(0:mainchian,1:oldchain,2:uncle) 	
	==========================================================================================================================================================================================================================
			0x123456-0xab123123-99999989-9998-89-0x789012															||   					0x123456-0xabcdef-10-1-0
	address: 0x123456
	token: 0xab123123 (ERC20 contract)
	blNumber: 10 => 99999989
	txIndex: 1 => 9998
	eventIndex: 10 => 89
	txHash: 0x789012

	from: 0x123456
	to: 0xabcdef
	value (amount of ERC20): 10
	sender: 0x123456 => flag = 1
	flagIsRemoved: 0 - mainchain
*/


/*
	History DB for Ethereum ERC20 Contract
		token - 100million minor blockNumber - 10000 minor txIndex - 100 minor eventIndex - transactionHash	||	 			from-to-value-flagIsRemoved(0:mainchian,1:oldchain,2:uncle)	
	======================================================================================================================================================================================
			0xab123123-99999989-9998-89-0x789012															||   							0x123456-0xabcdef-10-0
	token: 0xab123123 (ERC20 contract)
	blNumber: 10 => 99999990
	txIndex: 1 => 9999
	eventIndex: 10 => 90
	txHash: 0x789012

	from: 0x123456
	to: 0xabcdef
	value (amount of ERC20): 10
	flagIsRemoved: 0 - mainchain
*/
/*
	Transaction history format

	Key: prefix is transactionHash
	value: key of this transaction
*/
func (evmLogDb *EVMLogDb) storeEVMLogERC20(evmLog *EVMLog) (bool){
	// var blockNumberBigInt *big.Int
	// fmt.Println("Here 1")
	var bigNumber = big.NewInt(100000000)
	// fmt.Println("Init big int number", reflect.TypeOf(bigNumber), bigNumber)
	// fmt.Println("Get block number", reflect.TypeOf(evmLog.blockNumber), evmLog.blockNumber)
	bigNumber.Sub(bigNumber,evmLog.blockNumber)
	// fmt.Println("Here 2")
	blockNumber := bigNumber.String()

	var txIndex = txBigNumber - evmLog.txIndex
	var txHash = evmLog.txHash.String()

	var logs = evmLog.eventLog


	for i, log := range logs {
		// event Transfer(address indexed from, address indexed to, uint value)
		// Topic 0 is event signature
		// Topic 1 is from
		// Topic 2 is to
		topics := log.GetTopics()
		fmt.Println("Logs number", i, "has value: ", log, "and topics len is: ", len(topics))
		if len(topics) < 1 {
			return false
		}
		fmt.Println("Topic 0 of log ", i, "is: ", topics[0].String())

		if reflect.DeepEqual(topics[0],transferSig) && len(topics) >= 3 {
			batch := evmLogDb.customDb.NewBatch()
			contractAddress := log.GetAddress().String()
			sender := common.BytesToAddress(topics[1].Bytes()).String()
			receiver := common.BytesToAddress(topics[2].Bytes()).String()
			if common.IsHexAddress(sender) == false || common.IsHexAddress(receiver) == false {
				return false
			}
			z := new(big.Int)
			z.SetBytes(log.GetData())
			value := z.String()
			eventIndex := int(eventBigNumber - log.GetIndex())
			if sender == receiver {
				var keySender = string(sender + "-" + contractAddress + "-"+ blockNumber + "-" + strconv.Itoa(txIndex) + "-" + strconv.Itoa(eventIndex) + "-" + txHash)
				var valueSender = string(sender + "-" + receiver + "-" +  value + "-" + "2" + "-" + "0")

				var keyTransaction = string(txHash + "-" + sender + "-" + contractAddress + "-"+ blockNumber + "-" + strconv.Itoa(txIndex) + "-" + strconv.Itoa(eventIndex)) 
				batch.Put([]byte(keySender), []byte(valueSender))
				batch.Put([]byte(keyTransaction), []byte(keySender))
				batch.Write()
				getValueSender, err := evmLogDb.customDb.Get([]byte(keySender))
				if err != nil {
					fmt.Println("err in get erc20 log", err)
					return false
				}
				fmt.Println("Transfer ERC20 History of Sender", string(keySender), string(getValueSender))
			} else {
				var keySender = string(sender + "-" + contractAddress + "-"+ blockNumber + "-" + strconv.Itoa(txIndex) + "-" + strconv.Itoa(eventIndex) + "-" + txHash)
				var valueSender = string(sender + "-" + receiver + "-" +  value + "-" + "1" + "-" + "0")

				var keyTransactionSender = string(txHash + "-" + sender + "-" + contractAddress + "-"+ blockNumber + "-" + strconv.Itoa(txIndex) + "-" + strconv.Itoa(eventIndex))
				
				batch.Put([]byte(keySender), []byte(valueSender))
				batch.Put([]byte(keyTransactionSender), []byte(keySender))
				

				// generate key and value for receiver
				var keyReceiver = string(receiver + "-" + contractAddress + "-"+ blockNumber + "-" + strconv.Itoa(txIndex) + "-" + strconv.Itoa(eventIndex) + "-" + txHash)
				var valueReceiver = string(sender + "-" + receiver + "-" +  value + "-" + "0" + "-" + "0")
				
				var keyTransactionReceiver = string(txHash + "-" + receiver + "-" + contractAddress + "-"+ blockNumber + "-" + strconv.Itoa(txIndex) + "-" + strconv.Itoa(eventIndex))
				
				batch.Put([]byte(keyReceiver), []byte(valueReceiver))
				batch.Put([]byte(keyTransactionReceiver), []byte(keyReceiver))

				batch.Write()

				getValueSender, err := evmLogDb.customDb.Get([]byte(keySender))
				getValueReceiver, err := evmLogDb.customDb.Get([]byte(keyReceiver))
				if err != nil {
					fmt.Println("err in get erc20 log", err)
					return false
				}	
				fmt.Println("Transfer ERC20 History of Sender", string(keySender), string(getValueSender))
				fmt.Println("Transfer ERC20 History of Receiver", string(keyReceiver), string(getValueReceiver))
			}
			// generate key and value for token
			var keyToken = string(contractAddress + "-"+ blockNumber + "-" + strconv.Itoa(txIndex) + "-" + strconv.Itoa(eventIndex) + "-" + txHash)
			var valueToken = string(sender + "-" + receiver + "-" +  value + "-" + "0")

			var keyTransaction = string(txHash + "-" + contractAddress + "-"+ blockNumber + "-" + strconv.Itoa(txIndex) + "-" + strconv.Itoa(eventIndex)) 
			batch.Put([]byte(keyToken), []byte(valueToken))
			batch.Put([]byte(keyTransaction), []byte(keyToken))

			batch.Write()
			getValueToken, err := evmLogDb.customDb.Get([]byte(keyToken))
			if err != nil {
				fmt.Println("err in get erc20 log", err)
				return false
			}		
			fmt.Println("Transfer ERC20 History of Token", string(keyToken), string(getValueToken))
		}
	}
	return true
}


func (evmLogDb *EVMLogDb) MarkForkedOrUncleTransaction(txs types.Transactions, mark string) (bool){
	batch := evmLogDb.customDb.NewBatch()
	for i, tx := range txs {
		fmt.Println("Transaction", i, "is marked as forked: ", tx.Hash().String())
		var hash = tx.Hash().String()
		ldbDatabase := evmLogDb.customDb.(*ethdb.LDBDatabase)
		fmt.Println("transaction to get history: ", hash)
		iter := ldbDatabase.NewIterator()
		// Seek all the key relevant to this transaction and marked as forked
		for ok := iter.Seek([]byte(hash)); ok && strings.HasPrefix(string(iter.Key()), hash); ok = iter.Next(){
			// keyTx := iter.Key()
			valueTx := iter.Value()
			// valueTx is key to get real value
			value,_ := evmLogDb.customDb.Get(valueTx)
			valueString := string(value)
			changedValue := valueString[:len(valueString)-1] + mark
			batch.Put(valueTx, []byte(changedValue))
			fmt.Println("Forked Transaction", string(valueTx), "old value:", string(value), "new value:", string(changedValue))
		}
	}
	batch.Write()
	return true
}


// Test DB when init ETH API
func (evmLogDb *EVMLogDb) TestDb() (bool){
	key := []byte{'k','e','y'}
	value := []byte{'v','a','l','u','e'}
	evmLogDb.customDb.Put(key,value)
	result, err := evmLogDb.customDb.Get(key)
	if err != nil{
		return false
	}
	fmt.Println("TestDb", string(result))
	return true
}