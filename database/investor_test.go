// +build all

package database

import (
	"encoding/json"
	"github.com/pkg/errors"
	"log"
	"os"
	"testing"

	consts "github.com/YaleOpenLab/openx/consts"
	utils "github.com/YaleOpenLab/openx/utils"
	"github.com/boltdb/bolt"
)

// go test -run=XXX -tags="all" -bench=.
// Benchamrking functions follow
// note that we don't have any benchmarks for recipients since most functiosn are identical to
// that of investors. Functions for users are limited. Any optimization proposed to the
// db handler functions (such as moving defers) must improve on the benchmarks defined here.

// RetrieveInvestor retrieves a particular investor indexed by key from the database
func RetrieveInvestor2(key int) (Investor, error) {
	var inv Investor
	db, err := OpenDB()
	if err != nil {
		return inv, errors.Wrap(err, "failed to open db")
	}
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(InvestorBucket)
		x := b.Get(utils.ItoB(key))
		if x == nil {
			// no investor with the specific details
			return errors.New("No investor found with required credentials")
		}
		return json.Unmarshal(x, &inv)
	})
	db.Close()
	return inv, err
}

func PopulateDB() {
	// populate the db with artifical values
	CreateHomeDir()                     // create home directory if it doesn't exist yet
	os.Remove(consts.DbDir + "/yol.db") // remove the database file, if it exists
	db, err := OpenDB()
	if err != nil {
		log.Fatal(err)
	}
	db.Close()

	for i := 0; i < 10000; i++ {
		var x Investor
		x.U.Index = i
		x.U.Username = "q"
		x.U.Pwhash = "p"
		x.Save()
		x.U.Save()
	}

	iA, err := RetrieveAllInvestors()
	if err != nil {
		log.Fatal(err)
	}

	if len(iA) != 9999 {
		log.Fatal("Couldn't populate db, quitting!")
	}
}

// run these multiple times since the first function called will obviously run slower
// due to the first caller syndrome

func BenchmarkPopulateDB(b *testing.B) {
	PopulateDB()
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		_, _ = RetrieveInvestor(i)
	}
}

// compare implementations with and without defer
func BenchmarkRetrieveWarmup(b *testing.B) {
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		_, _ = RetrieveInvestor(i)
	}
}

func BenchmarkRetrieveWD1(b *testing.B) {
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		_, _ = RetrieveInvestor(i)
	}
}

func BenchmarkRetrieveWOD1(b *testing.B) {
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		_, _ = RetrieveInvestor2(i)
	}
}

func BenchmarkRetrieveWOD2(b *testing.B) {
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		_, _ = RetrieveInvestor2(i)
	}
}

func BenchmarkRetrieveWD2(b *testing.B) {
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		_, _ = RetrieveInvestor(i)
	}
}

func BenchmarkRetrieveWOD3(b *testing.B) {
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		_, _ = RetrieveInvestor2(i)
	}
}

func BenchmarkRetrieveWD3(b *testing.B) {
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		_, _ = RetrieveInvestor(i)
	}
}

func BenchmarkRetrieveUser(b *testing.B) {
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		_, _ = RetrieveUser(i)
	}
}

// test retrieveallinvestors function. Note taht we're going to test with retrieving 10000 investors

// RetrieveAllInvestors gets a list of all investors in the database
func RetrieveAllInvestors1() ([]Investor, error) {
	// this route is broken because it reads through keys sequentially
	// need to see keys until the length of the users database
	var arr []Investor
	temp, err := RetrieveAllUsers()
	if err != nil {
		return arr, errors.Wrap(err, "failed to retrieve all users")
	}
	limit := len(temp) + 1
	db, err := OpenDB()
	if err != nil {
		return arr, errors.Wrap(err, "failed to open db")
	}

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(InvestorBucket)
		for i := 1; i < limit; i++ {
			var rInvestor Investor
			x := b.Get(utils.ItoB(i))
			if x == nil {
				// this is where the key does not exist, we search until limit, so don't error out
				continue
			}
			err := json.Unmarshal(x, &rInvestor)
			if err != nil {
				// error in unmarshalling this struct, error out
				return errors.Wrap(err, "failed to unmarshal json")
			}
			arr = append(arr, rInvestor)
		}
		return nil
	})
	db.Close()
	return arr, err
}

func BenchmarkRetrieveAllWOD1(b *testing.B) {
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		_, _ = RetrieveAllInvestors1()
	}
}

func BenchmarkRetrieveAllWD1(b *testing.B) {
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		_, _ = RetrieveAllInvestors()
	}
}

func BenchmarkRetrieveAllWD2(b *testing.B) {
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		_, _ = RetrieveAllInvestors()
	}
}

func BenchmarkRetrieveAllWOD2(b *testing.B) {
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		_, _ = RetrieveAllInvestors1()
	}
}

func BenchmarkValidateUser(b *testing.B) {
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		_, _ = ValidateUser("q", "p")
	}
}

func BenchmarkValidateInvestor(b *testing.B) {
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		_, _ = ValidateInvestor("q", "p")
	}
}

func BenchmarkInvReputation(b *testing.B) {
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		_ = ChangeInvReputation(i, float64(i))
	}
}

func BenchmarkRetrieveAllTRInvestors(b *testing.B) {
	b.ResetTimer()
	for i := 1; i < 10; i++ {
		_, _ = TopReputationInvestors()
	}
}
