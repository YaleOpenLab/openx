## Anchors on stellar

Anchors are an idea on stellar that one can use to peg a specific asset to another asset or class of assets. For eg, say we want to anchor 1 BTC to 100000 XLM, the process would roughly be:

1. Find an escrow that is willing to exchange BTC for XLM
2. Sign relevant contracts and reneging clauses  with party
3. Use XLM

An anchor roughly does this and the parties that run anchors on Stellar mainnet have real world identities associated with them and are assumed to be liable outside the blockchain world (ie the legal framework) in case they cheat on their customers or show malice. For more info on assets, checkout [the stellar official documentation](https://www.stellar.org/developers/guides/concepts/assets.html)

We only need to identify one anchor for what we aim to do:

1. An anchor that is willing to accept USD and issue a stable coin on stellar for the same, which we can use.

On searching, there seem to be multiple options that deal with this. The below section highlights that with a short description of what the various options do, along with pdf links to their whitepapers in case one wants to delve deeper:

1. StrongholdUSD - seems to be the most popular stablecoin on stellar, backed by IBM's blockchain venture arm. They accept deposits in USD and then give back USD deposits that can be used on stellar to interact with the in-protocol DEX. They claim to be 100% USD backed, have an SEC certified custodian of assets named PrimeTrust based in Nevada and also claim to do KYC/AML for all users wishing to interact with StrongholdUSD. StableUSD is live: https://stellar.expert/explorer/public/asset/USD-GBSTRUSD7IRX73RQZBL3RQUH6KS3O4NYFY3QCALDLZD77XMZOPWAVTUK

2. Tempo - Jed McCaleb is on the board of this stablecoin. They accept deposits in EUR and have a stablecoin named EURT, which they claim is 1:1 backed by on hand Euro. EU regulated company based out of France and the primary European anchor on Stellar. They also claim to do KYC/AML for all their users. Tempo was funded by an ICO and has some pre-mine stuff as well, so we need to audit this part carefully. They also support smaller anchors like coins.ph on their platform.

3. RippleFox - claims to be a CNY anchor on Stellar. They have a one page website with an [address](https://steexp.com/account/GAREELUB43IRHWEASCFBLKHURCGMHE5IF6XSE7EXDLACYHGRHM43RFOX#effects). Account history seems to show some CNY deposits and withdrawals, which shows that they have a working implementation but overall seems not so good considering the total lack of information.

4. AnchorCoin - based out of Canada and accept deposits in CND and have a stablecoin named anchorcoin, but they do not claim they will hold cash deposits for all tokens. Instead, they propose that they will be holding funds in three types of investment - cash, term deposits and mortgages. This is a different model from those detailed above and is not considered standard. They claim that they will provide audits and disclosure regarding all their balances and claim that they are implementing a full KYC and compliance regime.

The above four coins seem to be the ones that explicitly state that they have stablecoins and that they are using them on stellar. [Other anchors](https://www.stellar.org/about/directory#anchors) do not state this explicitly, so we assume that this is not their primary use case / business model. We notice here that there is only one stablecoin per region, so our receivers would have to absolutely trust these entities when it comes to investments in the associated regions.

None of these stable coins are on mainnet, so we need to setup a model bank anchor on mainnet which would simulate their functionality in order to do extensive testing without burning an USD equivalent of coins. This should be easy to do and would require that we implement a new model named "stablecoin" to deal with this and issue them.
