package consts

import (
	"os"
)

var DonateBalance = "10"

// while an investor signs up on our platform, do we send them 10 XLM free?
// do we charge investors to be on our platform? if not, we shouldn't ideally
// be sending them free XLM. also, should the platform have some function for
// withdrawing XLM? if so, we'll become an exchange of sorts and have some
// legal stuff there. If not, we'll just be a custodian and would not have
// too much to consider on our side
var StablecoinTrustLimit = "100000"
var INVAssetPrefix = "INVTokens_"
var DEBAssetPrefix = "DEBTokens_"
var PBAssetPrefix = "PBTokens_"
var HomeDir = os.Getenv("HOME") + "/.opensolar" // home directory where we store seeds
var IpfsFileLength = "10"
