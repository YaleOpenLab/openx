package main

// test.go runs the PoC stellar implementation calling various functions
import (
	"fmt"
	"log"
	"os"

	ipfs "github.com/Varunram/essentials/ipfs"
	consts "github.com/YaleOpenLab/openx/consts"
	database "github.com/YaleOpenLab/openx/database"
	// platform "github.com/YaleOpenLab/openx/platforms"
	stablecoin "github.com/Varunram/essentials/crypto/stablecoin"
	wallet "github.com/Varunram/essentials/crypto/xlm/wallet"
	scan "github.com/Varunram/essentials/scan"
	utils "github.com/Varunram/essentials/utils"
	solar "github.com/YaleOpenLab/openx/platforms/opensolar"
	rpc "github.com/YaleOpenLab/openx/rpc"
	// xlm "github.com/Varunram/essentials/crypto/xlm"
	flags "github.com/jessevdk/go-flags"
)

// test.go drives the CLI interface and is intended to be a CLI client that can be used
// to interact with the openx platform
// TODO: move to the teller based config system mimicking the frontend once we have RPCs
// for functions that will be used by the frontend.
var opts struct {
	Port int `short:"p" description:"The port on which the server runs on"`
}

// ParseConfig parses ofcli config
func ParseConfig(args []string) error {
	_, err := flags.ParseArgs(&opts, args)
	if err != nil {
		return err
	}
	port := utils.ItoS(consts.DefaultRpcPort)
	if opts.Port != 0 {
		port = utils.ItoS(opts.Port)
	}
	log.Println("Starting RPC Server on Port: ", opts.Port)
	go rpc.StartServer(port, false) // run as go routine for now
	return nil
}

