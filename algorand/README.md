# Algorand

Algorand doesn't have an API like horizon, so we're required to setup a node before we can experiment with what they have. [The Algorand documentation](https://developer.algorand.org/docs/introduction-installing-node#start-node) is a good place to start and has instructions that work seamlessly. Once your node is up, you need to wait for  it to finish syncing up to current height and afterwards, you're free to experiment with the [Go SDK](https://github.com/algorand/go-algorand-sdk/)

Since Algorand doesn't control the testnet unlike Stellar, we have to run a private network to test our applications before porting to mainnet.

Start a private network, `cd` to `private/node` to fetch `algod.net` and `algod.token` to link with the SDK.

For getting the kmd access token, `cd` to `private` and `cat kmd-v0.5/kmd.token`.
