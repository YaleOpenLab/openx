package opensolar

import (
	platform "github.com/YaleOpenLab/openx/platforms"
)

func InitializePlatform() (string, string, error) {
	return platform.InitializePlatform()
}

// RefillPlatform checks whether the publicKey passed has any xlm and if its balance
// is less than 21 XLM, it proceeds to ask the friendbot for more test xlm
func RefillPlatform(publicKey string) error {
	return platform.RefillPlatform(publicKey)
}