func main() {
	err := ParseConfig(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	var investorSeed string
	var recipientSeed string

	consts.PlatformPublicKey, consts.PlatformSeed, err = StartPlatform()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("PLATFORM SEED IS: %s\n PLATFORM PUBLIC KEY IS: %s\n", consts.PlatformSeed, consts.PlatformPublicKey)
	// TODO: how much do we pay the investor?
	// Do we sell the REC created from the solar panels only to the investor? If so,
	// isn't that enough to propel investment in the solar contract itself?
	// TODO: need a server to run a public stellar node to test out stuff
	// change the API mapping
	// move current number of years metric to a separate package since that is
	// more suitable for a model like affordable housing.
	// look into what kind of data we get from the pi and checkout pi specific code
	// to see if we can get something from there.
	// TODO: Need to automatically cover breach scenarios in case the recipient doesn't
	// pay for a specific period of time
	// TODO: also write a Makefile so that its easy for people to get started with stuff
	// TODO: move most of this stuff to the new emulator struct
	fmt.Println("------------STELLAR HOUSE INVESTMENT CLI INTERFACE (RETIRED, USE EMULATOR)------------")

	// init stablecoin stuff
	err = stablecoin.InitStableCoin()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("------WHAT DO YOU WANT TO DO?------")
	fmt.Println("1. CREATE A NEW INVESTOR ACCOUNT")
	fmt.Println("2. CREATE A NEW RECIPIENT ACCOUNT")
	fmt.Println("3: ALREADY HAVE AN ACCOUNT")
	opt, err := scan.ScanForInt()
	if err != nil {
		fmt.Println("Couldn't read user input")
		return
	}
	switch opt {
	case 1:
		err := NewInvestorPrompt()
		if err != nil {
			log.Println(err)
			return
		}
	case 2:
		err := NewRecipientPrompt()
		if err != nil {
			log.Println(err)
			return
		}
	default:
	}

	investor, recipient, contractor, isRecipient, isContractor, err := LoginPrompt()
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
			fmt.Println("  0. Display all Locked Projects")               // done
			fmt.Println("  1. Display all Open Projects (STAGE 3)")       // done
			fmt.Println("  2. Display my Profile")                        // done
			fmt.Println("  3. Payback towards an Project (STAGE 6)")      // done
			fmt.Println("  4. Exchange XLM for USD")                      // done
			fmt.Println("  5. Finalize a specific Project (STAGE 2->3)")  // done
			fmt.Println("  6. Originate a specific Project (STAGE 0->1)") // done
			fmt.Println("  7. View all Projects (ALL STAGES)")            // done
			fmt.Println("  8. View all Origin Projects (STAGE 1)")        // done
			fmt.Println("  9. View All Balances ")                        // done
			fmt.Println("  10. Unlock Account")                           // done
			fmt.Println("  default: Exit")
			optI, err := scan.ScanForInt()
			if err != nil {
				fmt.Println("Couldn't read user input")
				break
			}
			switch optI {
			case 0:
				fmt.Println("CHOOSE A PROJECT TO UNLOCK")
				allProjects, err := solar.RetrieveLockedProjects()
				if err != nil {
					log.Println(err)
					break
				}
				log.Println(allProjects)
				pIndex, err := scan.ScanForInt()
				if err != nil {
					log.Println(err)
					break
				}

				fmt.Println("ENTER SEED PASSWORD:")
				// need to unlock the recipient account
				seedpwd, err := scan.ScanRawPassword()
				if err != nil {
					log.Println(err)
					break
				}
				// make sure that the seed provided is valid
				_, err = wallet.DecryptSeed(recipient.U.StellarWallet.EncryptedSeed, seedpwd)
				if err != nil {
					log.Println(err)
					break
				}
				// unlock the project
				err = solar.UnlockProject(recipient.U.Username, recipient.U.Pwhash, pIndex, seedpwd)
				if err != nil {
					log.Println(err)
				}
			case 1:
				// getFinalProjects RPC
				Stage3ProjectsDisplayPrompt()
			case 2:
				// validateRecipient RPC
				// retrieve again to get changes that may have occurred in between
				recipient, err = database.RetrieveRecipient(recipient.U.Index)
				if err != nil {
					log.Println(err)
					break
				}
				PrintRecipient(recipient)
			case 3:
				// payback RPC
				log.Println(recipient.ReceivedSolarProjects)
				fmt.Println("WHICH PROJECT DO YOU WANT TO PAY BACK TOWARDS? (ENTER PROJECT NUMBER)")
				projectNumber, err := scan.ScanForInt()
				if err != nil {
					log.Println("INPUT NOT AN INTEGER, TRY AGAIN")
					continue
				}
				// check if we can get the project using the project number that we have here
				rtContract, err := solar.RetrieveProject(projectNumber)
				if err != nil {
					log.Println("Couldn't retrieve project, try again!")
					continue
				}
				// so we can retrieve the project using the project Index, nice
				PrintProject(rtContract)
				fmt.Println("HOW MUCH DO YOU WANT TO PAYBACK?")
				paybackAmount, err := scan.ScanForStringWithCheckI()
				if err != nil {
					log.Println(err)
					break
				}
				fmt.Printf(" DO YOU WANT TO CONFIRM THAT YOU WANT TO PAYBACK %s TOWARDS THIS PROJECT? (PRESS N IF YOU DON'T WANT TO)\n", paybackAmount)
				confirmOpt, err := scan.ScanForString()
				if err != nil {
					log.Println(err)
					break
				}
				if confirmOpt == "N" || confirmOpt == "n" {
					fmt.Println("YOU HAVE DECIDED TO CANCEL THE PAYBACK ORDER")
					break
				}
				fmt.Printf("PAYING BACK %s TOWARDS PROJECT NUMBER: %d\n", paybackAmount, rtContract.Index) // use the rtContract.Params here instead of using projectNumber from long ago

				err = solar.Payback(recipient.U.Index, rtContract.Index, rtContract.DebtAssetCode, paybackAmount, recipientSeed)
				if err != nil {
					log.Println("PAYBACK TX FAILED, PLEASE TRY AGAIN!", err)
				}
			case 4:
				// getStableCoin RPC
				log.Println("Enter the amount you want to convert into STABLEUSD")
				convAmount, err := scan.ScanForStringWithCheckF()
				if err != nil {
					log.Println(err)
					break
				}
				err = stablecoin.Exchange(recipient.U.StellarWallet.PublicKey, recipientSeed, convAmount)
				if err != nil {
					log.Println(err)
				}
			case 5:
				var bestContract solar.Project
				var err error
				allContracts, err := solar.RetrieveRecipientProjects(solar.Stage2.Number, recipient.U.Index)
				if err != nil {
					log.Println(err)
					continue
				}
				PrintProjects(allContracts)
				// TODO: port this to the emulator
				fmt.Println("CHOOSE THE METRIC BY WHICH YOU WANT TO SELECT THE WINNING BID: ")
				fmt.Println("1. PRICE (BLIND)")
				fmt.Println("2. COMPLETION TIME (IN YEARS)")
				fmt.Println("3. SELECT MANUALLY")
				fmt.Println("4. PRICE (VICKREY)")
				fmt.Println("ENTER YOUR CHOICE AS A NUMBER (1 / 2 / 3)")
				opt, err := scan.ScanForInt()
				if err != nil {
					fmt.Println("Couldn't read user input")
					break
				}
				// the length of the above two slices must be the same
				for _, elem := range allContracts {
					log.Println("======================================================================================")
					PrintProject(elem)
				}
				switch opt {
				case 1:
					fmt.Println("YOU'VE CHOSEN TO SELECT BY BLIND AUCTION RULES")
					bestContract, err = solar.SelectContractBlind(allContracts)
					if err != nil {
						log.Println(err)
						continue
					}
				case 2:
					fmt.Println("YOU'VE CHOSEN TO SELECT BY NUMBER OF YEARS")
					bestContract, err = solar.SelectContractTime(allContracts)
					if err != nil {
						log.Println(err)
						continue
					}
				case 3:
					for i, contract := range allContracts {
						log.Println("BEST BID CHOICE NUMBER: ", i)
						PrintProject(contract)
					}
					fmt.Println("ENTER YOUR OPTION AS A NUMBER")
					opt, err := scan.ScanForInt()
					if err != nil {
						log.Println(err)
						break
					}
					bestContract = allContracts[opt]
				case 4:
					fmt.Println("YOU'VE CHOSEN TO SELECT BY VICKREY AUCTION RULES")
					// here we assume that the timeout period for the auction is up and that
					// price is the winning metric of a specific bid, like in traditional contract
					bestContract, err = solar.SelectContractVickrey(allContracts)
					if err != nil {
						log.Println(err)
						continue
					}
				default:
					break
				}
				err = bestContract.SetStage(3)
				if err != nil {
					log.Println(err)
					break
				}
				log.Println("BEST CONTRACT IS: ")
				PrintProject(bestContract)
				// now the contract is at stage 3
			case 6:
				fmt.Println("LIST OF ALL PRE ORIGIN PROJECTS BY ORIGINATORS (STAGE 0)")
				allMyProjects, err := solar.RetrieveProjectsAtStage(solar.Stage0.Number)
				if err != nil {
					log.Println(err)
					continue
				}
				PrintProjects(allMyProjects)
				fmt.Println("ENTER THE PROJECT INDEX")
				projectIndex, err := scan.ScanForInt()
				if err != nil {
					log.Println(err)
					continue
				}
				err = solar.RecipientAuthorize(projectIndex, recipient.U.Index)
				if err != nil {
					log.Println(err)
				}
			case 7:
				fmt.Println("PRINTING ALL PROJECTS: ")
				allContracts, err := solar.RetrieveAllProjects()
				if err != nil {
					log.Println(err)
					break
				}
				PrintProjects(allContracts)
			case 8:
				DisplayOriginProjects()
			case 9:
				BalanceDisplayPrompt(recipient.U.StellarWallet.PublicKey)
			case 10:
				// need to unlock the recipient account
				seedpwd, err := scan.ScanRawPassword()
				if err != nil {
					log.Println(err)
					break
				}
				seed, err := wallet.DecryptSeed(recipient.U.StellarWallet.EncryptedSeed, seedpwd)
				if err != nil {
					log.Println(err)
					break
				}
				recipientSeed = seed
				log.Println(" Seed successfully unlocked")
			default: // this default is for the larger switch case
				ExitPrompt()
			}
		}
		return
	} else if isContractor {
		var optI int
		if contractor.Contractor {
			optI = 1
		} else if contractor.Originator {
			optI = 2
		}
		switch optI {
		case 1:
			fmt.Println("-------------WELCOME BACK CONTRACTOR-------------")
			for {
				fmt.Println("WHAT WOULD YOU LIKE TO DO?")
				fmt.Println("  1. VIEW ALL ORIGINATED (STAGE 1) PROJECTS")  // done
				fmt.Println("  2. VIEW PROFILE")                            // done
				fmt.Println("  3. CREATE A PROPOSED (STAGE 2) PROJECT")     // added to TODO
				fmt.Println("  4. VIEW ALL MY PROPOSED (STAGE 2) PROJECTS") // done
				fmt.Println("  5. CREATE NEW TYPE OF COLLATERAL")           // done
				optI, err := scan.ScanForInt()
				if err != nil {
					log.Println(err)
					continue
				}
				switch optI {
				case 1:
					DisplayOriginProjects()
				case 2:
					PrintEntity(contractor)
				case 3:
					fmt.Println("YOU HAVE CHOSEN TO CREATE A NEW PROPOSED PROJECT")
					err = ProposeContractPrompt(&contractor)
					if err != nil {
						log.Println(err)
						continue
					}
				case 4:
					fmt.Println("LIST OF ALL PROPOSED CONTRACTS BY ME: ")
					allMyProjects, err := solar.RetrieveContractorProjects(solar.Stage2.Number, contractor.U.Index)
					if err != nil {
						log.Println(err)
						continue
					}
					PrintProjects(allMyProjects)
				case 5:
					fmt.Println("YOU HAVE CHOSEN TO ADD A NEW TYPE OF COLLATERAL")
					fmt.Println("Enter collateral amount")
					colAmount, err := scan.ScanForFloat()
					if err != nil {
						fmt.Println(err)
						continue
					}
					fmt.Println("Enter collateral data")
					colData, err := scan.ScanForString()
					if err != nil {
						fmt.Println(err)
						continue
					}
					err = contractor.AddCollateral(colAmount, colData)
					if err != nil {
						fmt.Println(err)
						continue
					}
					fmt.Println("Please upload documents for verification with KYC inspector")
				}
			}
		case 2:
			fmt.Println("-------------WELCOME BACK ORIGINATOR-------------")
			for {
				fmt.Println("WHAT WOULD YOU LIKE TO DO?")
				fmt.Println("  1. PROPOSE A PRE-ORIGIN (STAGE 0) PROJECT TO A RECIPIENT") // added to TODO
				fmt.Println("  2. VIEW PROFILE")                                          // done
				fmt.Println("  3. VIEW ALL MY PRE-ORIGINATED (STAGE 0) PROJECTS")         // done
				fmt.Println("  4. VIEW ALL MY ORIGINATED (STAGE 1) PROJECTS")             // done
				optI, err = scan.ScanForInt()
				if err != nil {
					log.Println(err)
					continue
				}
				switch optI {
				case 1:
					err := OriginContractPrompt(&contractor)
					if err != nil {
						fmt.Println("RETURNING BACK TO THE MAIN LOOP: ", err)
						continue
					}
				case 2:
					PrintEntity(contractor)
				case 3:
					allMyProjects, err := solar.RetrieveOriginatorProjects(solar.Stage0.Number, contractor.U.Index)
					if err != nil {
						fmt.Println("RETURNING BACK TO THE MAIN LOOP: ", err)
						continue
					}
					PrintProjects(allMyProjects)
				case 4:
					allMyProjects, err := solar.RetrieveOriginatorProjects(solar.Stage1.Number, contractor.U.Index)
					if err != nil {
						fmt.Println("RETURNING BACK TO THE MAIN LOOP: ", err)
						continue
					}
					PrintProjects(allMyProjects)
				default:
					ExitPrompt()
				}
			}
		}
		PrintEntity(contractor)
	} else {
		// User is an investor
		for {
			// Main investor loop
			fmt.Println("------------INVESTOR INTERFACE------------") //  done
			fmt.Println("----CHOOSE ONE OF THE FOLLOWING OPTIONS----")
			fmt.Println("  1. Display all Open Projects (STAGE 3)") // done
			fmt.Println("  2. Display my Profile")                  // done
			// fmt.Println("  3. Invest in an Project (STAGE 3)") // done
			fmt.Println("  4. Display All Balances")                               // done
			fmt.Println("  5. Exchange XLM for USD")                               // done
			fmt.Println("  6. Display all Origin (STAGE 1) Projects")              // done
			fmt.Println("  7. Vote towards a specific proposed project (STAGE 2)") // done
			fmt.Println("  8. Get ipfs hash of a contract")                        // done
			fmt.Println("  9. Display all Funded Projects")                        // done
			fmt.Println("  10. Unlock account")                                    // done
			fmt.Println("  11. KYC users (admin only)")                            // done
			fmt.Println("  default: Exit")
			optI, err := scan.ScanForInt()
			if err != nil {
				fmt.Println("Couldn't read user input")
				break
			}
			switch optI {
			case 1:
				Stage3ProjectsDisplayPrompt()
			case 2:
				PrintInvestor(investor)
			case 3:
				// investInProject RPC
				// This function has been removed from the CLI since once you invest in a particular order and it reaches
				// the limit, this function will not transfer assets back to the recipient, resulting
				// in an improper way of emulating the workflow. The only option is to call the route,
				// which will be called by the frontend, so we can emulate this successfully.
				// curl -X GET -H "Content-Type: application/x-www-form-urlencoded" -H "Origin: localhost" -H "Cache-Control: no-cache" "http://localhost:8080/investor/invest?username=john&password=9a768ace36ff3d1771d5c145a544de3d68343b2e76093cb7b2a8ea89ac7f1a20c852e6fc1d71275b43abffefac381c5b906f55c3bcff4225353d02f1d3498758&seedpwd=x&projIndex=1&amount=14000"
			case 4:
				BalanceDisplayPrompt(investor.U.StellarWallet.PublicKey)
			case 5:
				// Stablecoin/get route
				// curl -X GET -H "Content-Type: application/x-www-form-urlencoded" -H "Origin: localhost" -H "Cache-Control: no-cache" "http://localhost:8080/stablecoin/get?seed=SB2Z5GZASNF4ZR7263WWYISZP3UXSP7A6IP6ENZ44G4T44G6NVUCSVSP&amount=1"
				log.Println("Enter the amount you want to convert into STABLEUSD")
				convAmount, err := scan.ScanForStringWithCheckF()
				if err != nil {
					log.Println(err)
					break
				}
				err = stablecoin.Exchange(investor.U.StellarWallet.PublicKey, investorSeed, convAmount)
				if err != nil {
					log.Println(err)
				}
			case 6:
				DisplayOriginProjects()
			case 7:
				// this is the case where an investor can vote on a particular proposed project
				// ie stage 2 projects for the recipient to have an understanding about
				// which contracts are popular and can receive more investor money
				fmt.Println("LIST OF ALL PROPOSED ORDERS: ")
				allProposedProjects, err := solar.RetrieveProjectsAtStage(solar.Stage2.Number)
				if err != nil {
					log.Println(err)
					break
				}
				PrintProjects(allProposedProjects)
				fmt.Println("WHICH CONTRACT DO YOU WANT TO VOTE TOWARDS?")
				vote, err := scan.ScanForInt()
				if err != nil {
					log.Println(err)
					break
				}
				fmt.Println("YOUR AVAILABLE VOTING BALANCE IS: ", investor.VotingBalance)
				fmt.Println("HOW MANY VOTES DO YOU WANT TO DELEGATE TOWARDS THIS ORDER?")
				votes, err := scan.ScanForFloat()
				if err != nil {
					log.Println(err)
					break
				}
				err = solar.VoteTowardsProposedProject(investor.U.Index, votes, vote)
				if err != nil {
					log.Println(err)
					break
				}
				log.Println("You have voted for contract number: ", vote)
			case 8:
				fmt.Println("WELCOME TO THE IPFS HASHING INTERFACE")
				fmt.Println("ENTER THE STRING THAT YOU WOULD LIKE THE IPFS HASH FOR")
				// the UI should ideally have a menu that asks the user for a file and then
				// produces the hash of it. In this case, we shall use a sample ipfs file
				// and then hash it.
				// this uses the platform's ipfs key though, not the user's. If the user
				// wants to serve his own ipfs files, he is better off running a client on
				// his own
				hashString, err := scan.ScanForString()
				if err != nil {
					fmt.Println("Couldn't read user input, going back to the main loop")
				}
				hash, err := ipfs.AddStringToIpfs(hashString)
				if err != nil {
					fmt.Println("Couldn't hash user input, exiting to main menu", err)
					break
				}
				hashCheck, err := ipfs.GetStringFromIpfs(hash)
				if err != nil || hashCheck != hashString {
					fmt.Println("Hashed strings and retrieved strings don't match, don't use this hash!", err)
					break
				}
				// don't print this hash unless we can decrypt it and be sure that it behaves as expected
				log.Println("THE HASH OF THE PROVIDED STRING IS: ", hash)
				// try to retrieve the string back from ipfs and check if it works correctly
			case 9:
				fmt.Println("LIST OF ALL FUNDED PROJECTS: ")
				allFundedProjects, err := solar.RetrieveProjectsAtStage(solar.Stage3.Number)
				if err != nil {
					fmt.Println(err)
					break
				}
				PrintProjects(allFundedProjects)
			case 10:
				// need to unlock the recipient account
				seedpwd, err := scan.ScanRawPassword()
				if err != nil {
					log.Println(err)
					break
				}
				seed, err := wallet.DecryptSeed(investor.U.StellarWallet.EncryptedSeed, seedpwd)
				if err != nil {
					log.Println(err)
					break
				}
				investorSeed = seed
				log.Println(" Seed successfully unlocked: ", seed)
			case 11:
				if !investor.U.Inspector {
					fmt.Println("You do not have access to this page")
					break
				}
				fmt.Println("WELCOME TO THE KYC INTERFACE!!")
				fmt.Println("CHOOSE AN OPTION FROM THE FOLLOWING MENU")
				fmt.Println("1. VIEW ALL KYC'D USERS")
				fmt.Println("2. VIEW ALL NON KYC'D USERS")
				sInput, err := scan.ScanForInt()
				if err != nil {
					log.Println(err)
					break
				}
				switch sInput {
				case 1:
					allUsers, err := database.RetrieveAllUsersWithKyc()
					if err != nil {
						log.Println(err)
						break
					}
					PrintUsers(allUsers)
				case 2:
					allUsers, err := database.RetrieveAllUsersWithoutKyc()
					if err != nil {
						log.Println(err)
						break
					}
					PrintUsers(allUsers)
					fmt.Println("WHICH USER DO YOU WANT TO AUTHENTICATE WITH KYC?")
					uInput, err := scan.ScanForInt()
					if err != nil {
						log.Println(err)
						break
					}
					err = investor.U.Authorize(uInput)
					if err != nil {
						log.Println(err)
						break
					}
				default:
					log.Println("Invalid input, please enter valid input")
				}
			default:
				ExitPrompt()
			} // end of switch
		}
		// it should never arrive here
		return
	}
}
