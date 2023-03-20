# blockexplorer

A blockexplorer module written in golang

## General usage

Note: take a look at the "private repo notes" at the bottom of this file before proceeding.

To use this module, cd to your project and run the following command:

`go get code.cryptopower.dev/group/instantswap/blockexplorer`

import the module in the file(s) needing the module:

`import "code.cryptopower.dev/group/instantswap/blockexplorer"`

instantiate a new blockexplorer:

```
explorer, err := blockexplorer.NewExplorer("BTC", false)
if err != nil {
    return nil, err
}

```

### Available Methods

verify a tx based on the values passed in to the request params:

```
verificationInfo := blockexplorer.TxVerifyRequest{}
verification, err := explorer.VerifyTransaction(verificationInfo)
if err != nil {
    return nil, err
}
```

push a raw tx:

```
resp, err := explorer.PushTx(txHash)
if err != nil {
    return nil, err
}
```

get a transcation using the txID:

```
resp, err := explorer.GetTransaction(txID)
if err != nil {
    return nil, err
}
```

get all transactions belonging to a particular address:

```
resp, err := explorer.GetTxsForAddress(address, limit, viewKey)
if err != nil {
    return nil, err
}
```

## Private Repo Notes

In order to use this repo you will need to configure git to use ssh instead of https:

create ~/.gitconfig:
```
[user]
    name = Nane
    email = some@email.address
[url "git@code.cryptopower.dev:"]
	insteadOf = https://code.cryptopower.dev/
```

create ~/.netrc:
```
machine code.cryptopower.dev
login <current shared auth username>
password <current shared auth password>
```
For `go get` commands to work you will need to set the `GOPRIVATE` variable. 
Example:
`export GOPRIVATE="code.cryptopower.dev/"`