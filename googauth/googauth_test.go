// +build all travis

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
	"fmt"
	"testing"
	"time"
)

// Test vectors via:
// http://code.google.com/p/google-authenticator/source/browse/libpam/pam_google_authenticator_unittest.c
// https://google-authenticator.googlecode.com/hg/libpam/totp.html

var codeTests = []struct {
	secret string
	value  int64
	code   int
}{
	{"2SH3V3GDW7ZNMGYE", 1, 293240},
	{"2SH3V3GDW7ZNMGYE", 5, 932068},
	{"2SH3V3GDW7ZNMGYE", 10000, 50548},
}

func TestCode(t *testing.T) {

	for _, v := range codeTests {
		c := ComputeCode(v.secret, v.value)

		if c != v.code {
			t.Errorf("computeCode(%s, %d): got %d expected %d\n", v.secret, v.value, c, v.code)
		}

	}
}

func TestTotpCode(t *testing.T) {

	var cotp OTPConfig

	// reuse our test values from above
	cotp.Secret = "2SH3V3GDW7ZNMGYE"
	cotp.WindowSize = 5

	var windowTest = []struct {
		code   int
		t0     int
		result bool
	}{
		{50548, 9997, false},
		{50548, 9998, true},
		{50548, 9999, true},
		{50548, 10000, true},
		{50548, 10001, true},
		{50548, 10002, true},
		{50548, 10003, false},
	}

	for i, s := range windowTest {
		r := cotp.checkTotpCode(s.t0, s.code)
		if r != s.result {
			t.Errorf("counterCode(%d) (step %d) failed: got %t expected %t", s.code, i, r, s.result)
		}
	}

	cotp.DisallowReuse = make([]int, 0)
	var noreuseTest = []struct {
		code       int
		t0         int
		result     bool
		disallowed []int
	}{
		{50548 /* 10000 */, 9997, false, []int{}},
		{50548 /* 10000 */, 9998, true, []int{10000}},
		{50548 /* 10000 */, 9999, false, []int{10000}},
		{478726 /* 10001 */, 10001, true, []int{10000, 10001}},
		{646986 /* 10002 */, 10002, true, []int{10000, 10001, 10002}},
		{842639 /* 10003 */, 10003, true, []int{10001, 10002, 10003}},
	}

	for i, s := range noreuseTest {
		r := cotp.checkTotpCode(s.t0, s.code)
		if r != s.result {
			t.Errorf("timeCode(%d) (step %d) failed: got %t expected %t", s.code, i, r, s.result)
		}
		if len(cotp.DisallowReuse) != len(s.disallowed) {
			t.Errorf("timeCode(%d) (step %d) failed: disallowReuse len mismatch: got %d expected %d", s.code, i, len(cotp.DisallowReuse), len(s.disallowed))
		} else {
			same := true
			for j := range s.disallowed {
				if s.disallowed[j] != cotp.DisallowReuse[j] {
					same = false
				}
			}
			if !same {
				t.Errorf("timeCode(%d) (step %d) failed: disallowReused: got %v expected %v", s.code, i, cotp.DisallowReuse, s.disallowed)
			}
		}
	}
}

func TestAuthenticate(t *testing.T) {

	otpconf := &OTPConfig{
		Secret:     "2SH3V3GDW7ZNMGYE",
		WindowSize: 1,
	}

	type attempt struct {
		code   string
		result bool
	}

	var attempts = []attempt{
		{"foobar", false},  // not digits
		{"1fooba", false},  // not valid number
		{"1111111", false}, // bad length
	}

	for _, a := range attempts {
		r, _ := otpconf.Authenticate(a.code)
		if r != a.result {
			t.Errorf("bad result from code=%s: got %t expected %t\n", a.code, r, a.result)
		}
	}

	// I haven't mocked the clock, so we'll just compute one
	var t0 int64
	if otpconf.UTC {
		t0 = int64(time.Now().UTC().Unix() / 30)
	} else {
		t0 = int64(time.Now().Unix() / 30)
	}
	c := ComputeCode(otpconf.Secret, t0)
	code := fmt.Sprintf("%06d", c)

	attempts = []attempt{
		{code + "1", false},
		{code, true},
	}

	for _, a := range attempts {
		r, _ := otpconf.Authenticate(a.code)
		if r != a.result {
			t.Errorf("bad result from code=%s: got %t expected %t\n", a.code, r, a.result)
		}

		otpconf.UTC = true
		r, _ = otpconf.Authenticate(a.code)
		if r != a.result {
			t.Errorf("bad result from code=%s: got %t expected %t\n", a.code, r, a.result)
		}
		otpconf.UTC = false
	}

}

func TestGenerateURI(t *testing.T) {
	otpconf := OTPConfig{
		Secret: "x",
	}

	cases := []struct {
		user string
		out  string
	}{
		{"test", "otpauth://totp/OpenX:test?issuer=OpenX&secret=x"},
		{"blah", "otpauth://totp/OpenX:blah?issuer=OpenX&secret=x"},
	}

	for i, c := range cases {
		otpString, err := otpconf.GenerateURI(c.user)
		if err != nil {
			t.Fatal(err)
		}
		if otpString != c.out {
			t.Errorf("%d: want %q, got %q", i, c.out, otpString)
		}
	}
}
