package main

// test.go runs the PoC stellar implementation calling various functions
import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"syscall"

	assets "github.com/YaleOpenLab/smartPropertyMVP/stellar/assets"
	database "github.com/YaleOpenLab/smartPropertyMVP/stellar/database"
	rpc "github.com/YaleOpenLab/smartPropertyMVP/stellar/rpc"
	stablecoin "github.com/YaleOpenLab/smartPropertyMVP/stellar/stablecoin"
	utils "github.com/YaleOpenLab/smartPropertyMVP/stellar/utils"
	xlm "github.com/YaleOpenLab/smartPropertyMVP/stellar/xlm"
	flags "github.com/jessevdk/go-flags"
	"golang.org/x/crypto/ssh/terminal"
)

var opts struct {
	// Slice of bool will append 'true' each time the option
	// is encountered (can be set multiple times, like -vvv)
	// TOOD: define default values for each and then use them if not passed
	Verbose []bool `short:"v" long:"verbose" description:"Show verbose debug information"`
	// InvAmount int    `short:"i" description:"Desired investment" required:"true"`
	InvAmount int `short:"i" description:"Desired investment"`
	// RecYears  int    `short:"r" description:"Number of years the recipient wants to repay in. Can be 3, 5 or 7 years." required:"true"`
	RecYears int    `short:"r" description:"Number of years the recipient wants to repay in. Can be 3, 5 or 7 years."`
	Port     string `short:"p" description:"The port on which the server runs on"`
}

func ValidateInputs() {
	if (opts.RecYears != 0) && !(opts.RecYears == 3 || opts.RecYears == 5 || opts.RecYears == 7) {
		// right now payoff periods are limited, I guess they don't need to be,
		// but in this case jsut are
		log.Fatal(fmt.Errorf("Number of years not supported"))
	}
}

