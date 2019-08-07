# Algorand

This package implements handlers which can be used to interact with the Algorand blockchain. [The Algorand documentation](https://developer.algorand.org/docs/introduction-installing-node#start-node) is a good place to start and has instructions that work seamlessly. Once your node is up, you need to wait for it to finish syncing up to the current height and then you're free to experiment with the [Go SDK](https://github.com/algorand/go-algorand-sdk/)

To start a private network,
`cd private/node` to fetch `algod.net` and `algod.token` and link it with the SDK.

For getting a kmd access token,
`cd private` and `cat kmd-v0.5/kmd.token`.
