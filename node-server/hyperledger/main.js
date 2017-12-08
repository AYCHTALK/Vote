/**
 * Created by stephde on 02.12.17.
 */

let init = require("./initFabricClient.js")
let query = require("./query.js")
let transaction = require("./invoke.js")

const host = 'grpc://localhost:7051';
const channelId = 'mychannel';

function run() {
    let hlAdaper = init.initFabricClient(host, channelId);

    //sample usage of query execution
    // let queryPromise = query.executeQuery(hlAdaper.client, hlAdaper.channel, 'fabcar', 'queryAllCars', [''])
    //sample usage of transaction invocation
    // let transactionPromise = transaction.invokeTransaction(hlAdaper.client, hlAdaper.channel, 'fabcar', '', 'mychannel', [''])

    registerParties(["Green", "Blue", "Yellow", "Red"])
        .then(() => registerEligibleVoters(["voter1", "voter2"]))
        .then(() => configureVoting(new Date(), new Date().setHours(new Date().getHours() + 8)))
        .then(() => enableVote())
}

function registerParties(parties) {
    return Promise.resolve(() => true)
}

function registerEligibleVoters(parties) {
    return Promise.resolve(() => true)
}

function configureVoting(startTime, endTime) {
    return Promise.resolve(() => true)
}

function enableVote() {
    return Promise.resolve(() => true)
}

run();