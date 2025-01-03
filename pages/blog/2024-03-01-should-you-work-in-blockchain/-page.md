---
title: Should you work in Blockchain?
description: With the bad reputation of cryptocurrencies, should you work in blockchain?
tags: ['blockchain', 'cryptocurrency', 'ethics']
---

## Table of contents

<div class="toc">

{{% $.TOC %}}

</div>

<hr>

## Introduction

Hey, let's talk a little about Blockchain. Finally reputed for being a technology that hosts scams and illegal activities, I want to help you understand the technology behind it before criticizing it. Then, you can have actual arguments to criticize it.

Before starting on the subject, I want to clarify that I worked in the blockchain industry for the last three years. Looking at that industry from the "tech"-side, you should know that there are aspects of blockchain that is interesting to work with. If you are on the "critic"-side, you are probably wondering why there are still a lot of people working on Blockchain, and why they are not going to work on something more "humanitarian" or "useful".

This is why, I'm going to approach the subject on two different angles:

- The "technical" aspect of blockchain.
- The "financial" aspect of blockchain.
- The "ethical" aspect of blockchain.

Let's start by understanding the problem that blockchain is trying to solve.

## The problem

The problem is quite simple: it's trust. This problem is not new, and was never new to begin with. To resonate with you, let's define trust by taking a very simple example: a Messaging application.

**What is trust in a messaging app?**

Trust is the receipt of a message from someone, and the ability to know that the message was sent by that person. Blockchain tries to prove that receipt permanently. Now, let's ask ourselves this question:

**Why in the legal world, we barely trust anything that is digital?**

The answer is simple: because it's easy to fake. Why? Because the authority of service permits us to do it. For example, on Slack, you can edit your messages. Blockchain to decentralize that authority, or to be more simple, blockchain wants everyone to be witnesses of the message.

There you go, Blockchain wants to **decentralize authority**. A very simple mission, but ethically difficult to achieve.

## The technical aspect of blockchain

### A high-level approach

Let's go back to that example. "Blockchain wants everyone to be witnesses of the message". Isn't actually something very wrong? Where is privacy in that? Well, it's not wrong, but it's not right either. It's a trade-off. If a messaging application was actually developed on Blockchain, it would be something similar to solve the issue of empty promises. Again, it's not about privacy, it's about trust.

Now, technically speaking, the issue of trust exists since the beginning of the internet. Let's take a smaller example, the verrrry beginning of the internet: the email. You are sending an email to `john.doe@example.com`, but who is `example.com`? I mean, ok, `example.com` is the IP `93.184.216.34`... or is it?

You see, due to the "global" ecosystem, it's difficult to assert if one domain (`example.com`) is actually the right email provider. What if your Domain Name Server (DNS) is hacked? What if the IP is not the right one?

