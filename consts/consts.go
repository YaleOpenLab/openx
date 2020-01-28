package consts

import (
	"os"

	algorand "github.com/Varunram/essentials/algorand"
	ipfs "github.com/Varunram/essentials/ipfs"
	xlm "github.com/Varunram/essentials/xlm"
	stablecoin "github.com/Varunram/essentials/xlm/stablecoin"
)

// the consts package contains constants that are specific to openx. These constants
// can be accessed in other platforms by using the platform APIs. If a third party
// platform wants to set consts, it can do so using the SetConsts function below

// PlatformPublicKey is the Stellar public key of the openx platform
var PlatformPublicKey string

// PlatformSeed is the Stellar seed corresponding to the above Stellar public key
var PlatformSeed string

// PlatformEmail is the email of the platform used to send notifications related to openx
var PlatformEmail string

// PlatformEmailPass is the password for the emial account linked above
var PlatformEmailPass string

// KYCAPIKey is the KYC key for ComplyAdvantage, a leading KYC provider which is used with openx
var KYCAPIKey string

// Mainnet denotes if openx is running on Stellar mainnet / testnet
var Mainnet bool

// IpfsFileLength is the length of the temporary ipfs file created during upload of third party documents
var IpfsFileLength int

// RefillAmount is the amount used for refilling an account on Stellar testnet
var RefillAmount float64

// StablecoinCode is the code of the in house stablecoin that openx possesses
var StablecoinCode string

// StablecoinPublicKey is the public Stellar address of our in house stablecoin
var StablecoinPublicKey string

// StablecoinSeed is the Seed corresponding to the above Stablecoin Publickey
var StablecoinSeed string

// StableCoinSeedFile is the location where the encrypted seed is stored and decrypted each time the platform is started
var StableCoinSeedFile string

// StablecoinTrustLimit is the trust limit till which an account trusts openx's stablecoin
var StablecoinTrustLimit float64

// AnchorUSDCode is the code for AnchorUSD's stablecoin
var AnchorUSDCode string

// AnchorUSDAddress is the address associated with AnchorUSD
var AnchorUSDAddress string

// AnchorUSDTrustLimit is the trust limit till which an account trusts AnchorUSD's stablecoin
var AnchorUSDTrustLimit float64

// AnchorAPI is the URL of AnchorUSD's API
var AnchorAPI string

// AlgodAddress is the address of the Algod Daemon
var AlgodAddress string

// AlgodToken is the RPC token required to call Algod endpoints
var AlgodToken string

// KmdAddress is the Algorand Key Management Daemon's address
var KmdAddress string

// KmdToken is the token required to access the Algorand Key Management Daemon
var KmdToken string

// HomeDir is the directory where openx users and other elements specific to openx are stored
var HomeDir = os.Getenv("HOME") + "/.openx"

// DbDir is the directory where the openx database (boltDB) is stored
var DbDir = HomeDir + "/database/"

// DbName is the name of the openx database
var DbName = "openx.db"

// PlatformSeedFile is the location where PlatformSeedFile is stored and decrypted each time the platform is started
var PlatformSeedFile = HomeDir + "/platformseed.hex"

// Tlsport is the default SSL port on which openx starts
var Tlsport = 443

// AccessTokenLife is the life of a generated access token
var AccessTokenLife = int64(3600)

// AccessTokenLength is the length of a user generated access token
var AccessTokenLength = 32

// SetConsts sets the consts required for openx to operate. Third party platforms should
// call this before starting their platform.
func SetConsts(mainnet bool) {
	if !mainnet {
		HomeDir += "/testnet"
		DbDir = HomeDir + "/database/"
		PlatformSeedFile = HomeDir + "/platformseed.hex"

		StablecoinCode = "STABLEUSD"                                                     // this is constant across different pubkeys
		StablecoinPublicKey = "GBESYUIFJ2NKNSLXCDWJJ7YYXD7OTCPWDM57YK6R3U76YEVYS5F5HI37" // set this after running this the first time. replace for tests to run properly
		StablecoinSeed = "SCQET25QJSJU7WU56O2FQBJOXOC37WWUBDDCUMPL53AUII72JOV5JMS2"      // set this after running this the first time. replace for tests to run properly
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

		RefillAmount = 0
	}

	xlm.SetConsts(RefillAmount, Mainnet)
	IpfsFileLength = 10
	ipfs.SetConsts(IpfsFileLength)
	stablecoin.SetConsts(StablecoinCode, StablecoinPublicKey, StablecoinSeed, StableCoinSeedFile, StablecoinTrustLimit,
		AnchorUSDCode, AnchorUSDAddress, AnchorUSDTrustLimit, Mainnet)
}
