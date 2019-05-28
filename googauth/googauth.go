/*
Original package: https://github.com/dgryski/dgoogauth
Copyright (c) 2012 Damian Gryski damian@gryski.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Modifications Copyright (c) 2019-present MIT Digital Currency Initiative under the MIT License

package googauth

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"time"
)

// ComputeCode computes the response code for a 64-bit challenge 'value' using the secret 'secret'.
// To avoid breaking compatibility with the previous API, it returns an invalid code (-1) when an error occurs,
// but does not silently ignore them (it forces a mismatch so the code will be rejected).
func ComputeCode(secret string, value int64) int {

	key, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		return -1
	}

	hash := hmac.New(sha1.New, key)
	err = binary.Write(hash, binary.BigEndian, value)
	if err != nil {
		return -1
	}
	h := hash.Sum(nil)

	offset := h[19] & 0x0f

	truncated := binary.BigEndian.Uint32(h[offset : offset+4])

	truncated &= 0x7fffffff
	code := truncated % 1000000

	return int(code)
}

// OTPConfig is a one-time-password configuration.  This object will be modified by calls to
// Authenticate and should be saved to ensure the codes are in fact only used
// once.
type OTPConfig struct {
	Secret        string // 80-bit base32 encoded string of the user's secret
	WindowSize    int    // valid range: technically 0..100 or so, but beyond 3-5 is probably bad security
	DisallowReuse []int  // timestamps in the current window unavailable for re-use
	UTC           bool   // use UTC for the timestamp instead of local time
}

func (c *OTPConfig) checkTotpCode(t0, code int) bool {

	minT := t0 - (c.WindowSize / 2)
	maxT := t0 + (c.WindowSize / 2)
	for t := minT; t <= maxT; t++ {
		if ComputeCode(c.Secret, int64(t)) == code {

			if c.DisallowReuse != nil {
				for _, timeCode := range c.DisallowReuse {
					if timeCode == t {
						return false
					}
				}

				// code hasn't been used before
				c.DisallowReuse = append(c.DisallowReuse, t)

				// remove all time codes outside of the valid window
				sort.Ints(c.DisallowReuse)
				min := 0
				for c.DisallowReuse[min] < minT {
					min++
				}
				// FIXME: check we don't have an off-by-one error here
				c.DisallowReuse = c.DisallowReuse[min:]
			}

			return true
		}
	}

	return false
}

// Authenticate a one-time-password against the given OTPConfig
// Returns true/false if the authentication was successful.
// Returns error if the password is incorrectly formatted (not a zero-padded 6 or non-zero-padded 8 digit number).
func (c *OTPConfig) Authenticate(password string) (bool, error) {

	switch {
	case len(password) == 6 && password[0] >= '0' && password[0] <= '9':
		break
	default:
		return false, fmt.Errorf("invalid code, exiting")
	}

	code, err := strconv.Atoi(password)

	if err != nil {
		return false, fmt.Errorf("invalid code, exiting")
	}

	var t0 int
	// assume we're on Time-based OTP
	if c.UTC {
		t0 = int(time.Now().UTC().Unix() / 30)
	} else {
		t0 = int(time.Now().Unix() / 30)
	}
	return c.checkTotpCode(t0, code), nil
}

// GenerateURI generates a URI that can be turned into a QR code
// to configure a Google Authenticator mobile app. It respects the recommendations
// on how to avoid conflicting accounts.
//
// See https://github.com/google/google-authenticator/wiki/Conflicting-Accounts
func (c *OTPConfig) GenerateURI(user string) (string, error) {
	auth := "totp/"
	issuer := "OpenX"
	q := make(url.Values)
	q.Add("secret", c.Secret)
	q.Add("issuer", issuer)
	auth += issuer + ":"

	otpString := "otpauth://" + auth + user + "?" + q.Encode()
	return otpString, nil
}
