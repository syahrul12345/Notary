const Web3 = require('Web3')
web3 = new Web3(new Web3.providers.HttpProvider("https://ropsten.infura.io/v3/2a3f078d3755444b8777a0204e5f694a"))

async function main() {
    accountInfo = await web3.eth.accounts.create();
    const payload = {
        from:     "0xE9C0614F054FAd022e989034c00b136E507e162b",
        gas:      100000,
        gasPrice: '50000000000',
        value:    0,
        data:      "889072148fd9e7d1fef808c270420baf2dce16c44f0335cf4c4b8fbfa3891578"
    }
    const tx = await web3.eth.accounts.signTransaction(payload,"0x600afd241a7e7a0e36b0267c6ac0bd6aa9396b338b814cba3ce4be83c10bbce2")
    console.log(tx.rawTransaction)
    console.log("\n")
    web3.eth.sendSignedTransaction(tx.rawTransaction)
    .on("transactionHash",console.log)
    .on('receipt',console.log)
    
}
main()

