/*
User interface for Dero multisig wallet by thedudelebowski. 
Version 1.0

Please note: This UI was written for the Dero Stargate testnet. Use at your own risk! The code could be simplified in some areas, and error handling may not be 100% complete. 

Github link: https://github.com/lebowski1234/multisig-ui 
*/

package main

import (
	"bufio"	
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"net/http"	
	"strconv"	
	"github.com/dixonwille/wmenu"
	"github.com/tidwall/gjson"
)

//For menu
type menuItem int


//For daemon call
type PayloadKeys struct { 
	TxsHashes []string `json:"txs_hashes"` 
	ScKeys []string `json:"sc_keys"` 
}



//Wallet call function with just value (send Dero)
type Params struct {
	Mixin    int  `json:"mixin"`
	GetTxKey bool `json:"get_tx_key"`
	ScTx	 ScTx `json:"sc_tx"`
}

type ScTx struct {
	Entrypoint string `json:"entrypoint"`
	Scid       string `json:"scid"`
	Value      int64  `json:"value"`
} 

type PayloadDeposit struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      string `json:"id"`
	Method  string `json:"method"`
	Params  Params `json:"params"`
}


//Wallet call function with multiple strings and value
type Params3 struct {
	To	string `json:"To,omitempty"`
	Amount	string `json:"Amount,omitempty"`
	ID	string `json:"ID,omitempty"`
}

type Params2 struct {
	Mixin    int  `json:"mixin"`
	GetTxKey bool `json:"get_tx_key"`
	ScTx	 ScTx2 `json:"sc_tx"`
}

type ScTx2 struct {
	Entrypoint string `json:"entrypoint"`
	Scid       string `json:"scid"`
	Value      int64  `json:"value"`
	Params	   Params3 `json:"params"`
	
} 

type PayloadGeneral struct {
	Jsonrpc string `json:"jsonrpc"`
	ID      string `json:"id"`
	Method  string `json:"method"`
	Params  Params2 `json:"params"`
}


//For menu
const (
	enterSCID menuItem = iota
	viewBalance
	depositDero
	send
	sign
	displayUnsigned
	displaySigned
	displayAll
	displayID
	exit
)

//For menu
var menuItemStrings = map[menuItem]string{
	enterSCID:	"Enter smart contract ID (SCID)",
	viewBalance:	"View wallet balance",
	depositDero:   	"Deposit Dero in wallet",
	send:		"Create a new transaction",
	sign:		"Sign a transaction",
	displayUnsigned:	"Display unsent (open) transactions",
	displaySigned:	"Display signed (sent) transactions",
	displayAll:	"Display all transactions",
	displayID:	"Display a transaction by Index No.",
	exit:		"Exit",		
}


var SCID string


func main() {
	mm := mainMenu()
	err := mm.Run()
	if err != nil {
		wmenu.Clear()		
		fmt.Println(err)
		mm.Run()		
		
	}


	
}


/*-----------------------------------------------------------Menu-----------------------------------------------------------------*/


