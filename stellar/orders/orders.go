package orders

// the orders package maintains read / write operations to the orderbook
// we need an orderbook because there is no state on Stellar which makes it
// difficult for us to store this on the blockchain. Do we need to publish
// the orders / proof of these orders to the blockchain? We do need to
// publish some analogue of state..

// what do we need to store?
import (
	// "fmt"
	"encoding/binary"
	"encoding/json"
	"log"

	"github.com/boltdb/bolt" // coreOS fork of boltDB
)

type Order struct {
	Index        uint32
	PanelSize    string // size of the given panel, for diplsaying to the user who wants to bid stuff
	TotalValue   int    // the total money that we need from investors
	Location     string // where this specific solar panel is located
	MoneyRaised  int    // total money that has been raised until now
	Metadata     string // any other metadata can be stored here
	Live         bool   // check to see whether the current order is live or not
	INVAssetCode string // once all funds have been raised, we need to set assetCodes
	DEBAssetCode string // once all funds have been raised, we need to set assetCodes
	PBAssetCode  string // once all funds have been raised, we need to set assetCodes
	// Percentage raised is not stored in the database since that can be calculated by the UI
}

// do we store separate  investor and debt holder pubkeys or do we have a separate
// struct for investors, debtors and then store common fileds in them and lookup from
// there when necessary? Having a separate bucket is useful for less code
// complexity, but having less buckets might be good performance wise. But we do
// need to call the orders bucket each time and that would require locks and
// stuff

func Uint32toB(a uint32) []byte {
	// need to convert int to a byte array for indexing
	temp := make([]byte, 4)
	binary.LittleEndian.PutUint32(temp, a)
	return temp
}

func BToUint32(a []byte) uint32 {
	return binary.LittleEndian.Uint32(a)
}

// need locks over this to ensure no one's using the db while we are
func OpenDB() (*bolt.DB, error) {
	db, err := bolt.Open("yol.db", 0600, nil)
	if err != nil {
		log.Println("Couldn't open database, exiting!")
		return nil, err
	}
	return db, nil
}

// need locks over insert and retrieve operations since BOLTdb supports only
// one operation at a time.
func InsertOrder(order Order, db *bolt.DB) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("Orders")) // the orders bucket contains all our orders
		if err != nil {
			log.Fatal(err)
			return err
		}
		encoded, err := json.Marshal(order)
		if err != nil {
			log.Println("Failed to encode this data into json")
			return err
		}
		return b.Put([]byte(Uint32toB(order.Index)), encoded)
	})
	return err
}

func RetrieveOrder(key uint32, db *bolt.DB) (Order, error) {
	var rOrder Order
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Orders"))
		x := b.Get(Uint32toB(key))
		err := json.Unmarshal(x, &rOrder)
		if err != nil {
			return err
		}
		return nil
	})
	return rOrder, err
}

func DeleteOrder(key uint32, db *bolt.DB) error {
	// deleting order might be dangerous since that would mess with the retrieveAll
	// function, have it in here for now, don't do too much with it / fiox retrieve all
	// to handle this case
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Orders"))
		err := b.Delete(Uint32toB(key))
		if err != nil {
			return err
		}
		log.Println("Deleted order with key: ", key)
		return nil
	})
	return err
}

func RetrieveAll(db *bolt.DB) ([]Order, error) {
	var arr []Order
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("Orders"))
		i := uint32(1)
		for ; ; i++ {
			var rOrder Order
			x := b.Get(Uint32toB(i))
			err := json.Unmarshal(x, &rOrder)
			if err != nil && rOrder.Live == false {
				// we've reached the end of input, so this is not an error
				// ideal error would be "unexpected JSON input" or something similar
				return nil
			}
			arr = append(arr, rOrder)
		}
		return nil
	})
	return arr, err
}
