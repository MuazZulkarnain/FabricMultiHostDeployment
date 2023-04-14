const {
  Gateway,
  Wallets,
  TxEventHandler,
  GatewayOptions,
  DefaultEventHandlerStrategies,
  TxEventHandlerFactory,
} = require("fabric-network");
const fs = require("fs");
const path = require("path");
const log4js = require("log4js");
const logger = log4js.getLogger("BasicNetwork");
const util = require("util");

// const createTransactionEventHandler = require('./MyTransactionEventHandler.ts')

const helper = require("./helper");

// const createTransactionEventHandler = (transactionId, network) => {
//     /* Your implementation here */
//     const mspId = network.getGateway().getIdentity().mspId;
//     const myOrgPeers = network.getChannel().getEndorsers(mspId);
//     return new MyTransactionEventHandler(transactionId, network, myOrgPeers);
// }

const invokeTransaction = async (
  channelName,
  chaincodeName,
  fcn,
  args,
  username,
  org_name,
  transientData
) => {
  try {
    logger.debug(
      util.format(
        "\n============ invoke transaction on channel %s ============\n",
        channelName
      )
    );

    // load the network configuration
    // const ccpPath =path.resolve(__dirname, '..', 'config', 'connection-org1.json');
    // const ccpJSON = fs.readFileSync(ccpPath, 'utf8')
    const ccp = await helper.getCCP(org_name); //JSON.parse(ccpJSON);

    // Create a new file system based wallet for managing identities.
    const walletPath = await helper.getWalletPath(org_name); //path.join(process.cwd(), 'wallet');
    const wallet = await Wallets.newFileSystemWallet(walletPath);
    console.log(`Wallet path: ${walletPath}`);

    // Check to see if we've already enrolled the user.
    let identity = await wallet.get(username);
    if (!identity) {
      console.log(
        `An identity for the user ${username} does not exist in the wallet, so registering user`
      );
      await helper.getRegisteredUser(username, org_name, true);
      identity = await wallet.get(username);
      console.log("Run the registerUser.js application before retrying");
      return;
    }

    const connectOptions = {
      wallet,
      identity: username,
      discovery: { enabled: true, asLocalhost: false },
      eventHandlerOptions: {
        commitTimeout: 100,
        strategy: DefaultEventHandlerStrategies.NETWORK_SCOPE_ALLFORTX,
      },
      // transaction: {
      //     strategy: createTransactionEventhandler()
      // }
    };

    // Create a new gateway for connecting to our peer node.
    const gateway = new Gateway();
    await gateway.connect(ccp, connectOptions);

    const network = await gateway.getNetwork(channelName);
    const contract = network.getContract(chaincodeName);

    let message;
    let result;

    if (fcn === "createToken") {
      let existingToken = await contract.evaluateTransaction(
        "queryToken",
        args[0]
      );

      if (existingToken.toString()) {
        result = await contract.submitTransaction(
          "updateTokenVolume",
          args[0],
          args[1]
        );
      } else {
        result = await contract.submitTransaction(
          fcn,
          args[0],
          args[1],
          args[2],
          args[3]
        );
      }
      result = JSON.parse(result.toString());
      message = `Successfully updated ${args[1]} token to the ${args[0]}`;
    } else if (fcn === "transferToken") {
      let tokenKey = args[0];
      let newOwner = args[1];
      let transferAmount = args[2];

      // Deduct transfer amount from the current owner's token volume
      await contract.submitTransaction(
        "updateTokenVolume",
        tokenKey,
        -transferAmount
      );

      // Add transfer amount to the new owner's token volume
      let existingToken = await contract.evaluateTransaction(
        "queryToken",
        tokenKey
      );
      if (existingToken.toString()) {
        let token = JSON.parse(existingToken.toString());
        let newVolume = parseInt(token.volume) + parseInt(transferAmount);
        await contract.submitTransaction(
          "updateTokenVolume",
          tokenKey,
          newVolume.toString()
        );
      } else {
        return `Token with key ${tokenKey} does not exist`;
      }

      result = { txID: "transferTokenSuccess" };
      message = `Successfully transferred ${transferAmount} token from ${identity.mspId} to ${newOwner}`;
    } else if (fcn === "retireToken") {
      result = await contract.submitTransaction(fcn, args[0]);
      result = JSON.parse(result.toString());
      message = `Successfully retired token with key ${args[0]}`;
    } else {
      return `Invocation requires either createToken or transfer Token as function but got ${fcn}`;
    }

    await gateway.disconnect();

    let response = {
      message: message,
      TxId: result.txID,
    };

    return response;
  } catch (error) {
    console.log(`Getting error: ${error}`);
    return error.message;
  }
};

exports.invokeTransaction = invokeTransaction;
