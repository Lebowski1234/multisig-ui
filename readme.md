# User Interface for Multi-Signature Wallet

This is the user interface (UI) for the multisig wallet written for the Dero Stargate Smart Contract competition, located in this depository: [https://github.com/lebowski1234/dero-multisig](https://github.com/lebowski1234/dero-multisig). This document contains basic usage instructions.

Binaries for Windows and Linux (both 64 bit) are located here:

* Linux: [https://github.com/Lebowski1234/multisig-ui/raw/master/binaries/multisig-linux64.tar.gz](https://github.com/Lebowski1234/multisig-ui/raw/master/binaries/multisig-linux64.tar.gz)

* Windows: [https://github.com/Lebowski1234/multisig-ui/raw/master/binaries/multisig-windows64.rar](https://github.com/Lebowski1234/multisig-ui/raw/master/binaries/multisig-windows64.rar)

Or follow the instructions below to compile. 


## Compiling

All development was done in Ubuntu using Go version 1.11.4.

First, download dependencies:

```
$ go get -u github.com/tidwall/gjson
$ go get -u github.com/dixonwille/wmenu
```

Then build:

```
$ go build multisig.go
```


## Running

The Dero Stargate daemon and wallet must both be running first, with the standard RPC ports open. The wallet must be unlocked and have minimum 50 Dero unlocked balance available.

Get the Dero Stargate binaries here:

[https://git.dero.io/DeroProject/Dero_Stargate_testnet_binaries](https://git.dero.io/DeroProject/Dero_Stargate_testnet_binaries)


To run the Dero Stargate daemon (in Linux):

```
./derod-linux-amd64 --testnet
```

To run the Dero Stargate wallet:

```
./dero-wallet-cli-linux-amd64 --rpc-server --wallet-file testnetwallet.db --testnet
```

Finally, to run the multisig wallet user interface:

```
$ ./multisig
```

The instructions are the same for Windows, without the './'


## Usage

Refer to the multisig wallet smart contract [readme](https://github.com/lebowski1234/dero-multisig) for an explanation of how the wallet works. All options in the user interface are self explanatory and intuitive. After the smart contract has been deployed, run the user interface and choose from Options 1 to 10 (e.g. type '1' then enter):

### Option 1 - Enter Smart Contract ID (SCID)

Enter the smart contract ID (SCID), which was displayed in the daemon after deploying the smart contract. This must be done before any other options are selected. 

### Option 2 - View Wallet Balance

Balance held in smart contract is displayed on screen.

### Option 3 - Deposit Dero in Wallet

Enter the amount in Dero, up to 12 decimal places.

### Option 4 - Create a New Transaction

This option creates a new transaction within the smart contract. Enter the recipient Dero address, and the amount to send. The transaction will be stored in the wallet, but no Dero will be sent until the transaction has been signed by the required number of owners. The recipient must be a valid Dero address, and the value must be greater than 0. 

### Option 5 - Sign a Transaction

Enter the transaction index No. To get the transaction index No, use option 6 to list all unsent (open) transactions on screen.

### Option 6 - Display Unsent (Open) Transactions

This option lists all unsent (open) transactions on screen, showing index No, recipient, and value. 

### Option 7 - Display Signed (Sent) Transactions

Similar to option 6, but for sent (closed) transactions.

### Option 8 - Display All Transactions

Lists all transactions on screen, showing which are sent, and which are unsent. 

### Option 9 - Display Transaction by Index No.

Enter the index No. of a transaction to display recipient, value, sent / unsent status, and a list of which owners have signed / have not signed this transaction. 

### Option 10 - Exit

Exit the user interface. 


