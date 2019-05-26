package consts

import (
	"os"
	"time"
)

// not possible to write tests for this package
// This script relates to all constant values used in the whole code. Some relate to the specific simulation instance, which will be updated later. Some values are arbitrary and others relate to the Stellar model.
// DonateBalance is the minimum amount of lumens (XLM) an investor account should have to open itself and invest on a project. We charge this amount when they do their first investment.
// For more info about Stellar costs see: https://www.stellar.org/developers/guides/concepts/accounts.html
var DonateBalance = "10"                                                             // we send this amount free to invesotrs who signup on our platform to enable them to have trustlines. Maybe we should have a payment provider and take money from them?
var StablecoinTrustLimit = "1000000000"                                              // the maximum limit that the investor trusts the stablecoin issuer for
var InvestorAssetPrefix = "InvestorAssets_"                                          // the prefix that will be hashed to give an investor AssetID
var BondAssetPrefix = "BondAssets_"                                                  // the prefix that will be hashed to give a bond asset
var CoopAssetPrefix = "CoopAsset_"                                                   // the prefix that will be hashed to give the cooperative asset
var DebtAssetPrefix = "DebtAssets_"                                                  // the prefix that will be hashed to give a recipient AssetID
var SeedAssetPrefix = "SeedAssets_"                                                  // the prefix that will be hashed to give an ivnestor his seed id
var PaybackAssetPrefix = "PaybackAssets_"                                            // the prefix that will be hashed to give a payback AssetID
var HomeDir = os.Getenv("HOME") + "/.openx"                                          // home directory where we store the platform seed
var PlatformSeedFile = HomeDir + "/platformseed.hex"                                 // the path where the platform's seed is stored
var StableCoinSeedFile = HomeDir + "/stablecoinseed.hex"                             // the path where the stablecoin's seed is stored
var StablecoinPublicKey = ""                                                         // the publickey of the address issuing the asset STABLEUSD
var StablecoinSeed = ""                                                              // the seed of the address issuing the STABLEUSD asset
const StableCoinAddress = "GDJE64WOXDXLEK7RDURVYEJ5Y5XFHS6OQZCS3SHO4EEMTABEIJXF6SZ5" // the address of the stablecoin must be a constant for the payment listener to work properly
var Code = "STABLEUSD"                                                               // the code of the stablecoin we issue on the platform
var DbDir = os.Getenv("HOME") + "/.openx/database"                                   // the directory where the main assets of our platform are stored
var OpenSolarIssuerDir = HomeDir + "/projects/"                                      // the directory where we store issuer seeds related to the opensolar platforms
var OpzonesIssuerDir = HomeDir + "/opzones/"                                         // the directory where we store issuer seeds related to the opzones platform
var IpfsFileLength = 10                                                              // the length of the hash that we want our ipfs hashes to have
var TellerHomeDir = os.Getenv("HOME") + "/.openx/teller"                             // the home directory of the teller executable
var Tlsport = 443                                                                    // default port for ssl
var IssuerSeedPwd = "blah"                                                           // the issuer seed password for unlocking the encrypted file. This is a constant for now, can be varied later if required
var EscrowPwd = "blah"                                                               // the escrowpwd is the password used for locking the seed used by the escrow
var PaybackInterval = time.Duration(1 * 60 * 60 * 24 * 30)                           // second * minute * hour * day * number
var TestPaybackInterval = time.Duration(5)                                           // 5 seconds
var DefaultRpcPort = 8080                                                            // the default port on which the rpc server of the platform starts. Defaults to HTTPS
var LoginRefreshInterval = time.Duration(5 * 60)                                     // every 5 minutes we refresh the teller to import the changes on the platform
var PlatformPublicKey = ""                                                           // set this to empty and store during runtime
var PlatformEmail = ""                                                               // define a platform email so that we can send notifications to the platform when needed
var PlatformSeed = ""                                                                // set this to empty and store during runtime
const LockInterval = int64(1 * 60 * 60 * 24 * 3)                                     // time a recipient is given to unlock the project and redeem investment
const OneWeekInSecond = 604800                                                       // one week in seconds
const TwoWeeksInSecond = 1209600                                                     // one week in seconds, easier to have it here than call it in multiple places
const SixWeeksInSecond = 3628800                                                     // six months in seconds, send notification
const CutDownPeriod = 4838400                                                        // period when we direct power to the grid
const TellerPollInterval = 30000                                                     // this is in seconds, MW: What unit is the 300 interval on? Seconds?
const TellerMaxLocalStorageSize = 2000                                               // in bvy bytes, tweak this later to something like 10M after testing
var KYCAPIKey = ""                                                                   // the API key to call the KYC provider's API
var AnchorUSDCode = ""                                                               // The code for the AnchorUSD stablecoin
var AnchorUSDAddress = ""                                                            // The Address issuing AnchorUSD
var AnchorUSDTrustLimit = "1000000"                                                  // The Trust limit towards which a person can trust AnchorUSD
