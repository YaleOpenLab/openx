package consts

import (
	"os"
	"time"
)

// contains constants - some arbitrary, some forced due to stellar. Each account should have a minimum of 0.5 XLM and each trust line costs 0.5 XLM.
// For more info about Stellar costs see: https://www.stellar.org/developers/guides/concepts/accounts.html

// Platform consts
var RefillAmount = "10"                                     // we send this amount free to invesotrs who signup on our platform to enable them to have trustlines. Maybe we should have a payment provider and take money from them?
var InvestorAssetPrefix = "InvestorAssets_"                 // the prefix that will be hashed to give an investor AssetID
var BondAssetPrefix = "BondAssets_"                         // the prefix that will be hashed to give a bond asset
var CoopAssetPrefix = "CoopAsset_"                          // the prefix that will be hashed to give the cooperative asset
var DebtAssetPrefix = "DebtAssets_"                         // the prefix that will be hashed to give a recipient AssetID
var SeedAssetPrefix = "SeedAssets_"                         // the prefix that will be hashed to give an ivnestor his seed id
var PaybackAssetPrefix = "PaybackAssets_"                   // the prefix that will be hashed to give a payback AssetID
var HomeDir = os.Getenv("HOME") + "/.openx"                 // home directory where we store everything
var DbDir = os.Getenv("HOME") + "/.openx/database"          // the directory where the database is stored (project info, user info, etc)
var OpenSolarIssuerDir = HomeDir + "/projects/"             // the directory where we store opensolar projects' issuer seeds
var OpzonesIssuerDir = HomeDir + "/opzones/"                // the directory where we store ozpones projects' issuer seeds
var Tlsport = 443                                           // default port for ssl
var DefaultRpcPort = 8080                                   // the default port on which the rpc server of the platform starts. Defaults to HTTPS
var IpfsFileLength = 10                                     // the length of the hash that we want our ipfs hashes to have
var IssuerSeedPwd = "blah"                                  // the password for unlocking the encrypted file. This must be modified a compile time and kept secret
var EscrowPwd = "blah"                                      // the password used for locking the seed used by the escrow. This must be modified a compile time and kept secret
var PlatformPublicKey = ""                                  // set this to empty and store during runtime
var PlatformEmail = ""                                      // email so we can send notifications to the platform when needed
var PlatformSeed = ""                                       // set this to empty and store during runtime
var PlatformSeedFile = HomeDir + "/platformseed.hex"        // where the platform's seed is stored
var LockInterval = int64(1 * 60 * 60 * 24 * 3)              // time a recipient is given to unlock the project and redeem investment, right now at 3 days
var PaybackInterval = time.Duration(1 * 60 * 60 * 24 * 30)  // second * minute * hour * day * number, 30 days right now
var OneWeekInSecond = time.Duration(604800 * time.Second)   // one week in seconds
var TwoWeeksInSecond = time.Duration(1209600 * time.Second) // one week in seconds, easier to have it here than call it in multiple places
var SixWeeksInSecond = time.Duration(3628800 * time.Second) // six months in seconds, send notification
var CutDownPeriod = time.Duration(4838400 * time.Second)    // period when we direct power to the grid

// stablecoin related consts
var StablecoinCode = "STABLEUSD"                                                     // code of the stablecoin we issue on the platform
var StablecoinPublicKey = ""                                                         // publickey of the address issuing the asset STABLEUSD
var StablecoinSeed = ""                                                              // seed of the address issuing STABLEUSD
var StableCoinSeedFile = HomeDir + "/stablecoinseed.hex"                             // path where the stablecoin's seed is stored
const StableCoinAddress = "GDJE64WOXDXLEK7RDURVYEJ5Y5XFHS6OQZCS3SHO4EEMTABEIJXF6SZ5" // address of the stablecoin must be a constant for the payment listener daemon to work properly
var StablecoinTrustLimit = "1000000000"                                              // the limit that the investor trusts the stablecoin issuer for / the max number of STABLEUSD that can be granted to a specific user

// teller related consts
var TellerHomeDir = os.Getenv("HOME") + "/.openx/teller"       // the home directory of the teller executable
var TellerMaxLocalStorageSize = 2000                           // in bytes, tweak this later to something like 10M after testing
var TellerPollInterval = time.Duration(30000 * time.Second)    // frequency with which the teller of a particular system is polled
var LoginRefreshInterval = time.Duration(5 * 60 * time.Second) // every 5 minutes we refresh the teller to import the changes on the platform

// third party consts
var KYCAPIKey = ""                                                                // API key to call the KYC provider's API
var AnchorUSDCode = "USD"                                                         // code for the AnchorUSD stablecoin (ref: https://www.anchorusd.com/docs/api#introduction)
var AnchorUSDAddress = "GCKFBEIYV2U22IO2BJ4KVJOIP7XPWQGQFKKWXR6DOSJBV7STMAQSMTGG" // address issuing AnchorUSD (ref: https://www.anchorusd.com/.well-known/stellar.toml)
var AnchorUSDTrustLimit = "1000000"                                               // the limit that the investor trusts AnchorUSD for / the max amount of AnchorUSD that a person can be given.
