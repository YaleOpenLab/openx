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
var StablecoinTrustLimit = "100000"                                                  // the maximum limit that the investor trusts the stablecoin issuer for
var InvestorAssetPrefix = "InvestorAssets_"                                          // the prefix that will be hashed to give an investor AssetID
var BondAssetPrefix = "BondAssets_"                                                  // the prefix that will be hashed to give a bond asset
var CoopAssetPrefix = "CoopAsset_"                                                   // the prefix that will be hashed to give the cooperative asset
var DebtAssetPrefix = "DebtAssets_"                                                  // the prefix that will be hashed to give a recipient AssetID
var SeedAssetPrefix = "SeedAssets_"                                                  // the prefix that will be hashed to give an ivnestor his seed id
var PaybackAssetPrefix = "PaybackAssets_"                                            // the prefix that will be hashed to give a payback AssetID
var HomeDir = os.Getenv("HOME") + "/.openfinancing"                                  // home directory where we store the platform seed
var PlatformSeedFile = HomeDir + "/platformseed.hex"                                 // the path where the platform's seed is stored
var StableCoinSeedFile = HomeDir + "/stablecoinseed.hex"                             // the path where the stablecoin's seed is stored
var DbDir = os.Getenv("HOME") + "/.openfinancing/database"                           // the directory where the main assets of our platform are stored
var OpenSolarIssuerDir = HomeDir + "/projects/"                                      // the diriectory where we store the seeds related to the opensolar project
var IpfsFileLength = 10                                                              // the length of the hash that we want our ipfs hashes to have
const StableCoinAddress = "GBGM6IJH6Z54NIIJE6K7KGWLSFKBFNWWAQBD4FMTN27NC7TNYKCF5CFY" // the address of the stabellcoin must be a constant for the payment listener to work properly
var TellerHomeDir = os.Getenv("HOME") + "/.openfinancing/teller"                     // the home directory of the teller executable
var Tlsport = "443"                                                                  // default port for ssl
var IssuerSeedPwd = "blah"                                                           // the issuer seed password for unlocking the encrypted file. This is a constant for now, can be varied later if required
var PaybackInterval = time.Duration(1 * 60 * 60 * 24 * 30)                           // second * minute * hour * day * number
var TestPaybackInterval = time.Duration(5)                                           // 5 seconds
var DefaultRpcPort = 8080                                                            // the default port on which the rpc server of the platform starts
var LoginRefreshInterval = time.Duration(5 * 60)                                     // every 5 minutes we refresh the teller to import the changes on the platform
