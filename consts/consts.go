package consts

import (
	"os"
	"time"
)

// not possible to write tests for this package
// This script relates to all constant values used in the whole code. Some relate to the specific simulation instance, which will be updated later. Some values are arbitrary and others relate to the Stellar model.
// DonateBalance is the minimum amount of lumens (XLM) an investor account should have to open itself and invest on a project. We charge this amount when they do their first investment.
// For more info about Stellar costs see: https://www.stellar.org/developers/guides/concepts/accounts.html

// we send this amount free to invesotrs who signup on our platform to enable them to have trustlines. Maybe we should have a payment provider and take money from them?
// DonateBalance
var DonateBalance = "10"

// the maximum limit that the investor trusts the stablecoin issuer for
// StablecoinTrustLimit
var StablecoinTrustLimit = "1000000000"

// the prefix that will be hashed to give an investor AssetID
// InvestorAssetPrefix
var InvestorAssetPrefix = "InvestorAssets_"

// the prefix that will be hashed to give a bond asset
// BondAssetPrefix
var BondAssetPrefix = "BondAssets_"

// the prefix that will be hashed to give the cooperative asset
// CoopAssetPrefix
var CoopAssetPrefix = "CoopAsset_"

// the prefix that will be hashed to give a recipient AssetID
// DebtAssetPrefix
var DebtAssetPrefix = "DebtAssets_"

// the prefix that will be hashed to give an ivnestor his seed id
// SeedAssetPrefix
var SeedAssetPrefix = "SeedAssets_"

// the prefix that will be hashed to give a payback AssetID
// PaybackAssetPrefix
var PaybackAssetPrefix = "PaybackAssets_"

// home directory where we store the platform seed
// HomeDir
var HomeDir = os.Getenv("HOME") + "/.openx"

// the path where the platform's seed is stored
// PlatformSeedFile
var PlatformSeedFile = HomeDir + "/platformseed.hex"

// the path where the stablecoin's seed is stored
// StableCoinSeedFile
var StableCoinSeedFile = HomeDir + "/stablecoinseed.hex"

// the publickey of the address issuing the asset STABLEUSD
// StablecoinPublicKey
var StablecoinPublicKey = ""

// the seed of the address issuing the STABLEUSD asset
// StablecoinSeed
var StablecoinSeed = ""

// the address of the stablecoin must be a constant for the payment listener to work properly
// StableCoinAddress
const StableCoinAddress = "GCX4BSWDWHDDGLLTA6C73NENJCQBST7C4B4W5HZE7ZCSOVWML7VLLLT3"

// the code of the stablecoin we issue on the platform
// Code
var Code = "STABLEUSD"

// the directory where the main assets of our platform are stored
// DbDir
var DbDir = os.Getenv("HOME") + "/.openx/database"

// the directory where we store issuer seeds related to the opensolar platforms
// OpenSolarIssuerDir
var OpenSolarIssuerDir = HomeDir + "/projects/"

// the directory where we store issuer seeds related to the opzones platforms
// OpzonesIssuerDir
var OpzonesIssuerDir = HomeDir + "/opzones/"

// the length of the hash that we want our ipfs hashes to have
// IpfsFileLength
var IpfsFileLength = 10

// the home directory of the teller executable
// TellerHomeDir
var TellerHomeDir = os.Getenv("HOME") + "/.openx/teller"

// default port for ssl
// Tlsport
var Tlsport = 443

// the issuer seed password for unlocking the encrypted file. This is a constant for now, can be varied later if required
// IssuerSeedPwd
var IssuerSeedPwd = "blah"

// second * minute * hour * day * number
// PaybackInterval
var PaybackInterval = time.Duration(1 * 60 * 60 * 24 * 30)

// 5 seconds
// TestPaybackInterval
var TestPaybackInterval = time.Duration(5)

// the default port on which the rpc server of the platform starts
// DefaultRpcPort
var DefaultRpcPort = 8080

// every 5 minutes we refresh the teller to import the changes on the platform
// LoginRefreshInterval
var LoginRefreshInterval = time.Duration(5 * 60)

// set this to empty and store during runtime
// PlatformPublicKey
var PlatformPublicKey = ""

// define a platform email so that we can send notifications to the platform when needed
// PlatformEmail
var PlatformEmail = ""

// set this to empty and store during runtime
// PlatformSeed
var PlatformSeed = ""

// time a recipient is given to unlock the project and redeem investment
// LockInterval
const LockInterval = int64(1 * 60 * 60 * 24 * 3)

// one week in seconds
// OneWeekInSecond
const OneWeekInSecond = 604800

// one week in seconds, easier to have it here than call it in multiple places
// TwoWeeksInSecond
const TwoWeeksInSecond = time.Duration(1209600)

// six months in seconds, send notification
// SixWeeksInSecond
const SixWeeksInSecond = time.Duration(3628800)

// period when we direct power to the grid
// CutDownPeriod
const CutDownPeriod = time.Duration(4838400)

// this is in seconds, MWTODO: get feedback on polling interval
// TellerPollInterval
const TellerPollInterval = time.Duration(300)

// in bvy bytes, tweak this later to something like 10M after testing
// TellerMaxLocalStorageSize
const TellerMaxLocalStorageSize = 2000
