package main

// test.go runs the PoC stellar implementation calling various functions
import (
	"fmt"
	"log"
	"os"

	assets "github.com/YaleOpenLab/smartPropertyMVP/stellar/assets"
	consts "github.com/YaleOpenLab/smartPropertyMVP/stellar/consts"
	database "github.com/YaleOpenLab/smartPropertyMVP/stellar/database"
	rpc "github.com/YaleOpenLab/smartPropertyMVP/stellar/rpc"
	stablecoin "github.com/YaleOpenLab/smartPropertyMVP/stellar/stablecoin"
	utils "github.com/YaleOpenLab/smartPropertyMVP/stellar/utils"
	xlm "github.com/YaleOpenLab/smartPropertyMVP/stellar/xlm"
	flags "github.com/jessevdk/go-flags"
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

func main() {
	var err error
	_, err = flags.ParseArgs(&opts, os.Args)
	if err != nil {
		log.Fatal(err)
	}

	// Open the database
	platformPublicKey, platformSeed, err := StartPlatform()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("PLATFORM SEED IS: %s\n PLATFORM PUBLIC KEY IS: %s", platformSeed, platformPublicKey)
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
	// need to update collections to directly hold orders, similar to the investor
	// class that we have already
	// need to implement the contract stuff as described earlier, so that people
	// can advertise bids, get paid for it, etc.
	// move current number of years metric to a separate package since that is
	// more suitable for a model like affordable housing.
	// look into what kind of data we get from the pi and checkout pi specific code
	// to see if we can get something from there.

	fmt.Println("------------STELLAR HOUSE INVESTMENT CLI INTERFACE------------")

	// init stablecoin stuff
	err = stablecoin.InitStableCoin()
	if err != nil {
		log.Fatal(err)
	}

	// start a goroutine to listen for stablecoin payments and issuance
	go stablecoin.ListenForPayments()
	// don't have an error catching thing because if this fails, the platform should
	// not initialize
	// insert an investor with the relevant details
	// add dummy investor and recipient data for the demo
	// uname: john, password: password
	// need to ask for user role as well here, to know whether the user is an investor
	// or recipient so that we can show both sides
	// After this, ask what the user wants to do - there are roughly three options:
	// 1. Create a new investor account
	// 2. Create a new recipient account
	// 3. Login (Are you an investor / recipient)
	fmt.Println("------WHAT DO YOU WANT TO DO?------")
	fmt.Println("1. CREATE A NEW INVESTOR ACCOUNT")
	fmt.Println("2. CREATE A NEW RECIPIENT ACCOUNT")
	fmt.Println("3: ALREADY HAVE AN ACCOUNT")
	opt, err := utils.ScanForInt()
	if err != nil {
		fmt.Println("Couldn't read user input")
		return
	}
	switch opt {
	case 1:
		log.Println("You have chosen to create a new investor account, welcome")
		err := NewInvestorPrompt()
		if err != nil {
			log.Println(err)
			return
		}
	case 2:
		log.Println("You have chosen to create a new recipient account, welcome")
		err := NewRecipientPrompt()
		if err != nil {
			log.Println(err)
			return
		}
	default:
	}

	investor, recipient, isRecipient, err := LoginPrompt()
	if err != nil {
		log.Println(err)
		return
	}
	// check if the user is a recipient here
	if isRecipient {
		// we already have the recipient, so no need to make a call to the database
		// have a for loop here with various options
		for {
			fmt.Println("------------RECIPIENT INTERFACE------------")
			fmt.Println("----CHOOSE ONE OF THE FOLLOWING OPTIONS----")
			fmt.Println("  1. Display all Open Orders")
			fmt.Println("  2. Display my Profile")
			fmt.Println("  3. Payback towards an Order")
			fmt.Println("  4. Exchange XLM for USD")
			fmt.Println("  5. Finalize a specific contract")
			fmt.Println("  default: Exit")
			optI, err := utils.ScanForInt()
			if err != nil {
				fmt.Println("Couldn't read user input")
				break
			}
			switch optI {
			case 1:
				fmt.Println("------------LIST OF ALL AVAILABLE ORDERS------------")
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
				database.PrettyPrintPBOrders(recipient.ReceivedOrders)
				fmt.Println("WHICH ORDER DO YOU WANT TO PAY BACK TOWARDS? (ENTER ORDER NUMBER)")
				orderNumber, err := utils.ScanForInt()
				if err != nil {
					log.Println("INPUT NOT AN INTEGER, TRY AGAIN")
					continue
				}
				// check if we can get the order using the order number that we have here
				rtOrder, err := database.RetrieveOrder(orderNumber)
				if err != nil {
					log.Println("Couldn't retrieve order, try again!")
					continue
				}
				// so we can retrieve the order using the order Index, nice
				database.PrettyPrintPBOrder(rtOrder)
				fmt.Println("HOW MUCH DO YOU WANT TO PAYBACK?")
				paybackAmount, err := utils.ScanForStringWithCheckI()
				if err != nil {
					log.Println(err)
					break
				}
				fmt.Printf(" DO YOU WANT TO CONFIRM THAT YOU WANT TO PAYBACK %s TOWARDS THIS ORDER? (PRESS N IF YOU DON'T WANT TO)\n", paybackAmount)
				confirmOpt, err := utils.ScanForString()
				if err != nil {
					log.Println(err)
					break
				}
				if confirmOpt == "N" || confirmOpt == "n" {
					fmt.Println("YOU HAVE DECIDED TO CANCEL THIS ORDER")
					break
				}
				fmt.Printf("PAYING BACK %s TOWARDS ORDER NUMBER: %d\n", paybackAmount, rtOrder.Index) // use the rtOrder here instead of using orderNumber from long ago
				// now we need to call back the payback function to payback the asset
				// Here, we will simply payback the DEBTokens that was sent to us earlier
				if rtOrder.DEBAssetCode == "" {
					log.Fatal("Order not found")
				}
				err = recipient.Payback(rtOrder, rtOrder.DEBAssetCode, platformPublicKey, paybackAmount)
				// TODO: right now, the payback asset directly sends back, change
				if err != nil {
					log.Println("PAYBACK TX FAILED, PLEASE TRY AGAIN!")
					break
				}
				// now send back the PBToken from the platform to the issuer
				// this function is optional and can be deleted in case we don't need PBAssets
				err = assets.SendPBAsset(rtOrder, recipient.U.PublicKey, paybackAmount, platformSeed, platformPublicKey)
				if err != nil {
					log.Println("PBAsset sending back FAILED, PLEASE TRY AGAIN!")
					break
				}
				// check if we can get the order using the order number that we have here
				rtOrder, err = database.RetrieveOrder(orderNumber)
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
				convAmount, err := utils.ScanForStringWithCheckF()
				if err != nil {
					log.Println(err)
					break
				}
				hash, err := assets.TrustAsset(stablecoin.StableUSD, consts.StablecoinTrustLimit, recipient.U.PublicKey, recipient.U.Seed)
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
			case 5:
				// we shouild  finalize the contract that we want
				// can be imagined as some sort of voting mechanism to choose the winning
				// contract.
				// now we display a list of options for the recipient to choose which parameter
				// he would like to decide the winning contract
				// 1. Price
				// 2. Completion time
				// 3. Select Manually
				fmt.Println("CHOOSE THE METRIC BY WHICH YOU WANT TO SELECT THE WINNING BID: ")
				fmt.Println("1. PRICE")
				fmt.Println("2. COMPLETION TIME (IN YEARS)")
				fmt.Println("3. SELECT MANUALLY")
				fmt.Println("ENTER YOUR CHOICE AS A NUMBER (1 / 2 / 3)")
				opt, err := utils.ScanForInt()
				if err != nil {
					fmt.Println("Couldn't read user input")
					break
				}
				contractorsForBid, contractsForBid, err := database.RetrieveAllProposedContracts(4)
				// retrieve all contracts towards the order
				if err != nil {
					log.Fatal(err)
				}
				// the length of the above two slices must be the same
				for i, contractor := range contractorsForBid {
					log.Println("======================================================================================")
					log.Println("Contractor Name: ", contractor.U.Name)
					log.Println("Proposed Contract: ")
					database.PrettyPrintProposedContract(contractsForBid[i].O)
				}
				switch opt {
				case 1:
					fmt.Println("YOU'VE CHOSEN TO SELECT BY LEAST PRICE")
					// here we assume that the timeout period for the auction is up and that
					// price is the winning metric of a specific bid, like in traditional contract
					bestContract, err := database.SelectContractByPrice(contractsForBid)
					if err != nil {
						log.Fatal(err)
					}
					log.Println("BEST CONTRACT IS: ")
					database.PrettyPrintProposedContract(bestContract.O)
					// now we need to replace the originator order with this order, order
					// indices are same, so we can insert
					err = database.InsertOrder(bestContract.O)
					if err != nil {
						log.Println(err)
						break
					}
					// now at this point, we need to mark this specific contract as completed.
					// do we set a flag? db entry? how do we do that
				case 2:
					fmt.Println("YOU'VE CHOSEN TO SELECT BY NUMBER OF YEARS")
					bestContract, err := database.SelectContractByTime(contractsForBid)
					if err != nil {
						log.Fatal(err)
					}
					log.Println("BEST CONTRACT IS: ")
					database.PrettyPrintProposedContract(bestContract.O)
					// finalize this order and open to investors
					err = database.InsertOrder(bestContract.O)
					if err != nil {
						log.Println(err)
						break
					}
				case 3:
					fmt.Println("HERE ARE A LIST OF AVAILABLE OPTIONS, CHOOSE ANY ONE CONTRACT")
					for i, contractor := range contractorsForBid {
						log.Println("======================================================================================")
						log.Println("CONTRACTOR NUMBER: ", i+1) // +1 to skip 0
						log.Println("Contractor Name: ", contractor.U.Name)
						log.Println("Proposed Contract: ")
						database.PrettyPrintProposedContract(contractsForBid[i].O)
					}
					fmt.Println("ENTER YOUR OPTION AS A NUMBER")
					opt, err := utils.ScanForInt()
					if err != nil {
						log.Println(err)
						break
					}
					// here we get the ContractEntity that the user has chosen
					// we need to insert the ContractEntity's Proposed Order
					err = database.InsertOrder(contractsForBid[opt-1].O) // -1  for order index
					if err != nil {
						log.Println(err)
						break
					}
					// now we need to set the winning order's chosen flag as true
					bestOrder, err := database.RetrieveOrder(opt)
					if err != nil {
						log.Println(err)
						break
					}
					// we also need to udpate the contract, once differences are clearer
					err = database.InsertOrder(bestOrder)
					if err != nil {
						log.Println(err)
						break
					}
					// TODO: have a flag filter
					log.Println("Marking order number: ", opt, " as completed")
				default:
					break
				}
				// now we need to finalize this person, potentially
				// move funds from the investor money to this person and so on.
				// another question is that whether we raise money before and then we have a
				// blind auction or whether we take in their feedback and then present this to
				// investors. Investors would ideally want to know more about what they are
				// investing in, so I guess the second option is better for now.
			default: // this default is for the larger switch case
				// check whether he wants to go back to the display all screen again
				fmt.Println("DO YOU REALLY WANT TO EXIT? (PRESS Y TO CONFIRM)")
				exitOpt, err := utils.ScanForString()
				if err != nil {
					log.Println(err)
					break
				}
				if exitOpt == "Y" || exitOpt == "y" {
					fmt.Println("YOU HAVE DECIDED TO EXIT")
					log.Fatal("")
				}
				break
			}
		}
		database.PrettyPrintRecipient(recipient)
		return
	}

	// User is an investor
	for {
		// Main investor loop
		fmt.Println("------------INVESTOR INTERFACE------------")
		fmt.Println("----CHOOSE ONE OF THE FOLLOWING OPTIONS----")
		fmt.Println("  1. Display all Open Orders")
		fmt.Println("  2. Display my Profile")
		fmt.Println("  3. Invest in an Order")
		fmt.Println("  4. Display All Balances")
		fmt.Println("  5. Exchange XLM for USD")
		fmt.Println("  6. Display all Originated Orders")
		fmt.Println("  7. Vote towards a specific proposed order")
		fmt.Println("  default: Exit")
		optI, err := utils.ScanForInt()
		if err != nil {
			fmt.Println("Couldn't read user input")
			break
		}
		switch optI {
		case 1:
			fmt.Println("------------LIST OF ALL AVAILABLE ORDERS------------")
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
			fmt.Println("----WHICH ORDER DO YOU WANT TO INVEST IN? (ENTER ORDER NUMBER WITHOUT SPACES)----")
			oNumber, err := utils.ScanForInt()
			if err != nil {
				fmt.Println("Couldn't read user input")
				break
			}
			// now the user has decided to invest in the asset with index uInput
			// we need to retrieve the order and ask for confirmation
			uOrder, err := database.RetrieveOrder(oNumber)
			if err != nil {
				log.Fatal("Order with specified index not found in the database")
			}
			database.PrettyPrintOrder(uOrder)
			fmt.Println(" HOW MUCH DO YOU WANT TO INVEST?")
			investmentAmount, err := utils.ScanForStringWithCheckI()
			if err != nil {
				log.Println(err)
				break
			}
			fmt.Println(" DO YOU WANT TO CONFIRM THIS ORDER? (PRESS N IF YOU DON'T WANT TO)")
			confirmOpt, err := utils.ScanForString()
			if err != nil {
				log.Println(err)
				break
			}
			if confirmOpt == "N" || confirmOpt == "n" {
				fmt.Println("YOU HAVE DECIDED TO CANCEL THIS ORDER")
				break
			}
			// when I am creating an account, I will have a PublicKey and Seed, so
			// don't need them here
			// check whether the investor has XLM already
			balance, err := xlm.GetXLMBalance(platformPublicKey)
			// balance is in string, convert to int
			balanceI := utils.StoF(balance)
			log.Println("Platform's balance is: ", balanceI)
			if balanceI < 21 { // 1 to account for fees
				// get coins if balance is this low
				log.Println("Refilling platform balance")
				err := xlm.GetXLM(platformPublicKey)
				// TODO: in future, need to refill platform sufficiently well and interact
				// with a cold wallet that we have previously set
				if err != nil {
					log.Fatal(err)
				}
			}

			balance, err = xlm.GetXLMBalance(platformPublicKey)
			log.Println("Platform balance updated is: ", balance)
			fmt.Printf("Platform seed is: %s and platform's publicKey is %s", platformSeed, platformPublicKey)
			log.Println("Investor's publickey is: ", investor.U.PublicKey)
			balance, err = xlm.GetXLMBalance(investor.U.PublicKey)
			if balance == "" {
				// means we need to setup an account first
				// Generating a keypair on stellar doesn't mean that you can send funds to it
				// you need to call the CreateAccount method in order to be able to send funds
				// to it
				log.Println("Investor balance empty, refilling!")
				_, _, err = xlm.SendXLMCreateAccount(investor.U.PublicKey, consts.DonateBalance, platformSeed)
				if err != nil {
					log.Println("Investor Account doesn't have funds")
					log.Fatal(err)
				}
			}
			// balance is in string, convert to float
			balance, err = xlm.GetXLMBalance(investor.U.PublicKey)
			balanceI = utils.StoF(balance)
			log.Println("Investor balance is: ", balanceI)
			if balanceI < 3 { // to setup trustlines
				_, _, err = xlm.SendXLM(investor.U.PublicKey, consts.DonateBalance, platformSeed)
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
				_, _, err = xlm.SendXLMCreateAccount(recipient.U.PublicKey, consts.DonateBalance, platformSeed)
				if err != nil {
					log.Println("Recipient Account doesn't have funds")
					log.Fatal(err)
				}
			}
			balance, err = xlm.GetXLMBalance(recipient.U.PublicKey)
			// balance is in string, convert to float
			balanceI = utils.StoF(balance)
			log.Println("Recipient balance is: ", balanceI)
			if balanceI < 3 { // to setup trustlines
				_, _, err = xlm.SendXLM(recipient.U.PublicKey, consts.DonateBalance, platformSeed)
				if err != nil {
					log.Println("Recipient Account doesn't have funds")
					log.Fatal(err)
				}
			}
			log.Println("The issuer's public key and private key are: ", platformPublicKey, " ", platformSeed)
			log.Println("The investor's public key and private key are: ", investor.U.PublicKey, " ", investor.U.Seed)
			log.Println("The recipient's public key and private key are: ", recipient.U.PublicKey, " ", recipient.U.Seed)

			log.Println(&investor, &recipient, investmentAmount, uOrder)
			// so now we have three entities setup, so we create the assets and invest in them
			cOrder, err := assets.InvestInOrder(platformPublicKey, platformSeed, &investor, &recipient, investmentAmount, uOrder) // assume payback period is 5
			if err != nil {
				log.Println(err)
				continue
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
			convAmount, err := utils.ScanForStringWithCheckF()
			if err != nil {
				log.Println(err)
				return
			}
			// maybe don't trust asset again when you've trusted it already? check if that's
			// possible and save on the tx fee for a single transaction. But I guess its
			// difficult to retrieve trustlines, so we'll go ahead with it
			hash, err := assets.TrustAsset(stablecoin.StableUSD, consts.StablecoinTrustLimit, investor.U.PublicKey, investor.U.Seed)
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
		case 6:
			// TODO: add voting scheme here
			fmt.Println("LIST OF ALL ORIGINATED ORDERS: ")
			originatedOrders, err := database.RetrieveAllOriginatedOrders()
			if err != nil {
				log.Println(err)
				break
			}
			database.PrettyPrintOrders(originatedOrders)
		case 7:
			// this is the case where an investor can vote on a particular proposed order
			fmt.Println("LIST OF ALL PROPOSED ORDERS: ")
			// now the investor can only put his funds in one potential proposed order.
			// his voting weight will be the amount of stableUSD that he has.
			// for testing, we assume that weights are 100 each.
			allContractors, err := database.RetrieveAllContractEntities("contractor")
			if err != nil {
				log.Println(err)
				break
			}
			for index, contractor := range allContractors {
				log.Println("CONTRACTOR NAME: ", contractor.U.Name)
				log.Println("CONTRACTOR INDEX: ", contractor.U.Index)
				for _, contract := range allContractors[index].ProposedContracts {
					database.PrettyPrintOrder(contract.O)
				}
			}
			// we need to get the vote of the investor here, but how do you get the vote?
			// because you have two arrays, you need to have some kind of a common index
			// between the two and then you can fetch the stuff back.
			// lets ask for the user index, which will be unique and the order number,
			// which can get us the specific proposed contract
			fmt.Println("WHICH INVESTOR DO YOU WANT TO VOTE TOWARDS?")
			vote, err := utils.ScanForInt()
			if err != nil {
				log.Println(err)
				break
			}
			log.Println("You have voted for user number: ", vote)
			fmt.Println("WHICH PROPOSED ORDER OF HIS DO YOU WANT TO VOTE TOWARDS?")
			pOrderN, err := utils.ScanForInt()
			if err != nil {
				log.Println(err)
				break
			}
			// we need to go through the contractor's proposed orders to find an order
			// with index pOrderN
			for _, elem := range allContractors {
				if elem.U.Index == vote {
					// we have the specific contractor
					log.Println("FOUND CONTRACTOR!")
					for i, pcs := range elem.ProposedContracts {
						if pcs.O.Index == pOrderN {
							// this is the order the user voted towards
							// need to update the vote stuff
							// check whether user has already voted
							// now an investor can vote  up to a max of his balance in StableUSD
							// the final call for selecting still falls on to the recipient, but
							// the recipient can get soem idea on which proposed contracts are
							// popular.
							fmt.Println("YOUR AVAILABLE VOTING BALANCE IS: ", investor.VotingBalance)
							fmt.Println("HOW MANY VOTES DO YOU WANT TO DELEGATE TOWARDS THIS ORDER?")
							votes, err := utils.ScanForInt()
							if err != nil {
								log.Println(err)
								break
							}
							if votes > investor.VotingBalance {
								log.Println("Can't vote with an amount greater than available balance")
								break
							}
							pcs.O.Votes += votes
							// an order's votes can exceed the total amount because it only shows
							// that many peopel feel that contract to be doing the right thing
							elem.ProposedContracts[i] = pcs
							// and we need to update the contractor, not order, since these
							// are not in the order database yet)
							err = database.InsertContractEntity(elem)
							if err != nil {
								log.Println(err)
								break
							}
							err = investor.DeductVotingBalance(votes)
							if err != nil {
								// fatal since we need to exit if there's a bug in updating votes
								log.Fatal(err)
								break
							}
							log.Println("CASTING VOTE! ", pcs)
						}
					}
				}
			}
		default:
			// check whether he wants to go back to the display all screen again
			fmt.Println("DO YOU REALLY WANT TO EXIT? (PRESS Y TO CONFIRM)")
			exitOpt, err := utils.ScanForString()
			if err != nil {
				log.Println(err)
				break
			}
			if exitOpt == "Y" || exitOpt == "y" {
				fmt.Println("YOU HAVE DECIDED TO EXIT")
				log.Fatal("")
			}
		} // end of switch
	}
	log.Fatal("")
	rpc.StartServer(opts.Port) // this must be towards the end
}