func mainMenu() *wmenu.Menu {
	menu := wmenu.NewMenu("Dero multisig wallet: choose an option.")
	menu.Option(menuItemStrings[enterSCID], enterSCID, false, func(opt wmenu.Opt) error { //change false to true to make default option
		wmenu.Clear()		
		getSCID() 	
		mm := mainMenu()
		return mm.Run()
	})

	menu.Option(menuItemStrings[viewBalance], viewBalance, false, func(opt wmenu.Opt) error {
		wmenu.Clear()
		if SCID != "" {		
			displayBalance(SCID)		
		} else {
			fmt.Println("Please enter a SCID (Menu Option 1)\n")
		}			
		mm := mainMenu()
		return mm.Run()
	})

	menu.Option(menuItemStrings[depositDero], depositDero, false, func(opt wmenu.Opt) error {
		wmenu.Clear()		
		amount, proceed:=getDepositAmount()
		if proceed == false {
			fmt.Println("Transaction cancelled")
		} else {
			if SCID != "" {		
			deposit(SCID, amount)		
			} else {
				fmt.Println("Please enter a SCID (Menu Option 1)\n")
			}	
		}		
		mm := mainMenu()
		return mm.Run()
	})

	menu.Option(menuItemStrings[send], send, false, func(opt wmenu.Opt) error {
		wmenu.Clear()		
		recipient, amount, proceed:= getSendParams()
		if proceed == false {
			fmt.Println("Transaction cancelled")
		} else {
			if SCID != "" {		
			sendTransaction(SCID, "Send", recipient, amount, "")		
			} else {
				fmt.Println("Please enter a SCID (Menu Option 1)\n")
			}	
		}				
		mm := mainMenu()
		return mm.Run()
	})

	menu.Option(menuItemStrings[sign], sign, false, func(opt wmenu.Opt) error {
		wmenu.Clear()		
		id, proceed:= getSignParams()
		if proceed == false {
			fmt.Println("Transaction cancelled")
		} else {
			if SCID != "" {		
			sendTransaction(SCID, "Sign", "", 0, id)		
			} else {
				fmt.Println("Please enter a SCID (Menu Option 1)\n")
			}	
		}				
		mm := mainMenu()
		return mm.Run()
	})

	menu.Option(menuItemStrings[displayUnsigned], displayUnsigned, false, func(opt wmenu.Opt) error {
		wmenu.Clear()		
		displayTransactions(SCID, 1, "")	
		pressToContinue()	
		mm := mainMenu()
		return mm.Run()
	})

	menu.Option(menuItemStrings[displaySigned], displaySigned, false, func(opt wmenu.Opt) error {
		wmenu.Clear()		
		displayTransactions(SCID, 2, "")
		pressToContinue()		
		mm := mainMenu()
		return mm.Run()
	})

	menu.Option(menuItemStrings[displayAll], displayAll, false, func(opt wmenu.Opt) error {
		wmenu.Clear()		
		displayTransactions(SCID, 0, "")
		pressToContinue()		
		mm := mainMenu()
		return mm.Run()
	})

	menu.Option(menuItemStrings[displayID], displayID, false, func(opt wmenu.Opt) error {
		wmenu.Clear()	
		txIDno:=getID()	
		displayByID(SCID, txIDno)		
		mm := mainMenu()
		return mm.Run()
	})

	menu.Option(menuItemStrings[exit], exit, false, func(opt wmenu.Opt) error {
		wmenu.Clear()		
		return nil //Exit		
		
	})
	
	menu.Action(func(opts []wmenu.Opt) error {
		if len(opts) != 1 {
			return errors.New("wrong number of options chosen")
		}
		wmenu.Clear()
		mm := mainMenu()
		return mm.Run()
		
	})
	return menu
}



//Get SCID, save to memory
func getSCID() {
	scanner := bufio.NewScanner(os.Stdin)
	var text string
	fmt.Print("Enter SCID: ")
	scanner.Scan()
	text = scanner.Text()
	wmenu.Clear()	
	SCID = text
	fmt.Println("SCID entered: ", text)
	fmt.Print("Press 'Enter' to continue...")
  	bufio.NewReader(os.Stdin).ReadBytes('\n')
      
}



//Get tx ID to display
func getID() string {
	scanner := bufio.NewScanner(os.Stdin)
	var text string
	fmt.Print("Enter Transaction ID No: ")
	scanner.Scan()
	text = scanner.Text()
	wmenu.Clear()	
	fmt.Println("Transaction ID No entered: ", text)
	fmt.Print("Press 'Enter' to continue...")
  	bufio.NewReader(os.Stdin).ReadBytes('\n')
	return text
      
}



func pressToContinue() {
	fmt.Print("Press 'Enter' to continue...")
  	bufio.NewReader(os.Stdin).ReadBytes('\n')
	//wmenu.Clear()
      
}

