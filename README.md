# Ontology Layer2

English | [中文](README_CN.md)

## Key Terminology

### Layer2 Transactions

A transaction or an execution request carried out on Layer2. The user authorizes the action by signing it. A transaction transmitted on the Layer2 may or may not have the same format as a standard transaction on the Ontology main chain.

### Node

A Layer2 node collects all the user transactions, verifies them, and then executes them. The node is mainly responsible for executing the transactions that are part of a newly generated block, updates the state, and generates a state report that is readable by a Layer2 smart contract in order to set a security measure in place.

### Layer2 Block

A node periodically generates Layer2 blocks by collecting the transactions for a particular cycle and encapsulating them in a block.

### Layer2 State

When a node generates a block containing the transactions for a cycle and updates the state, all the relevant updated state data is sorted to generate a merkle tree. This merkle tree's root hash is then calculated, which is the Layer2 state for the corresponding block.

### Operator

The operator is a security daemon on Layer2 that monitors the transactions on Layer2 and verifies whether the sent tokens are transferred successfully to the Ontology main chain. The operator also periodically sends the Layer2 state proof to the Ontology main chain as evidence of transactions having taken place on Layer2.

### Challenger

The challenger verifies the state proof submitted by the operator to the Ontology main chain. This would require the challenger to synchronize all the transactions that have taken place on Layer2 and maintain a complete state record. Upon synchronizing the executed transactions and updating the state, the challenger can confirm the validity of the state proof submitted by the operator to the main net. If found invalid, the challenger can generate a proof of fraud that can be read by the Layer2 contract to dispute the operator.


### Account State Proof

An account state proof includes the account state along with it's merkle proof. The state proof can be fetched from the operator and the challenger, since they both maintain complete state records.

### Proof of Fraud

The proof of fraud contains the account state proof before the current Layer2 block updated it. Since both state proofs are available, the legitimacy of the state proof before the current block updation can be verified. The current state proof can be verified by running the current block.

## Operation Process Flow

### Deposit on Layer2

1.	The user carries out a deposit action on the Ontology main chain. The contract locks the user's deposit amount and records it in the Layer2 state. The current state of this deposit amount is "unreleased".

2.	The operator detects the deposit action on the main chain and submits the deposit transaction details to the Layer2 node.

3.	The node creates a new Layer2 block that encapsulates the deposit transaction along with transactions from different users. The Layer2 state is then updated when this block is executed.

4.	An operator monitors the Layer2 nodes and the new blocks generated. When the updated Layer2 state is sent to the Ontology main chain, along with the request to release the deposit amount on the main chain.

5.	The main chain releases the deposit amount and updates the deposit amount state to released.

<div align=center><img width="360" height="450" src="doc/pic/user_deposit.png"/></div>

### Withdrawal to Ontology

1.	The user creates a Layer2 withdrawal transaction and sends it to the node.

2.	The node updates the state based on the withdrawal amount and generates a Layer2 block including the rest of the transactions.

3.	The operator sends the Layer2 block state to the Ontology main chain along with the withdrawal request.

4.	The main chain contract carries out the withdrawal transaction request and records the withdrawal action, and then sets the state to not released.

5.	Once the state is confirmed, the user sends a withdrawal release request.

6.	The main chain processes and executes the withdrawal release request, sends the withdrawal amount to the target account, and sets the withdraw state to released.

<div align=center><img width="499" height="450" src="doc/pic/user_withdraw.png"/></div>

###	User Layer Transactions

1.	The user creates a Layer2 transfer transaction and sends it to the node.

2.	The node encapsulates the transactions including this transfer transaction to create a Layer2 block. The transactions of this block are executed by Layer2, and the updated block state is then sent to the Ontology main chain.

3.	System then awaits state confirmation.

### Transactions and Security Measures

<div align=center><img width="650" height="450" src="doc/pic/system.png"/></div>

## Account Infrastructure

An account is created using a trackable Merkle tree data structure. A merkle tree contains and maintains the original root hash and the updated account record.

<div align=center><img width="450" height="660" src="doc/pic/account.png"/></div>

Every state root corresponds to a fixed height starting from 0 at the bottom and progressively increasing towards the top. Each state root corresponding to each height is recorded on the chain.

### How to Prove the Validity of an Account State?

Say an account state gets updated at some arbitrary height. At this height, the state tree contains the current account state and exists as a Merkle tree in itself, a sub-tree with respect to the complete merkle tree. Now, the root hash of this sub-tree exists on the chain and can be used to generate the merkle proof. Thus, this root hash can be used to validate this account's current state.

But there is a possibility that this account state isn't the latest, since when the account state is updated so is the merkle tree subsequently. The challenge mechanism can be used to determine the account validity by submitting the account merkle proof. If the calculated root hash of this state is higher than that of the current block state, the challenge is considered successful.

### Why is such a Merkle Tree Necessary?

The state root of all the heights of the merkle tree are recorded on the chain. This information can be used to verify the validity of account state at any given point of time as it gets updated. The change that takes place in account state can be represented to be a function in two variables in the following manner:

**`S = F(S'，Txs)`**

Here, `S'` is any arbitrary state, `Txs` is a transaction, and `S` is the new state that is obtained by updating `S'` upon execution of the transaction `Txs`. 

The simplest on-chain method of establishing the validity of a particular state is by sending the complete `S'` and `Txs` record to the chain. The validity of the updated state is established by comparing its state root to that of the on-chain state `S'`, which is calculated by executing the `Txs` transactions.

However, there are certain issues with this method, as stated below:

1.	Calculating the state by executing all the transaction records is a process that usually has limitations, especially considering the fact that the total state usually consists of a very large number of transactions.

2. Under normal circumstances there is a very small no. of transactions that actually have an effect on the account state.
   

The small no. of accounts can be used to track the state changes, and this can be a new implementation method.

Instead of sending the complete `S'` it is possible to send the `Txs` transactions that update the state `S'` and their respective merkle proofs. We have already discussed how to validate this particular state. It is possible to calculate the new state root for `S` by partly executing the transactions, generating a new state, and then adding the `S'` state root to it. This can then be used to determine the validity the state change.

The advantages to using the above stated method are:

1.	No matter how many total states there are, the number of sub-trees and their merkle proofs remains small if the number of updated states isn't too large.

2.	Lower state verification overhead with high efficiency since only a small number of transactions and merkle proofs needed to be provided in order to determine state validity.

3.	It is possible to determine the updation process of a state.
