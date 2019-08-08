package consts

import (
	"os"

	email "github.com/Varunram/essentials/email"
	ipfs "github.com/Varunram/essentials/ipfs"
	algorand "github.com/YaleOpenLab/openx/chains/algorand"
	stablecoin "github.com/YaleOpenLab/openx/chains/stablecoin"
	xlm "github.com/YaleOpenLab/openx/chains/xlm"
)

// contains constants - some arbitrary, some forced due to stellar. Each account should have a minimum of 0.5 XLM and each trust line costs 0.5 XLM.
// For more info about Stellar costs see: https://www.stellar.org/developers/guides/concepts/accounts.html

// Platform consts
var PlatformPublicKey string // set this to empty and store during runtime
var PlatformEmail string     // email so we can send notifications to the platform when needed
var PlatformEmailPass string // the password associated with the platform's email id
var PlatformSeed string      // set this to empty and store during runtime
var KYCAPIKey string         // API key to call the KYC provider's API
var Mainnet bool             // bool to denote whether this is mainnet or testnet
var IpfsFileLength int       // length of the ipfs file hash
var RefillAmount float64     // refill amount (testnet only)

// stablecoin related consts
var StablecoinCode string
var StablecoinPublicKey string
var StablecoinSeed string
var StableCoinSeedFile string
var StablecoinTrustLimit float64
var AnchorUSDCode string
var AnchorUSDAddress string
var AnchorUSDTrustLimit float64
var AnchorAPI string

// algorand consts
var AlgodAddress string
var AlgodToken string
var KmdAddress string
var KmdToken string

func SetConsts(mainnet bool) {
	if !mainnet {

		HomeDir += "/testnet"
		DbDir = HomeDir + "/database/"
		PlatformSeedFile = HomeDir + "/platformseed.hex"

		StablecoinCode = "STABLEUSD"                                                     // this is constant across different pubkeys
		StablecoinPublicKey = "GDUNL6KJKSCBEFTVRXDN7ZPLCYXUDZUTJRZ4IUHD7DQ5YHGO6EPPYBTZ" // set this after running this the first time. replace for tests to run properly
		StablecoinSeed = "SDGD3L3KVINB5OICLAQRQLGTNDDWRBELTQZG6FHFYCEZRQPWKEOQOWR6"      // set this after running this the first time. replace for tests to run properly
		StableCoinSeedFile = DbDir + "/stablecoinseed.hex"
		StablecoinTrustLimit = 1000000000
		// testnet anchor params
		AnchorUSDCode = "USD"
		AnchorUSDAddress = "GCKFBEIYV2U22IO2BJ4KVJOIP7XPWQGQFKKWXR6DOSJBV7STMAQSMTGG"
		AnchorUSDTrustLimit = 1000000
		AnchorAPI = "https://sandbox-api.anchorusd.com/"

		// algorand stuff is only enabled with stellar testnet and not mainnet
		AlgodAddress = "http://localhost:50435"
		AlgodToken = "df6740f7618f699b0417f764b6447fa7e690f9514c73cd60184314ae16141030"
		KmdAddress = "http://localhost:51976"
		KmdToken = "755071c9616f4ebac31512e4db7993dc056f12790d94d634e978a66dfc44ce9b"
		algorand.SetConsts(AlgodAddress, AlgodToken, KmdAddress, KmdToken)

		RefillAmount = 10
		xlm.SetConsts(RefillAmount, Mainnet)

		stablecoin.SetConsts(StablecoinCode, StablecoinPublicKey, StablecoinSeed, StableCoinSeedFile, StablecoinTrustLimit,
			AnchorUSDCode, AnchorUSDAddress, AnchorUSDTrustLimit, Mainnet)

		email.SetConsts(PlatformEmail, PlatformEmailPass)

		IpfsFileLength = 10
		ipfs.SetConsts(IpfsFileLength)
	} else {
		// init mainnet params
		HomeDir += "/mainnet"
		DbDir = HomeDir + "/database/"
		PlatformSeedFile = HomeDir + "/platformseed.hex"

		// set in house stablecoin params to zero to not trade in it
		StablecoinPublicKey = ""
		StableCoinSeedFile = ""
		StablecoinSeed = ""
		StablecoinTrustLimit = 0

		// set anchor mainnet params to exchange
		AnchorUSDCode = "USD"
		AnchorUSDAddress = "GDUKMGUGDZQK6YHYA5Z6AY2G4XDSZPSZ3SW5UN3ARVMO6QSRDWP5YLEX"
		AnchorUSDTrustLimit = 10000 // conservative limit of USD 10000 set for investments on mainnet. Can be increased or decreased as necessary
		AnchorAPI = "https://api.anchorusd.com/"

		stablecoin.SetConsts(StablecoinCode, StablecoinPublicKey, StablecoinSeed, StableCoinSeedFile, StablecoinTrustLimit,
			AnchorUSDCode, AnchorUSDAddress, AnchorUSDTrustLimit, Mainnet)

		RefillAmount = 0
		xlm.SetConsts(RefillAmount, Mainnet)

		email.SetConsts(PlatformEmail, PlatformEmailPass)

		IpfsFileLength = 10
		ipfs.SetConsts(IpfsFileLength)
	}
}

// directories
var HomeDir = os.Getenv("HOME") + "/.openx"          // home directory where we store everything
var DbDir = HomeDir + "/database/"                   // the directory where the database is stored (project info, user info, etc)
var DbName = "openx.db"                              // the name of the db that we want to store stuff in
var PlatformSeedFile = HomeDir + "/platformseed.hex" // where the platform's seed is stored

var Tlsport = 443         // default port for ssl
var DefaultRpcPort = 8080 // the default port on which the rpc server of the platform starts. Defaults to HTTPS

// var OpenSolarIssuerDir = HomeDir + "/projects/"      // the directory where we store opensolar projects' issuer seeds
// var OpzonesIssuerDir = HomeDir + "/opzones/"         // the directory where we store ozpones projects' issuer seeds

// ports + number consts