//Enter deposit amount, return value. 
func getDepositAmount() (int64, bool) {
	scanner := bufio.NewScanner(os.Stdin)
	var amountString string
	fmt.Print("Enter deposit amount in Dero: ")
	scanner.Scan()
	amountString = scanner.Text()
	wmenu.Clear()	
	fmt.Printf("Do you want to deposit %s Dero? Enter Y/N (Yes/No)", amountString)
	confirmed:=askForConfirmation()
	if confirmed == true {
		amountFloat, err := strconv.ParseFloat(amountString, 64) //convert to float64 
		if err != nil {
			fmt.Println(err)
			return 0, false
		}
		amount:= int64(amountFloat * 1000000000000)
		return amount, true 		

	} else {
		return 0, false
	}
		
      
}



//Enter recipient address and amount, return values. 
func getSendParams() (string, int64, bool) {
	scanner := bufio.NewScanner(os.Stdin)
	var recipient string
	var amountString string

	fmt.Print("Enter recipient address: ")
	scanner.Scan()
	recipient = scanner.Text()
	wmenu.Clear()	
		
	fmt.Print("Enter deposit amount in Dero: ")
	scanner.Scan()
	amountString = scanner.Text()
	wmenu.Clear()	
	fmt.Printf("Do you want to send %s Dero to %s? Enter Y/N (Yes/No)", amountString, recipient)
	confirmed:=askForConfirmation()
	if confirmed == true {
		amountFloat, err := strconv.ParseFloat(amountString, 64) //convert to float64 
		if err != nil {
			fmt.Println(err)
			return "", 0, false
		}
		amount:= int64(amountFloat * 1000000000000)
		return recipient, amount, true 		

	} else {
		return "", 0, false
	}
		
      
}



//Enter transaction ID, return value. 
func getSignParams() (string, bool) {
	scanner := bufio.NewScanner(os.Stdin)
	var id string
	fmt.Print("Enter transaction ID: ")
	scanner.Scan()
	id = scanner.Text()
	wmenu.Clear()	
	fmt.Printf("Do you want to sign transaction %s? Enter Y/N (Yes/No)", id)
	confirmed:=askForConfirmation()
	if confirmed == true {
		return id, true 		

	} else {
		return "", false
	}
		
      
}


// The following 3 functions were taken directly from https://gist.github.com/albrow/5882501

// askForConfirmation uses Scanln to parse user input. A user must type in "yes" or "no" and
// then press enter. It has fuzzy matching, so "y", "Y", "yes", "YES", and "Yes" all count as
// confirmations. If the input is not recognized, it will ask again. The function does not return
// until it gets a valid response from the user. Typically, you should use fmt to print out a question
// before calling askForConfirmation. E.g. fmt.Println("WARNING: Are you sure? (yes/no)")
func askForConfirmation() bool {
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		fmt.Println("Error")
	}
	okayResponses := []string{"y", "Y", "yes", "Yes", "YES"}
	nokayResponses := []string{"n", "N", "no", "No", "NO"}
	if containsString(okayResponses, response) {
		return true
	} else if containsString(nokayResponses, response) {
		return false
	} else {
		fmt.Println("Please type yes or no and then press enter:")
		return askForConfirmation()
	}
}

// posString returns the first index of element in slice.
// If slice does not contain element, returns -1.
func posString(slice []string, element string) int {
	for index, elem := range slice {
		if elem == element {
			return index
		}
	}
	return -1
}

// containsString returns true if slice contains element
func containsString(slice []string, element string) bool {
	return !(posString(slice, element) == -1)
}





/*-----------------------------------------------------------RPC Functions-----------------------------------------------------------------*/



