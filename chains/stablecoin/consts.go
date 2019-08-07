package stablecoin

var (
	// StablecoinCode is the code of the test stablecoin
	StablecoinCode string
	// StablecoinPublicKey is the publickey of the stablecoin
	StablecoinPublicKey string
	// StablecoinSeed is the seed of the stablecoin
	StablecoinSeed string
	// StableCoinSeedFile denotes the file location of the stablecoin's seed file
	StableCoinSeedFile string
	// StablecoinTrustLimit denotes the trust limit of the stablecoin
	StablecoinTrustLimit float64
	// AnchorUSDCode is the code of Anchor's stablecoin
	AnchorUSDCode string
	// AnchorUSDAddress denotes the address of Anchor's stablecoin
	AnchorUSDAddress string
	// AnchorUSDTrustLimit is the default trust limit for trusting Anchor
	AnchorUSDTrustLimit float64
	// Mainnet is a boolena value that should be set to switch to the mainnet Anchor API
	Mainnet bool
)

func SetConsts(code string, pubkey string, seed string, seedfile string, trustLimit float64,
	anchorUSDCode string, anchorUSDAddress string, anchorUSDTrustLimit float64, mainnet bool) {

	StablecoinCode = code
	StablecoinPublicKey = pubkey
	StablecoinSeed = seed
	StableCoinSeedFile = seedfile
	StablecoinTrustLimit = trustLimit
	AnchorUSDCode = anchorUSDCode
	AnchorUSDAddress = anchorUSDAddress
	AnchorUSDTrustLimit = anchorUSDTrustLimit
	Mainnet = mainnet
}
