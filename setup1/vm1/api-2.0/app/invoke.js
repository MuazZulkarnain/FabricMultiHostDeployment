const { Gateway, Wallets, TxEventHandler, GatewayOptions, DefaultEventHandlerStrategies, TxEventHandlerFactory } = require('fabric-network');
const fs = require('fs');
const path = require("path")
const log4js = require('log4js');
const logger = log4js.getLogger('BasicNetwork');
const util = require('util')

// const createTransactionEventHandler = require('./MyTransactionEventHandler.ts')

const helper = require('./helper')

// const createTransactionEventHandler = (transactionId, network) => {
//     /* Your implementation here */
//     const mspId = network.getGateway().getIdentity().mspId;
//     const myOrgPeers = network.getChannel().getEndorsers(mspId);
//     return new MyTransactionEventHandler(transactionId, network, myOrgPeers);
// }

const invokeTransaction = async (channelName, chaincodeName, fcn, args, username, org_name) => {
    try {
        logger.debug(util.format('\n============ invoke transaction on channel %s ============\n', channelName));

        // load the network configuration
        // const ccpPath =path.resolve(__dirname, '..', 'config', 'connection-org1.json');
        // const ccpJSON = fs.readFileSync(ccpPath, 'utf8')
        const ccp = await helper.getCCP(org_name) //JSON.parse(ccpJSON);

        // Create a new file system based wallet for managing identities.
        const walletPath = await helper.getWalletPath(org_name) //path.join(process.cwd(), 'wallet');
        const wallet = await Wallets.newFileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        let identity = await wallet.get(username);
        if (!identity) {
            console.log(`An identity for the user ${username} does not exist in the wallet, so registering user`);
            await helper.getRegisteredUser(username, org_name, true)
            identity = await wallet.get(username);
            console.log('Run the registerUser.js application before retrying');
            return;
        }



        const connectOptions = {
            wallet, identity: username, discovery: { enabled: true, asLocalhost: false },
            eventHandlerOptions: {
                commitTimeout: 100,
                strategy: DefaultEventHandlerStrategies.NETWORK_SCOPE_ALLFORTX
            }
            // transaction: {
            //     strategy: createTransactionEventhandler()
            // }
        }

        // Create a new gateway for connecting to our peer node.
        const gateway = new Gateway();
        await gateway.connect(ccp, connectOptions);

        const network = await gateway.getNetwork(channelName);
        const contract = network.getContract(chaincodeName);

        let result;
        let message;
        let tokenId = args[0];
        let tokenExist = await contract.evaluateTransaction("queryToken", tokenId);
      

        if (fcn === "generateToken") {
          if (tokenExist.toString()) {
              // assign variables
            let currentOwner = args[0];
            let amount = parseFloat(args[1]);
          
            // Check if token exists
            tokenExist = await contract.evaluateTransaction("queryToken", currentOwner);
            if (tokenExist.toString()) {
            } else {
              return `Token with key ${currentOwner} does not exist`;
            }

            let currentOwnerToken = JSON.parse(tokenExist.toString());
            let currentOwnerAmount = currentOwnerToken.amount;

            result = await contract.submitTransaction("updateTokenVolume", currentOwner, currentOwnerAmount + amount, "Token Generated");
            result = JSON.parse(result.toString());
            message = `Successfully generated ${args[1]} tokens for ${currentOwner}!`;
          } else {
              result = await contract.submitTransaction("createToken", tokenId, args[1], args[2], args[3], args[4], "Token Generated");
              result = JSON.parse(result.toString());
              message = `Successfully generated ${args[1]} tokens for ${tokenId}!`;
          }
        } else if (fcn === "transferToken") {
            // assign variables
            let currentOwner = args[0];
            let newOwner = args[1];
            let amount = parseFloat(args[2]);
          
            // Check if token exists
            tokenExist = await contract.evaluateTransaction("queryToken", currentOwner);
            if (tokenExist.toString()) {
            } else {
              return `Token with key ${currentOwner} does not exist`;
            }

            let currentOwnerToken = JSON.parse(tokenExist.toString());
            let currentOwnerAmount = currentOwnerToken.amount;

            tokenExist = await contract.evaluateTransaction("queryToken", newOwner);
            if (tokenExist.toString()) {
            } else {
              return `Token with key ${newOwner} does not exist`;
            }
          
            let newOwnerToken = JSON.parse(tokenExist.toString());
            let newOwnerAmount = newOwnerToken.amount;

            if (currentOwnerAmount < amount ) {
              return `${currentOwner} do not have enough tokens for the transfer!`;
            }

            // Transfer the tokens
            result = await contract.submitTransaction("updateTokenVolume", currentOwner, currentOwnerAmount - amount, "Token Transferred");
            result = await contract.submitTransaction("updateTokenVolume", newOwner, newOwnerAmount + amount, "Token Transferred");
            result = JSON.parse(result.toString());
            message = `Successfully transferred ${amount} tokens to ${newOwner}`;
        } else if (fcn === "createToken") {
            result = await contract.submitTransaction(fcn, args[0], args[1], args[2], args[3], args[4]);
            result = JSON.parse(result.toString());
            message = `Successfully added the token asset with key ${args[0]}`

        } else if (fcn === "changeTokenOwner") {
            result = await contract.submitTransaction(fcn, args[0], args[1]);
            result = JSON.parse(result.toString());
            message = `Successfully changed token owner with key ${args[0]}`

        } else if (fcn === "updateTokenVolume") {
            result = await contract.submitTransaction(fcn, args[0], args[1]);
            result = JSON.parse(result.toString());
            message = `Successfully updated token volume with key ${args[0]}`

        } else if (fcn === "retireToken") {
            result = await contract.submitTransaction(fcn, args[0]);
            result = JSON.parse(result.toString());
            message = `Successfully retired token with key ${args[0]}`

        } else if (fcn === "retirePartialToken") {
            tokenExist = await contract.evaluateTransaction("queryToken", tokenId);
            if (tokenExist.toString()) {
            } else {
              return `Token with key ${tokenId} does not exist`;
            }

            let amount = parseFloat(args[1]);
            let currentToken = JSON.parse(tokenExist.toString());
            let currentAmount = currentToken.amount;
            result = await contract.submitTransaction("updateTokenVolume", tokenId, currentAmount - amount, "Token Retired");
            result = JSON.parse(result.toString());
            message = `Successfully retired ${amount} tokens!`

        } else {
            return `Invocation requires either generateToken/transferToken/createToken/changeTokenOwner/updateTokenVolume/retireToken/retirePartialToken as function but got ${fcn}`
        }

        await gateway.disconnect();

        let response = {
            message: message,
            TxId: result.txID
        }

        return response;

    } catch (error) {
        console.log(`Getting error: ${error}`)
        return error.message
    }
}

exports.invokeTransaction = invokeTransaction;