func main() {
	var err error
	_, err = flags.ParseArgs(&opts, os.Args)
	if err != nil {
		log.Fatal(err)
	}

	// Open the database
	allOrders, err := database.RetrieveAllOrders()
	if err != nil {
		log.Println("Error retrieving all orders from the database")
	}

	if len(allOrders) == 0 {
		log.Println("Populating database with test values")
		err = database.InsertDummyData()
		if err != nil {
			log.Fatal(err)
		}
	}
	// TODO: how much do we pay the investor? how does it work
	// Do we sell the REC created from the solar panels only to the investor? If so,
	// isn't that enough to propel investment in the solar contract itself?
	// TODO: need a server to run a public stellar node to test out stuff
	// things to consider:
	// while an investor signs up on our platform, do we send them 10 XLM free?
	// do we charge investors to be on our platform? if not, we shouldn't ideally
	// be sending them free XLM. also, should the platform have some function for
	// withdrawing XLM? if so, we'll become an exchange of sorts and have some
	// legal stuff there. If not, we'll just be a custodian and would not have
	// too much to consider on our side
	// need to spin up a local stellar node and test if things run fine if we just
	// change the API mapping
	// need to create different entities and create db mappings for them.
	// need to update collections to directly hold orders, similar to the investor
	// class that we have already
	// need to implement the contract stuff as described earlier, so that people
	// can advertise bids, get paid for it, etc.
	// move current number of years metric to a separate package since that is
	// more suitable for a model like affordable housing.
	// look into what kind of data we get from the pi and checkout pi specific code
	// to see if we can get something from there.

	fmt.Println("------------STELLAR HOUSE INVESTMENT CLI INTERFACE------------")
	ValidateInputs()

	// setup issuer account if the platform doesn't  already exist
	// check whether the platform exists
	test, err := database.RetrievePlatform()
	if err != nil {
		log.Fatal(err)
	}
	if len(test.PublicKey) == 0 {
		// weird way to test, but still
		// this is the first time we're initializing a platform
		log.Println("Creating a new platform")
		platform, err := database.NewPlatform()
		if err != nil {
			log.Fatal(err)
		}
		// insert this into the database
		err = database.InsertPlatform(platform)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		log.Println("Platform already exists, using existing one")
	}
	platform, err := database.RetrievePlatform()
	if err != nil {
		log.Fatal(err)
	}
	// ask for the platform's password
	// now here, we must decrypt the seed before using it in other places
	fmt.Printf("%s: ", "ENTER PASSWORD TO UNLOCK THE PLATFORM")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	psPassword := string(bytePassword)
	platformSeed := database.GetSeedFromEncryptedSeed("seed.hex", psPassword)
	// init stablecoin stuff
	err = stablecoin.InitStableCoin()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(stablecoin.StableUSD)
	go stablecoin.ListenForPayments()
	// don't have an error catching thing because if this fails, the platform should
	// not initialize
	// insert an investor with the relevant details
	// add dummy investor and recipient data for the demo
	// uname: john, password: password
	/*
		nInvestor, err := database.NewInvestor("john",
			"e9a75486736a550af4fea861e2378305c4a555a05094dee1dca2f68afea49cc3a50e8de6ea131ea521311f4d6fb054a146e8282f8e35ff2e6368c1a62e909716",
			"John", true)
		err = database.InsertInvestor(nInvestor)
		if err != nil {
			log.Fatal(err)
		}

		// need to add a dummy recipient here as well
		// uname: martin, password:password
		nRecipient, err := database.NewRecipientWithoutSeed("martin",
			"8a56bac869374c669443a1626ff0967af258123f83faf6b55e31dd541e6bbd90308a3385713294bf2e8861bc8cf8f8feda41f9c4db19d5811a6b5de85eac9870",
			"Martin")
		err = database.InsertRecipient(nRecipient)
		if err != nil {
			log.Fatal(err)
		}
	*/
	// need to ask for user role as well here, to know whether the user is an investor
	// or recipient so that we can show both sides
	// After this, ask what the user wants to do - there are roughly three options:
	// 1. Create a new investor account
	// 2. Create a new recipient account
	// 3. Login (Are you an investor / recipient)
	fmt.Println("------WHAT DO YOU WANT TO DO?------")
	fmt.Println("1. CREATE A NEW INVESTOR ACCOUNT")
	fmt.Println("2. CREATE A NEW RECIPIENT ACCOUNT")
	fmt.Println("deafult: ALREADY HAVE AN ACCOUNT")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	if scanner.Err() != nil {
		fmt.Println("Couldn't read user input")
		return
	}
	switch scanner.Text() {
	case "1":
		log.Println("You have chosen to create a new investor account, welcome")
		log.Println("ENTER YOUR REAL NAME")
		scanner.Scan()
		if scanner.Err() != nil {
			fmt.Println("Couldn't read user input")
			break
		}
		invName := scanner.Text()

		log.Println("ENTER YOUR USERNAME")
		scanner.Scan()
		if scanner.Err() != nil {
			fmt.Println("Couldn't read user input")
			break
		}
		invLoginUserName := scanner.Text()

		log.Println("ENTER DESIRED PASSWORD, YOU WILL NOT BE ASKED TO CONFIRM THIS")
		bytePassword, err = terminal.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		tempString := string(bytePassword)
		invLoginPassword := utils.SHA3hash(tempString)

		inv, err := database.NewInvestor(invLoginUserName, invLoginPassword, invName)
		if err != nil {
			log.Println("FAILED TO SETUP ACCOUNT, TRY AGAIN")
			break
		}
		err = database.InsertInvestor(inv)
		if err != nil {
			log.Println("FAILED TO SETUP ACCOUNT, TRY AGAIN")
			break
		}
		// need to send fudns to this guy so that he can setup trustlines
		break
	case "2":
		log.Println("You have chosen to create a new recipient account, welcome")
		log.Println("ENTER YOUR REAL NAME")
		scanner.Scan()
		if scanner.Err() != nil {
			fmt.Println("Couldn't read user input")
			break
		}
		invName := scanner.Text()

		log.Println("ENTER YOUR USERNAME")
		scanner.Scan()
		if scanner.Err() != nil {
			fmt.Println("Couldn't read user input")
			break
		}
		invLoginUserName := scanner.Text()

		log.Println("ENTER DESIRED PASSWORD, YOU WILL NOT BE ASKED TO CONFIRM THIS")
		bytePassword, err = terminal.ReadPassword(int(syscall.Stdin))
		fmt.Println()
		tempString := string(bytePassword)
		invLoginPassword := utils.SHA3hash(tempString)

		inv, err := database.NewRecipient(invLoginUserName, invLoginPassword, invName)
		if err != nil {
			log.Println("FAILED TO SETUP ACCOUNT, TRY AGAIN")
			break
		}
		err = database.InsertRecipient(inv)
		if err != nil {
			log.Println("FAILED TO SETUP ACCOUNT, TRY AGAIN")
			break
		}
		break
	default:
		// don't add the entire file as a switch case because it would be ugly. we can
		// fall through, shouldn't be an issue
	}
	fmt.Println("---ARE YOU AN INVESTOR (I) OR RECIPIENT (R)? ---")
	scanner.Scan()
	rbool := false
	if scanner.Text() == "I" || scanner.Text() == "i" {
		fmt.Println("WELCOME BACK INVESTOR")
	} else if scanner.Text() == "R" || scanner.Text() == "r" {
		fmt.Println("WELCOME BACK RECIPIENT")
		rbool = true
	} else {
		log.Fatal("INVALID INPUT, EXITING!")
	}
	// ask for username and password combo here
	fmt.Printf("%s", "ENTER YOUR USERNAME: ")
	scanner.Scan()
	if scanner.Err() != nil {
		fmt.Println("Couldn't read user input")
		return
	}
	invLoginUserName := scanner.Text() // read user input regarding which option
	fmt.Printf("%s", "ENTER YOUR PASSWORD: ")
	bytePassword, err = terminal.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	invLoginPassword := utils.SHA3hash(string(bytePassword))
	// check for ibool vs rbool here
	if rbool {
		// handle the recipient case here because its simpler
		recipient, err := database.ValidateRecipient(invLoginUserName, invLoginPassword)
		if err != nil {
			log.Fatal("had trouble retrieving the username")
		}
		log.Println("RECIPIENT IS: ", recipient)
		if recipient.U.LoginPassword != invLoginPassword { // should rework to check the password, this is just a temp hack
			log.Printf("INGLOGIN: %s, LOGINP: %s", invLoginPassword, recipient.U.LoginPassword)
			log.Fatal("Passwords don't match")
		}
		// at this point, we have verified the recipient
		// have a for loop here with various options
		for {
			fmt.Println("------------RECIPIENT INTERFACE------------")
			fmt.Println("----CHOOSE ONE OF THE FOLLOWING OPTIONS----")
			fmt.Println("  1. Display all Open Assets")
			fmt.Println("  2. Display my Profile")
			fmt.Println("  3. Payback towards an Asset")
			fmt.Println("  4. Exchange XLM for USD")
			fmt.Println("  default: Exit")
			scanner.Scan()
			if scanner.Err() != nil {
				fmt.Println("Couldn't read user input")
				break
			}
			menuInput, err := strconv.Atoi(scanner.Text())
			if err != nil {
				log.Fatal(err)
			}
			switch menuInput {
			case 1:
				fmt.Println("------------LIST OF ALL AVAILABLE ASSETS------------")
				allOrders, err := database.RetrieveAllOrders()
				if err != nil {
					log.Println("Error retrieving all orders from the database")
				}
				database.PrettyPrintOrders(allOrders)
				break
			case 2:
				database.PrettyPrintRecipient(recipient)
				break
			case 3:
				// regarding payback, we need to first check if the STABLEUSD asset
				// is present on our account. If yes, we assume that it is 1:1 with USD
				// and then we need to pay the same balance in DEBTokens.
				// one could transfer the STRONGUSD token directly as well but we wouldn't
				// have a proof of token ownership / receipt
				// TODO: This does not account for the XLM/USD exchange and a recipient
				// can payback infinitely for now
				database.PrettyPrintPBOrders(recipient.ReceivedOrders)
				fmt.Println("WHICH ORDER DO YOU WANT TO PAY BACK TOWARDS? (ENTER ORDER NUMBER)")
				scanner.Scan()
				if scanner.Err() != nil {
					fmt.Println("Couldn't read user input")
					break
				}
				// user input must be an integer, else quit
				orderNumber, err := strconv.Atoi(scanner.Text())
				if err != nil {
					log.Println("INPUT NOT AN INTEGER, TRY AGAIN")
					continue
				}
				// check if we can get the roder using the order number that we have here
				rtOrder, err := database.RetrieveOrder(uint32(orderNumber))
				if err != nil {
					log.Println("Couldn't retrieve order, try again!")
					continue
				}
				// so we can retrieve the order using the order Index, nice
				database.PrettyPrintPBOrder(rtOrder)
				fmt.Println("HOW MUCH DO YOU WANT TO PAYBACK?")
				scanner.Scan()
				if scanner.Err() != nil {
					fmt.Println("Couldn't read user input")
					break
				}
				// user input must be an integer, else quit
				pbAmountS := scanner.Text()
				_, err = strconv.Atoi(pbAmountS) // TODO: assumes whole numbers
				if err != nil {
					log.Println("PAYBACK AMOUNT NOT AN INTEGER, TRY AGAIN")
					continue
				}
				fmt.Printf(" DO YOU WANT TO CONFIRM THAT YOU WANT TO PAYBACK %s TOWARDS THIS ORDER? (PRESS N IF YOU DON'T WANT TO)\n", pbAmountS)
				scanner.Scan()
				if scanner.Text() == "N" || scanner.Text() == "n" {
					fmt.Println("YOU HAVE DECIDED TO CANCEL THIS ORDER")
					break
				}
				fmt.Printf("PAYING BACK %s TOWARDS ORDER NUMBER: %d\n", pbAmountS, rtOrder.Index) // use the rtOrder here instead of using orderNumber from long ago
				// now we need to call back the payback function to payback the asset
				// TODO: need to check if the recipient has the required USDTokens.
				// Here, we will simply payback the DEBTokens that was sent to us earlier
				if rtOrder.DEBAssetCode == "" {
					log.Fatal("Asset not found")
				}
				err = recipient.Payback(rtOrder, rtOrder.DEBAssetCode, platform.PublicKey, pbAmountS)
				// TODO: right now, the payback asset directly sends back, change
				if err != nil {
					log.Println("PAYBACK TX FAILED, PLEASE TRY AGAIN!")
					break
				}
				fmt.Println("UPDATED ORDER: ")
				// check if we can get the roder using the order number that we have here
				rtOrder, err = database.RetrieveOrder(uint32(orderNumber))
				if err != nil {
					log.Println("Couldn't retrieve updated order, check again!")
					continue
				}
				// we should update the local slice to keep track of the changes here
				recipient.UpdateOrderSlice(rtOrder)
				// so we can retrieve the order using the order Index, nice
				database.PrettyPrintOrder(rtOrder)
				// print the order in a nice way
				break
			case 4:
				log.Println("Enter the amount you want to convert into STABLEUSD")
				scanner.Scan()
				convAmount := scanner.Text()
				if utils.StringToFloat(convAmount) == 0 {
					log.Println("Amount entered is not a float, quitting")
					break
				}
				hash, err := assets.TrustAsset(stablecoin.StableUSD, "1000000000", recipient.U.PublicKey, recipient.U.Seed)
				if err != nil {
					log.Fatal(err)
				}
				log.Println("tx hash for trusting stableUSD: ", hash)
				// now send coins across and see if our tracker detects it
				_, hash, err = xlm.SendXLM(stablecoin.Issuer.PublicKey, convAmount, recipient.U.Seed)
				if err != nil {
					log.Fatal(err)
				}
				log.Println("tx hash for sent xlm: ", hash, "pubkey: ", recipient.U.PublicKey)
				break
			default:
				// check whether he wants to go back to the display all screen again
				fmt.Println("DO YOU REALLY WANT TO EXIT? (PRESS Y TO CONFIRM)")
				scanner.Scan()
				if scanner.Text() == "Y" || scanner.Text() == "y" {
					fmt.Println("YOU HAVE DECIDED TO EXIT")
					log.Fatal("")
				}
				break
			}
		}
		database.PrettyPrintRecipient(recipient)
		return
	}

	investor, err := database.ValidateInvestor(invLoginUserName, invLoginPassword)
	if err != nil {
		log.Fatal("had trouble retrieving user from db, Username / password doesn't match")
	}

	for {
		// Main investor loop
		// password is password
		// ask for the investor username and password
		/*
			err = database.DeleteOrder(1, db)
			if err != nil {
				log.Println("Error deleting entry 1 from database")
			}
			err = database.DeleteOrder(2, db)
			if err != nil {
				log.Println("Error deleting entry 2 from database")
			}
		*/
		// have a menu here that will ask you what you want to do
		fmt.Println("------------INVESTOR INTERFACE------------")
		fmt.Println("----CHOOSE ONE OF THE FOLLOWING OPTIONS----")
		fmt.Println("  1. Display all Open Assets")
		fmt.Println("  2. Display my Profile")
		fmt.Println("  3. Invest in an Asset")
		fmt.Println("  4. Display All Balances")
		fmt.Println("  5. Exchange XLM for USD")
		fmt.Println("  default: Exit")
		scanner.Scan()
		if scanner.Err() != nil {
			fmt.Println("Couldn't read user input")
			break
		}
		menuInput, err := strconv.Atoi(scanner.Text())
		if err != nil {
			log.Fatal(err)
		}
		switch menuInput {
		case 1:
			fmt.Println("------------LIST OF ALL AVAILABLE ASSETS------------")
			allOrders, err := database.RetrieveAllOrders()
			if err != nil {
				log.Println("Error retrieving all orders from the database")
			}
			database.PrettyPrintOrders(allOrders)
			break
		case 2:
			database.PrettyPrintInvestor(investor)
			break
		case 3:
			fmt.Println("----WHICH ASSET DO YOU WANT TO INVEST IN? (ENTER NUMBER WITHOUT SPACES)----")
			scanner.Scan()
			if scanner.Err() != nil {
				fmt.Println("Couldn't read user input")
				break
			}
			// they want to choose
			// also check whether received user input is an integer
			uInput, err := strconv.Atoi(scanner.Text())
			if err != nil {
				log.Fatal("user input is not a number")
			}
			// now the user has decided to invest in the asset with index uInput
			// we need to retrieve the order and ask for confirmation
			uOrder, err := database.RetrieveOrder(uint32(uInput))
			if err != nil {
				log.Fatal("Order with specified index not found in the database")
			}
			database.PrettyPrintOrder(uOrder)
			fmt.Println(" HOW MUCH DO YOU WANT TO INVEST?")
			scanner.Scan()
			investedAmountS := scanner.Text()
			_, err = strconv.Atoi(investedAmountS)
			if err != nil {
				fmt.Println("AMOUNT INVESTED IS NOT AN INTEGER, EXITING!")
				break
			}
			fmt.Println(" DO YOU WANT TO CONFIRM THIS ORDER? (PRESS N IF YOU DON'T WANT TO)")
			scanner.Scan()
			if scanner.Text() == "N" || scanner.Text() == "n" {
				fmt.Println("YOU HAVE DECIDED TO CANCEL THIS ORDER")
				break
			}
			// when I am creating an account, I will have a PublicKey and Seed, so
			// don't need them here
			// check whether the investor has XLM already
			balance, err := xlm.GetXLMBalance(platform.PublicKey)
			// balance is in string, convert to int
			balanceI := utils.StringToFloat(balance)
			log.Println("Platform's balance is: ", balanceI)
			if balanceI < 21 { // 1 to account for fees
				// get coins if balance is this low
				log.Println("Refilling platform balance")
				err := xlm.GetXLM(platform.PublicKey)
				// TODO: in future, need to refill platform sufficiently well and interact
				// with a cold wallet that we have previously set
				if err != nil {
					log.Fatal(err)
				}
			}

			balance, err = xlm.GetXLMBalance(platform.PublicKey)
			log.Println("Platform balance updated is: ", balance)
			log.Printf("Platform seed is: %s and platform's publicKey is %s", platformSeed, platform.PublicKey)
			log.Println("Investor's publickey is: ", investor.U.PublicKey)
			balance, err = xlm.GetXLMBalance(investor.U.PublicKey)
			if balance == "" {
				// means we need to setup an account first
				// Generating a keypair on stellar doesn't mean that you can send funds to it
				// you need to call the CreateAccount method in order to be able to send funds
				// to it
				log.Println("Investor balance empty, refilling!")
				_, _, err = xlm.SendXLMCreateAccount(investor.U.PublicKey, "10", platformSeed)
				if err != nil {
					log.Println("Investor Account doesn't have funds")
					log.Fatal(err)
				}
			}
			// balance is in string, convert to float
			balance, err = xlm.GetXLMBalance(investor.U.PublicKey)
			balanceI = utils.StringToFloat(balance)
			log.Println("Investor balance is: ", balanceI)
			if balanceI < 3 { // to setup trustlines
				_, _, err = xlm.SendXLM(investor.U.PublicKey, "10", platformSeed)
				if err != nil {
					log.Println("Investor Account doesn't have funds")
					log.Fatal(err)
				}
			}

			recipient := uOrder.OrderRecipient
			// from here on, reference recipient
			balance, err = xlm.GetXLMBalance(recipient.U.PublicKey)
			if balance == "" {
				// means we need to setup an account first
				// Generating a keypair on stellar doesn't mean that you can send funds to it
				// you need to call the CreateAccount method in order to be able to send funds
				// to it
				_, _, err = xlm.SendXLMCreateAccount(recipient.U.PublicKey, "10", platformSeed)
				if err != nil {
					log.Println("Recipient Account doesn't have funds")
					log.Fatal(err)
				}
			}
			balance, err = xlm.GetXLMBalance(recipient.U.PublicKey)
			// balance is in string, convert to float
			balanceI = utils.StringToFloat(balance)
			log.Println("Recipient balance is: ", balanceI)
			if balanceI < 3 { // to setup trustlines
				_, _, err = xlm.SendXLM(recipient.U.PublicKey, "10", platformSeed)
				if err != nil {
					log.Println("Recipient Account doesn't have funds")
					log.Fatal(err)
				}
			}
			log.Println("The issuer's public key and private key are: ", platform.PublicKey, " ", platformSeed)
			log.Println("The investor's public key and private key are: ", investor.U.PublicKey, " ", investor.U.Seed)
			log.Println("The recipient's public key and private key are: ", recipient.U.PublicKey, " ", recipient.U.Seed)

			log.Println(&platform, platformSeed, &investor, &recipient, investedAmountS, uOrder)
			// so now we have three entities setup, so we create the assets and invest in them
			cOrder, err := assets.InvestInOrder(&platform, platformSeed, &investor, &recipient, investedAmountS, uOrder) // assume payback period is 5
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println("YOUR ORDER HAS BEEN CONFIRMED: ")
			database.PrettyPrintOrder(cOrder)
			fmt.Println("PLEASE CHECK A BLOCKHAIN EXPLORER TO CONFIRM BALANCES TO CONFIRM: ")
			fmt.Println("https://testnet.steexp.com/account/" + investor.U.PublicKey + "#balances")
			break
		case 4:
			balances, err := xlm.GetAllBalances(investor.U.PublicKey)
			if err != nil {
				log.Fatal(err)
			}
			// need to pr etty print this, experiment out with stuff
			xlm.PrettyPrintBalances(balances)
			break
		case 5:
			// this should be expanded in the future to make use of the inbuilt DEX
			// on stellar (checkout stellarterm)
			log.Println("Enter the amount you want to convert into STABLEUSD")
			// this would also mean that you need to check whether we have the balance
			// here and then proceed further
			scanner.Scan()
			convAmount := scanner.Text()
			if utils.StringToFloat(convAmount) == 0 {
				log.Println("Amount entered is not a float, quitting")
				break
			}
			// pay this stablecoin account to test
			/*
				reuse this to write unit tests
				privkey, pubkey, err := xlm.GetKeyPair()
				if err != nil {
					log.Fatal(err)
				}
				log.Println("created a new pirvkey, pubkrey pair", privkey, pubkey)
				// setup this account for testing
				err = xlm.GetXLM(pubkey)
				if err != nil {
					log.Fatal(err)
				}
				// first trust the stablecoin issuer
			*/
			// maybe don't trust asset again when you've trusted it already? check if that's
			// possible and save on the tx fee for a single transaction
			hash, err := assets.TrustAsset(stablecoin.StableUSD, "10000", investor.U.PublicKey, investor.U.Seed)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("tx hash for trusting stableUSD: ", hash)
			// now send coins across and see if our tracker detects it
			_, hash, err = xlm.SendXLM(stablecoin.Issuer.PublicKey, convAmount, investor.U.Seed)
			if err != nil {
				log.Fatal(err)
			}

			log.Println("tx hash for sent xlm: ", hash, "pubkey: ", investor.U.PublicKey)
			rpc.StartServer("8080") // run this in order to check whether the go routine is running
			break
		default:
			// check whether he wants to go back to the display all screen again
			fmt.Println("DO YOU REALLY WANT TO EXIT? (PRESS Y TO CONFIRM)")
			scanner.Scan()
			if scanner.Text() == "Y" || scanner.Text() == "y" {
				fmt.Println("YOU HAVE DECIDED TO EXIT")
				log.Fatal("")
			}
		} // end of switch
	}
	log.Fatal("")
	// we now have the order the user wants to confirm, pretty print this order
	log.Fatal("All good")
	// start the rpc server
	rpc.StartServer("8080")    // run this in order to check whether the go routine is running
	rpc.StartServer(opts.Port) // this must be towards the end
}