One solution is to use a "trusted" third-party, also known as "Certificate Authorities" (CA) (if we take the example of the messaging app, it's the witness). Every OS has a list of trusted CA (`/etc/ssl/certs`, or in Windows Registry `\SOFTWARE\Microsoft\SystemCertificates\`). And that CA certifies that `example.com` is actually `93.184.216.34` by doing a key exchange. The key exchange process involves two sets of keys: the CA key which is used to **certify**, and the domain's key which is used to **encrypt or sign the data**.

TODO: diagram

Remember this well, because blockchain is very similar:

- **An authority certifies that something is true.**
- **A client sends encrypted or signed data.**

Now, small issue, whether it's about high availability or about the CA being too "centralized", we want to share the same "state" for each CA. This is NOT decentralization, this is simply replication. There are a lot of strategies for replication: master writer and read-only replicas, DNS caching, election algorithms, etc. Blockchain wants that everyone share the same state, but without a master writer. So, how do we do that?

Well, we can use a "consensus algorithm". Consensus algorithms are **very** old, and are used in a lot of different fields. For example, in the legal world, a consensus algorithm is a "jury". In the IT world, a consensus algorithm is a "quorum". In the blockchain world, a consensus algorithm is a "proof of work" or a "proof of stake".

But, they all have the same goal: to **make sure that everyone agrees on the same state**.

And blockchain is not the only one to use consensus algorithms. For example, the "Raft" consensus algorithm is used in the "etcd" key-value store.

TODO: diagram

Now, for the final part of the technical aspect, let's talk about the "blockchain" itself. Why do we want to "chain" blocks? What's the point, why not use a simple database?

Let's take either mail or messaging app as an example. To have proof that a message was sent, we need a history of "commands". Commands like "someone has edited a message", "someone has sent a message", "someone has deleted a message". This is the **major** difference between a database and a blockchain. A database with consensus is about protecting "state", **a blockchain with consensus is about protecting "history".**

But again, why protecting history?

Well, if you say that question more loudly, isn't that obvious? Isn't protecting history one of the most important things in the world? It does feel more "humanitarian", doesn't it?

### A low-level approach

Now, lets summary the previous section in a more technical way. A blockchain needs:

- A consensus algorithm to make sure that everyone agrees on the same state.
- Cryptography to make sure that the data is not tampered with (i.e. a way to certify).
- A way to store history.

Again, there are a lot of consensus algorithms. Let's just focus on the [Hashcash](http://www.hashcash.org/papers/announce.txt) algorithm, which is used in the Bitcoin blockchain. The Hashcash algorithm is a "proof of work" algorithm. It's a very simple algorithm: you need to find a number that, when hashed, gives a number with a certain number of leading zeroes. This is called a "nonce". The first one to find the nonce is the one who can add a block to the blockchain. The action of finding the nonce is also called "mining".

TODO: diagram

Now, let's talk about the cryptography. Remember, everyone is their own "CA", therefore everyone has their own set of key. Blockchain is asymetric cryptography, which means that you have a public key and a private key. The public key is used to encrypt data, and the private key is used to decrypt/sign data. The private key is also used to sign data, and the public key is used to verify the signature. Bitcoin uses the [Elliptic Curve Digital Signature Algorithm](https://en.bitcoin.it/wiki/Elliptic_Curve_Digital_Signature_Algorithm) (ECDSA) for its cryptography. And the addresses are generated using the [SHA-256](https://en.bitcoin.it/wiki/Address) algorithm of the public key.

TODO: diagram

Finally, let's talk about the way to store history. A blockchain is a linked list of blocks. Quite simple. Each block contains a hash of the previous block, and a list of transactions. Some transaction can contain a payload (a smart-contract for example). The hash of the previous block is used to make sure that the history is not tampered with. If you change a block, you need to change the hash of the previous block, and so on. This is called a "Merkle tree". The transactions are stored in a "Merkle tree" as well, and the root of the tree is stored in the block. This is used to make sure that the transactions are not tampered with.

TODO: diagram

As you can see, it's technically simple, and could even be interesting to work with. Other technologies use different philosophies to achieve decentralization, like P2P networks, or federation. But, unlike these technologies, blockchain tries to include everyone in the same state, and tries to protect history.

P2P networks are voluntary, and federation is "semi-centralized". Blockchain is purely "decentralized".

## The financial aspect of blockchain

Let's go back to mining (or any consensus algorithm using rewards). Mining has a cost (energetic, hardware, etc.). To incentivize people to mine, the blockchain rewards the miner with a certain amount of cryptocurrency. This is called the "block reward". Because mining has a cost, the cryptocurrency has a value. And because cryptocurrency has a value, there is a market.

Now, do note that blockchain may not need a cryptocurrency. For example, private blockchain with zero-gas fees (i.e. no block reward) can exist. Why? Well, think about scope. Sometime, you do not need the whole world to certify that something is true. Sometime, you just need a group of people. Let's say, a company wants to use blockchain to certify that a document is true. The company can use a private blockchain, and the employees can be the witnesses. To avoid mining, the company can use a "proof of authority" consensus algorithm, where the authority is shared with specific nodes.

But, let's go back to the public blockchain. One usage of such public blockchain combined with cryptocurrency is to avoid the centralization of the financial system. I won't talk about the benefits of such "feature" (deflationary, etc.), because it's highly political. What I can say, is that if we are able to decentralize the financial system, then, we can decentralize the whole web. Which is why people talk about "Web 3.0".

Because the authority is decentralized, what you "own" is truly proven. This is why people are trying to use blockchain to decentralize "everything". From the financial system, to the web, to the storage, to the legal system, to the voting system, etc.

Financially speaking, blockchain has a lot of potential. It's early and conflicts with the current financial system, but it's very novel (and certainly, very "left-wing").

## The ethical aspect of blockchain

There is one major drawback to any decentralized system: everyone is mixed with everyone. Same as Tor, same as BitTorrent, same as ActivityPub. And like any decentralized system, you cannot moderate it. This is why, the blockchain is used for illegal activities, scams, money laundering etc. Your anonymity is protected because everyone is mixed with everyone, that's the whole point.

However, it's not like there is only one blockchain. There are other blockchains that have a different point of view. All of them tries to ask the same questions:

- Privacy vs Transparency: Do you allow money laundering? Or do you prefer to hide your identity?
- Sustainability vs Scalability: Do you want to protect the environment? Or do you want a lot of people to participate?
- Governance: Do you allow a certain authority to change the rules (either by vote, or simply no one)?

Blockchain has developed a lot. We now have smart-contracts, multi-layered blockchains, etc. For each of these developments, there are always bad cases: Blockchain getting hacked, exchanges scamming millions of users, cryptocurrencies with empty promises, and, obviously, NFTs.

Blockchain is also quite young. So young that there are not really a decentralized exchange that can convert Euros to a cryptocurrency. Paradoxically, the only way to convert Euros to a cryptocurrency is to use a centralized exchange. This is why, the "decentralized finance" (DeFi) is a very hot topic. But, it's also a very dangerous topic. There are a lot of scams, and a lot of people lost a lot of money.

Ethically speaking, blockchain is a very difficult subject. It's almost certain that the government will regulate it. But, is it the future? After all, a decentralized financial system could be something to explore.

## What is happening now?

DEX, smart-contracts, NFTs, ... What technology has potential? Which one is a scam? Is it really worth it to work in blockchain?

- DeFi: At least, at small scale, having a transparency about the "credit" system of some website could worth it. Or just the fact of using proven cryptography to make exchanges is already a good advancement. If you haven't tried it, I recommend installing MetaMask and do some transactions.
- NFTs: It's a mostly scams. NFTs or Non-Fungible Tokens tries to give a unique value to a digital asset. The thing is that if the digital asset is not unique, then the NFT is not unique, which mean it's a scam. One way to give uniqueness to a digital asset is not via blockchain, but via its usage. Think like skins in video games. Their uniqueness comes from how that they are used in a game. NFTs can only work by its usage, not because they are NFTs. (People "buying" PNGs are just retarded.) Buying a "certificate of ownership" is not a true NFT, it's a scam.
- Smart-contracts: This is where everything has potential, and it's the core for everything (NFT, DeFi, ..., anything decentralized). Smart-contracts are immutable piece of software that can be audited. It's a very good way to make sure that a program can be trusted, like decentralized storage, or digital signatures, ...

## Conclusion

So, should you work in blockchain? Well, if you want to work in a field with cryptography, distributed system and anything related to decentralization, it's certainly worth it for the technical aspect. Then, for the financial aspect, it's more like a political choice. And finally, for the ethical aspect, the technology is quite young, and like Tor, you will be associated with scammers.

To be fair, blockchain is certainly over-hyped from both side. It's like very similar to ETCD, or Turso. Does some company says: "We are using f-king distributed SQL, so our tech is better than yours, and we are going to be hell-a rich"? No. When working on Blockchain, just think about its benefit, the problem it tries to solve: trust. You can use anything else if the problem is not about trust. Want to roll back code? Use Git. Want to roll back a database? Use backups.

If you're looking for working in Blockchain, their solely reason to say "We are using blockchain" is because they want to solve the problem of trust. If they say "We are using blockchain" because "Web3", because "Unique Value", because "we're gonna be rich" or anything not related to "trust", then it's a scam.

Anyway, as [JB Kempf, the creator of VLC said, "could be technically interesting"](https://youtu.be/H0VgfQMkBIw?si=-q5ywhPnxHCLCg5s&t=2373), but in the same video, "don't believe \[too much] in a technology, because it may not last". Don't make Blockchain a religion, it's just a technology, like SQL.
