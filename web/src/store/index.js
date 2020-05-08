import Vue from 'vue'
import Vuex from 'vuex'
// import {
//     Crypto, RpcClient, utils, Parameter, ParameterType, TransactionBuilder,
//     WebsocketClient, OntAssetTxBuilder
// } from 'ontology-ts-sdk'

const Crypto = Ont.Crypto
const RpcClient = Ont.RpcClient;
const utils = Ont.utils;
const Parameter = Ont.Parameter, ParameterType = Ont.ParameterType;
const TransactionBuilder = Ont.TransactionBuilder;
const WebsocketClient = Ont.WebsocketClient;
const OntAssetTxBuilder = Ont.OntAssetTxBuilder;

import BigNumber from 'bignumber.js'
Vue.use(Vuex)

const onchain_tokens = [
    {
        token: 'ONT',
        balance: 0
    },
    {
        token: 'ONG',
        balance: 0.00
    }
    
]
const offchain_tokens = [
    {
        token: 'XONT',
        balance: 0
    },
    {
        token: 'XONG',
        balance: 0.00
    }
]

const actions = {
 
}

const privateKey = sessionStorage.getItem('layer2_privateKey') || ''
const address = sessionStorage.getItem('layer2_address') || ''
export default new Vuex.Store({
    state: {
        onchain_tokens: onchain_tokens,
        offchain_tokens: offchain_tokens,
        layer2_rpc: process.env.VUE_APP_LAYER2_RPC, //'http://172.168.3.59:40336',
        layer2_socket: process.env.VUE_APP_LAYER2_SOCKET,  //'ws://172.168.3.59:40335',
        contract_address: process.env.VUE_APP_CONTRACT_HASH,
        address: address,
        privateKey: privateKey
    },
    mutations: {
        UPDATE_OFFCHAIN_BALANCE(state, { ont, ong }) {
            const ongV = new BigNumber(ong).div(1e9).toNumber()
            const offchain_tokens = [
                { token: 'XONT', balance: Number(ont) },
                {token: 'XONG', balance: Number(ongV)}
            ]
            state.offchain_tokens = offchain_tokens
        },
        UPDATE_ONCHAIN_BALANCE(state, { ont, ong }) {
            const ongV = new BigNumber(ong).div(1e9).toNumber()
            const onchain_tokens = [
                { token: 'ONT', balance: Number(ont) },
                {token: 'ONG', balance: Number(ongV)}
            ]
            state.onchain_tokens = onchain_tokens;
        },
        UPDATE_PRIVATEkEY(state, privateKey) {
            sessionStorage.setItem('layer2_privateKey', privateKey)
            state.privateKey = privateKey
        },
        UPDATE_ADDRESS(state, address) {
            sessionStorage.setItem('layer2_address', address)
            state.address = address
        }
    },
    actions: {
        async getBalance({commit, state}, address) {
            const addr = new Crypto.Address(address)
            const rpcClient_layer2 = new RpcClient(state.layer2_rpc)
            const balance_offchain = await rpcClient_layer2.getBalance(addr)
            console.log(balance_offchain)
            if(balance_offchain.error === 0) {
                const { ont, ong } = balance_offchain.result
                commit('UPDATE_OFFCHAIN_BALANCE', {ont, ong})
            }
            const rpcClient = new RpcClient()
            const balance_onchain = await rpcClient.getBalance(addr)
            console.log(balance_onchain)
            if(balance_onchain.error === 0) {
                const { ont, ong } = balance_onchain.result
                commit('UPDATE_ONCHAIN_BALANCE', {ont, ong})
            }
        },
        async deposit({ commit, dispatch, state }, { amount, asset }) {
            const contract = utils.reverseHex(state.contract_address);
            const contractAddr = new Crypto.Address(contract);
            let amountV = amount;
            const method = 'deposit';
                const payer = new Crypto.Address(state.address)
                let assetAddress = ''
                if (asset === 'ONT') {
                    assetAddress = '0000000000000000000000000000000000000001'
                    amountV = Number(amount)
                } else {
                    assetAddress = '0000000000000000000000000000000000000002'
                    amountV = new BigNumber(amount).times(1e9).toNumber()
                }
            const params = [
                new Parameter('player', ParameterType.Address, payer),
                new Parameter('amount', ParameterType.Integer, amountV),
                new Parameter('asset', ParameterType.ByteArray, assetAddress)
            ]

            const tx = TransactionBuilder.makeInvokeTransaction(method, params, contractAddr, '500', '200000', payer);
            const privateKey = new Crypto.PrivateKey(state.privateKey)
            TransactionBuilder.signTransaction(tx, privateKey);
            const socketClient = new WebsocketClient()
            try {
                const res = await socketClient.sendRawTransaction(tx.serialize(), false, true);
                console.log(JSON.stringify(res));
                dispatch('getBalance', state.address)
                return res;
            } catch (err) {
                console.log(err)
                return err
            }
            
        },
        async withdraw({ commit, dispatch, state }, { amount, asset }) {
            debugger
            const payer = new Crypto.Address(state.address);
            const to = new Crypto.Address('0000000000000000000000000000000000000000');
            let assetV, amountV;
            if (asset === 'XONT') {
                assetV = 'ONT'
                amountV = Number(amount)
            } else {
                assetV = 'ONG'
                amountV = new BigNumber(amount).times(1e9).toNumber()
            }
            const tx = OntAssetTxBuilder.makeTransferTx(assetV, payer, to, amountV, '0', '200000', payer);
            tx.isLayer2Node = true
            const privateKey = new Crypto.PrivateKey(state.privateKey)
            TransactionBuilder.signTransaction(tx, privateKey);
            const socketClient = new WebsocketClient(state.layer2_socket)
            try {
                const response = await socketClient.sendRawTransaction(tx.serialize(), false, true);
                // tslint:disable:no-console
                console.log(JSON.stringify(response));
                dispatch('getBalance', state.address)
                return response
            } catch (err) {
                console.log(err)
                return err
            }
            
        },
        async sendToken({ commit, dispatch, state }, { amount, to, token }) {
            const payer = new Crypto.Address(state.address);
            const receiver = new Crypto.Address(to);
            let assetV, amountV;
            if (token === 'XONT' || token === 'ONT') {
                assetV = 'ONT'
                amountV = Number(amount)
            } else if(token === 'XONG' || token === 'ONG') {
                assetV = 'ONG'
                amountV = new BigNumber(amount).times(1e9).toNumber()
            }
            debugger
            if (token === 'XONT' || token === 'XONG') {
                const tx = OntAssetTxBuilder.makeTransferTx(assetV, payer, receiver, amountV, '0', '200000', payer);
                tx.isLayer2Node = true;
                const privateKey = new Crypto.PrivateKey(state.privateKey)
                TransactionBuilder.signTransaction(tx, privateKey);
                const socketClient = new WebsocketClient(state.layer2_socket)
                try {
                    const response = await socketClient.sendRawTransaction(tx.serialize(), false, true);
                    // tslint:disable:no-console
                    console.log(JSON.stringify(response));
                    dispatch('getBalance', state.address)
                    return response
                } catch (err) {
                    return err
                }
                
            } else {
                const tx = OntAssetTxBuilder.makeTransferTx(assetV, payer, receiver, amountV, '500', '200000', payer);
                const privateKey = new Crypto.PrivateKey(state.privateKey)
                TransactionBuilder.signTransaction(tx, privateKey);
                const socketClient = new WebsocketClient()
                try {
                    const response = await socketClient.sendRawTransaction(tx.serialize(), false, true);
                    // tslint:disable:no-console
                    console.log(JSON.stringify(response));
                    dispatch('getBalance', state.address)
                    return response
                } catch (err) {
                    return err
                }
                
            }
            
        }
    },
    modules: {
    }
})