//sendTransaction: send a transaction to the wallet or sign a transaction. entry should be "Send" or "Sign". 
func sendTransaction(scid string, entry string, to string, amount int64, id string) {
	
	walletURL:= "http://127.0.0.1:30309/json_rpc"
	var amountString string	
	if amount == 0 {
		amountString = ""
	} else {	
		amountString = strconv.FormatInt(amount, 10)
	}
	data:=  PayloadGeneral{
		Jsonrpc: "2.0", 
		ID: "0",
		Method: "transfer_split",
		Params: Params2{
			Mixin: 5,
			GetTxKey: true,
			ScTx: ScTx2{
				Entrypoint: entry,
				Scid: scid,
				Value: 0,
				Params: Params3{
						To: to,
						Amount: amountString,
						ID: id,
				},
			}, 
		},
	}

	
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	body := bytes.NewReader(payloadBytes)
	
	_, err=rpcPost(body, walletURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	
	//println(result)	
	fmt.Println("Transaction sent to wallet!")
	
}



//deposit: Deposit Dero to SC
func deposit(scid string, amount int64) {
	
	walletURL:= "http://127.0.0.1:30309/json_rpc"
	
	data:=  PayloadDeposit{
		Jsonrpc: "2.0", 
		ID: "0",
		Method: "transfer_split",
		Params: Params{
			Mixin: 5,
			GetTxKey: true,
			ScTx: ScTx{
				Entrypoint: "Deposit", 
				Scid: scid,
				Value: amount,
			}, 
		},
	}

	
	payloadBytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	body := bytes.NewReader(payloadBytes)
	
	_, err=rpcPost(body, walletURL)
	if err != nil {
		fmt.Println(err)
		return
	}
	
	//println(result)	
	fmt.Println("Deposit sent to wallet!")
	
}



//getKeysFromDaemon: send RPC call with list of keys, do error checking, return raw data in string form for JSON extraction
func getKeysFromDaemon(scid string, keys []string) string {
	
	deamonURL:= "http://127.0.0.1:30306/gettransactions"
	txHashes:= []string{scid}
	
	data := PayloadKeys{
		TxsHashes: txHashes,
		ScKeys: keys, 
	}

	payloadBytes, err := json.Marshal(data)
	if err != nil {
		println("Error in function getKeysFromDaemon:")
		fmt.Println(err)
		return ""
	}
	body := bytes.NewReader(payloadBytes)
	
	result, err:=rpcPost(body, deamonURL)
	if err != nil {
		println("Error in function getKeysFromDaemon:")
		fmt.Println(err)
		return ""
	}
	
	//Check to see if we have got the expected response from the daemon:
	if !gjson.Valid(result) {
		println("Error, result not in JSON format")
		return ""
	}
	
	validResponse := gjson.Get(result, "txs_as_hex") 	
	//fmt.Printf("Array value is: %s", validResponse.Array()[0])
	
	validResponse0:= validResponse.Array()[0]

	if validResponse0.String() == "" { //array position 0 value will be [""] if SCID not found or invalid
		println("Error, SCID not found")
		return ""
	}
	
	
	return result
	
			
}


//displayBalance: Display SC balance from daemon
func displayBalance(scid string) {
	
	scKeys:= []string{"null"} //Don't get any keys, just sc_balance which is returned by default

	result:= getKeysFromDaemon(scid, scKeys)
			
	if result == "" {return}
	
	//Response ok, extract keys from JSON
	balance := gjson.Get(result, "txs.#.sc_balance")
	value:= balance.Array()[0]
	intValue:= value.Uint()
		
	f:= float64(intValue)/1000000000000
	fmt.Printf("Wallet balance is %6.12f Dero \n\n", f)

	pressToContinue()
				
}





//displayByID: displays transaction data stored in smart contract, for one TX only. 
func displayByID(scid string, ID string) {
		
	scKeys:= []string{"numberOfOwners", "txCount"}
	result:= getKeysFromDaemon(scid, scKeys)
	if result == "" {return}


	//Response ok, extract keys from JSON
	
	
	numberOfOwners := gjson.Get(result, "txs.#.sc_keys.numberOfOwners")
	numberOfOwnersArray:= numberOfOwners.Array()[0]
	numberOfOwnersInt:= numberOfOwnersArray.Int()
	//fmt.Printf("Number of owners: %d\n", numberOfOwnersInt)

	//make sure ID No. exists
	txCount := gjson.Get(result, "txs.#.sc_keys.txCount")
	txCountArray:= txCount.Array()[0]
	txCountInt:= txCountArray.Int()
	IDint, err := strconv.ParseInt(ID, 10, 64) //convert to int64 
	if err != nil {
		fmt.Println(err)
		return 
	}

	if IDint >= txCountInt {
		fmt.Println("Error, transaction ID does not exist")
		return
	}


	//fmt.Printf("Tx Count: %d\n", txCountInt)


	//Make a slice of keys so we can request in RPC call
	x:= int(numberOfOwnersInt) 	
	keySlice:= make([]string, x) 
	
	for i:=0; i<x; i++ {
		z:= strconv.Itoa(i+1) //number of owners starts at 1, not 0
		keySlice[i] = "tx" + ID + "_signer" + z
		
	}
		
	//fmt.Println(keySlice)
	
	//request keys
	result= getKeysFromDaemon(scid, keySlice)
			
	if result == "" {return}

	displayTransactions(scid, 3, ID) //display transaction data
	
	//display which owners have signed
	for i:=0; i<x; i++ {
		z:= strconv.Itoa(i+1)
		
		keyName:= "txs.#.sc_keys.tx" + ID + "_signer" + z
		raw:= gjson.Get(result, keyName)
		resultString:=raw.String() 
		if resultString != "[]" { //weird way of checking value exists, raw.Exists() doesn't seem to work in this context
			//signer:=raw.Array()[0].String()
			resultInt:=raw.Array()[0].Int()
			if resultInt == 1 {
				fmt.Printf("Owner %s status: Signed \n", z)
			} else {
				fmt.Printf("Owner %s status: Unsigned \n", z)
			}			
			
		}

						
	}
	fmt.Printf("\n") //print blank line

	pressToContinue()

}




//displayTransactions: displays transaction data stored in smart contract. Option 1 = unsigned, 2 = signed, 3 = one ID only, default (e.g. 0) = all.
func displayTransactions(scid string, option int, ID string) {
		
	scKeys:= []string{"numberOfOwners", "txCount"}
	result:= getKeysFromDaemon(scid, scKeys)
	if result == "" {return}


	//Response ok, extract keys from JSON
	

	txCount := gjson.Get(result, "txs.#.sc_keys.txCount")
	txCountArray:= txCount.Array()[0]
	txCountInt:= txCountArray.Int()
	//fmt.Printf("Tx Count: %d\n", txCountInt)

	//Make a slice of keys so we can request in RPC call
	x:= int(txCountInt) //txCount in wallet smart contract is always 1 ahead of actual number of transactions	
	x4:= x * 4	
	keySlice:= make([]string, x4) 
	
	for i:=0; i<x; i++ {
		z:= strconv.Itoa(i)
		keySlice[i] = "txIndex_" + z
		keySlice[i+x] = "recipient_" + z
		keySlice[i+(x*2)] = "amount_" + z
		keySlice[i+(x*3)] = "sent_" + z
	}
		
	//fmt.Println(keySlice)
	displayData(scid, keySlice, x, option, ID)


}


//displayData: called from displayTransactions: displays transaction data stored in smart contract. 
func displayData(scid string, keys []string, x int, option int, ID string) {
	
	result:= getKeysFromDaemon(scid, keys)
			
	if result == "" {return}

	switch option {
		case 1: //Show all unsent transactions
			fmt.Printf("Listing all unsent transactions:\n")

		case 2: //Show all sent transactions
			fmt.Printf("Listing all sent transactions:\n")

		case 3: //Show only one transaction, by ID
			//fmt.Printf("Listing transaction with Index %s:\n", ID)
			fmt.Printf("Listing transaction data:\n")
									
		default: //Show all transactions
			fmt.Printf("Listing all transactions:\n")
						
			
		}

	
	for i:=0; i<x; i++ {
		z:= strconv.Itoa(i)
		
		var txID string
		var recipient string
		var amount int64
		var sent bool
		var amountDero float64
		
		//get recipient		
		keyName:= "txs.#.sc_keys.txIndex_"+z
		raw:= gjson.Get(result, keyName)
		resultString:=raw.String() 
		if resultString != "[]" { //weird way of checking value exists, raw.Exists() doesn't seem to work in this context
			txID=raw.Array()[0].String()
			//fmt.Printf("TX ID %s: ", txID)
			
		}

		//get recipient		
		keyName= "txs.#.sc_keys.recipient_"+z
		raw= gjson.Get(result, keyName)
		resultString=raw.String() 
		if resultString != "[]" { //weird way of checking value exists, raw.Exists() doesn't seem to work in this context
			recipient=raw.Array()[0].String()
			//fmt.Printf("Recipient = %s ", recipient)
			
		}

		//get amount 		
		keyName= "txs.#.sc_keys.amount_"+z
		raw= gjson.Get(result, keyName)
		resultString=raw.String() 
		if resultString != "[]" { //weird way of checking value exists, raw.Exists() doesn't seem to work in this context
			amount=raw.Array()[0].Int()
			amountDero= float64(amount)/1000000000000
			//fmt.Printf("Amount = %6.12f Dero ", amountDero)
			
		}

		//get sent 		
		keyName= "txs.#.sc_keys.sent_"+z
		raw= gjson.Get(result, keyName)
		resultString=raw.String() 
		if resultString != "[]" { //weird way of checking value exists, raw.Exists() doesn't seem to work in this context
			sentInt:=raw.Array()[0].Int()
			if sentInt == 1 {
				sent = true
			} else {
				sent = false
			}
			//fmt.Printf("Sent = %t \n", sent)
			
		}


		switch option {
		case 1: //Show all unsent transactions
			if sent == false {
				fmt.Printf("TX Index %s:\n", txID)
				fmt.Printf("Recipient = %s\n", recipient)
				fmt.Printf("Amount = %6.12f Dero\n\n", amountDero)
				//fmt.Printf("Sent = %t ", sent)
				//fmt.Printf("Failed = %t \n", failed)
				
			}

		case 2: //Show all sent transactions
			if sent == true {
				fmt.Printf("TX Index %s:\n", txID)
				fmt.Printf("Recipient = %s\n", recipient)
				fmt.Printf("Amount = %6.12f Dero\n\n", amountDero)
				//fmt.Printf("Sent = %t \n", sent)
				
			}

		case 3: //Show only one transaction, by ID
			if z == ID {
				fmt.Printf("TX Index %s:\n", txID)
				fmt.Printf("Recipient = %s\n", recipient)
				fmt.Printf("Amount = %6.12f Dero\n", amountDero)
				fmt.Printf("Sent = %t\n\n", sent)
				//fmt.Printf("Failed = %t\n", failed)
				
			}
		
				
		default: //Show all transactions
			fmt.Printf("TX Index %s:\n", txID)
			fmt.Printf("Recipient = %s\n", recipient)
			fmt.Printf("Amount = %6.12f Dero\n", amountDero)
			fmt.Printf("Sent = %t\n\n", sent)
			//fmt.Printf("Failed = %t \n", failed)
						
			
		}
				
	}

	//pressToContinue()

}



//rpcPost: Send RPC request, return response body as string 
func rpcPost(body *bytes.Reader, url string) (string, error) {

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		println("Error in function rpcPost:")
		fmt.Println(err)		
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")


	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		println("Error in function rpcPost:")
		fmt.Println(err)
		return "", err
	}
	defer resp.Body.Close()

	response, err:= ioutil.ReadAll(resp.Body)
	result:=string(response)

	return result, err
}



