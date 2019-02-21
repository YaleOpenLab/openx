# Energy Web

This document serves as a review of the features and capabilities offered by the energyweb platform. This document will also list the resources available on energyweb.

Energy web is a blockchain platform that is specifically designed for the energy sector. They also run something called the energy web foundation which acts as an entity which has a say / opinion on the architecture of these systems. Currently, their software is still in beta and their test network [Tobaloba](https://tobalaba.etherscan.com) is live now. THis test network does not have many participants except the dev team as shown by the [stats page](http://netstats.energyweb.org) but nodes do seem to be run by other individuals with identifiers representing their respective companies. The chain seems to be forked off from ethereum and seems to follow a proof of authority model (the same consensus model used by geth's testnet) and there seem to be no plans to move off PoA for a while. This is a good decision since PoW would be incompatible with what the company aims to do primarily because there's nothing to prove and there's only stuff to verify, which can be run by select people while still maintaining consensus.

Their [blockchain](https://energyweb.org/blockchain) page goes into detail on what has been done and what they plan to do as part of their roadmap. It is to be noted that their roadmap while not defined on the website is defined in a page in the whitepaper. As of now, the chain seems to have a testnet running, have a client which can be used to interface with their blockchain and has provisions for secret transactions (although this may be a misnomer since these still require something called private validators). In the fuiture, they plan to expand the chain to work with Polkadot (whose aim is to act like a blockchain of blockchains). There also seem to be plans to add payment channels, although that would not make much sense for nun fungible assets such as RECs, which the chain aims to foster. They cite various companies to be using their technology but links to such declaration are absent so it is difficult to verify these claims. The services that would be of use to the openx platform seem to be:

1. EW Origin - reference application for RECs and carbon accounting markets (similar to Swytch. It is also to be noted that Swytch is an "affiliate" of the EWF)
2. EW Link - set of architectures and standards for connecting devices to the blockchain (similar to Atonomi)


EWLink seems to have a few interesting points that may prove useful for this document:

1. Governance by gas - as more people start using a dapp on the EWF blockchain, more gas is consumed. The governance model of the chain is in such a way that application developers of these well used apps have more say in the governance process. All stakeholders can propose updates and the protocol implementation team (PIT) builds the propose protocol. Who constitutes this team is not entirely clear but for the purpose of this document, it can be assumed that the EWF would employ these resources.

2. No Conflict of interest - developers and organizations still have to go through the legal framework on their side to ensure compliance on those parts. Only developers and organizations that are approved and identified either by an existing legal system or by the EWF are allowed to propose dapps. While this may seem like a centralizing influence, it makes sense to do this because the governance mechanism proposed hinges directly on the developers and having byzantine developers blocking specific protocol updates is not a good solution. It is to b e stated here that the conditions for becoming a validator are specified in the whitepaper.

3. Mix of on chain and off chain governance - both peopler and code have say in these systems, which makes any predefined rules malleable.

The token allocation model is defined in the whitepaper with 28.5 percent of the tokens being pre-allocated to the team's founders, initial backers and the dev team itself.

Overall, the project seems exciting and useful to what we aim to do with the openx platform. Their relation with swytch is to be explored and we need to see what kinds of working relationships they have with the companies mentioned as affiliates along with the impact of this EWF (which we might need to consult in the future if we want to do this stuff).

Resources:
1. http://netstats.energyweb.org
2. https://tobalaba.etherscan.com
3. https://energyweb.org/wp-content/uploads/2018/10/EWF-Paper-TheEnergyWebChain-v1-201810-FINAL.pdf
