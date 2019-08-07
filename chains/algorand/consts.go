package algorand

var (
	// AlgodAddress is the address of the algod daemon
	AlgodAddress string
	// AlgodToken is the auth token needed to call algod's rpc functions
	AlgodToken string
	// KmdAddress is the address of the key management daemon
	KmdAddress string
	// KmdToken is the auth token needed to call the kmd's rpc functions
	KmdToken string
)

func SetConsts(address string, token string, kmdaddress string, kmdtoken string) {
	AlgodAddress = address
	AlgodToken = token
	KmdAddress = kmdaddress
	KmdToken = kmdtoken
}
