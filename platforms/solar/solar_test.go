// +build all travis

package solar

import (
	"log"
	"os"
	"testing"

	database "github.com/OpenFinancing/openfinancing/database"
)

// TODO: rewrite how this works and split between platforms and database
// go test --tags="all" -coverprofile=test.txt .
func TestDb(t *testing.T) {
	var err error
	os.Remove(os.Getenv("HOME") + "/.openfinancing/database/" + "/yol.db")
	err = os.MkdirAll(os.Getenv("HOME")+"/.openfinancing/database", os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	db, err := database.OpenDB()
	if err != nil {
		t.Fatal(err)
	}
	db.Close() // close immmediately after check
	var project1 SolarParams
	var dummy SolarProject
	// investors entity testing over, test recipients below in the same way
	// now we repeat the same tests for all other entities
	// connections and the other for non RPC connections
	// now we repeat the same tests for all other entities
	// connections and the other for non RPC connections
	inv, err := database.NewInvestor("inv1", "blah", "blah", "cool")
	if err != nil {
		t.Fatal(err)
	}
	err = inv.Save()
	if err != nil {
		t.Fatal(err)
	}
	log.Println("INDEX IS: ", inv.U.Index)
	recp, err := database.NewRecipient("user1", "blah", "blah", "cool")
	if err != nil {
		t.Fatal(err)
	}
	err = recp.Save()
	if err != nil {
		t.Fatal(err)
	}

	newCE2, err := NewOriginator("OrigTest2", "pwd", "blah", "NameOrigTest2", "123 ABC Street", "OrigDescription2")
	if err != nil {
		t.Fatal(err)
	}
	// tests for contractors
	newCE, err := NewContractor("ConTest", "pwd", "blah", "NameConTest", "123 ABC Street", "ConDescription") // use and test this as well
	if err != nil {
		t.Fatal(err)
	}
	project1.Index = 1
	project1.PanelSize = "100 1000 sq.ft homes each with their own private spaces for luxury"
	project1.TotalValue = 14000
	project1.Location = "India Basin, San Francisco"
	project1.MoneyRaised = 0
	project1.Metadata = "India Basin is an upcoming creative project based in San Francisco that seeks to invite innovators from all around to participate"
	project1.INVAssetCode = ""
	project1.DEBAssetCode = ""
	project1.PBAssetCode = ""
	project1.DateInitiated = ""
	project1.Years = 3
	project1.ProjectRecipient = recp
	dummy.Params = project1
	dummy.Contractor = newCE
	dummy.Originator = newCE2
	dummy.Stage = 3
	err = dummy.Save()
	if err != nil {
		t.Errorf("Inserting an project into the database failed")
		// shouldn't really fatal here, but this is in main, so we can't return
	}
	project, err := RetrieveProject(dummy.Params.Index)
	if err != nil {
		log.Println(err)
		t.Errorf("Retrieving project from the database failed")
		// again, shouldn't really fat a here, but we're in main
	}
	klx1, err := RetrieveProject(1000)
	if klx1.Params.Index != 0 {
		t.Fatalf("REtrieved project which does not exist, quitting!")
	}
	log.Println("Retrieved project: ", project)
	if project.Params.Index != dummy.Params.Index {
		t.Fatalf("Indices don't match, quitting!")
	}
	dummy.Params.Index = 2 // change index and try inserting another project
	err = dummy.Save()
	if err != nil {
		log.Println(err)
		t.Errorf("Inserting an project into the database failed")
		// shouldn't really fatal here, but this is in main, so we can't return
	}
	projects, err := RetrieveAllProjects()
	if err != nil {
		log.Println("Retrieve all error: ", err)
		t.Errorf("Failed in retrieving all projects")
	}
	if projects[0].Params.Index != 1 {
		t.Fatalf("Index of first element doesn't match, quitting!")
	}
	log.Println("Retrieved projects: ", projects)
	oProjects, err := RetrieveProjects(OriginProject)
	if err != nil {
		log.Println("Retrieve all error: ", err)
		t.Errorf("Failed in retrieving all originated projects")
	}
	if len(oProjects) != 0 {
		log.Println("OPROJECTS: ", len(oProjects))
		t.Fatalf("Originated projects present!")
	}
	err = database.DeleteKeyFromBucket(dummy.Params.Index, database.ProjectsBucket)
	if err != nil {
		log.Println(err)
		t.Errorf("Deleting an order from the db failed")
	}
	log.Println("Deleted project")
	// err = DeleteProject(1000) this would work becuase the key will not be found and hence would not return an error
	pp1, err := RetrieveProject(dummy.Params.Index)
	if err == nil && pp1.Params.Index != 0 {
		log.Println(err)
		// this should fail because we're trying to read an empty key value pair
		t.Errorf("Found deleted entry, quitting!")
	}

	inv, err = database.NewInvestor("inv2", "b921f75437050f0f7d2caba6303d165309614d524e3d7e6bccf313f39d113468d30e1e2ac01f91f6c9b66c083d393f49b3177345311849edb026bb86ee624be0", "blah", "cool")
	if err != nil {
		t.Fatal(err)
	}
	err = inv.Save()
	if err != nil {
		t.Fatal(err)
	}
	recp, err = database.NewRecipient("user1", "b921f75437050f0f7d2caba6303d165309614d524e3d7e6bccf313f39d113468d30e1e2ac01f91f6c9b66c083d393f49b3177345311849edb026bb86ee624be0", "blah", "cool")
	if err != nil {
		t.Fatal(err)
	}
	err = recp.Save()
	if err != nil {
		t.Fatal(err)
	}
	// tests for originators
	newCE, err = NewEntity("OrigTest", "pwd", "blah", "NameOrigTest", "123 ABC Street", "OrigDescription", "originator")
	if err != nil {
		t.Fatal(err)
	}
	allOrigs, err := RetrieveAllContractEntities("originator")
	if err != nil {
		t.Fatal(err)
	}
	acz1, err := RetrieveAllContractEntities("random")
	if len(acz1) != 0 {
		log.Println(acz1)
		t.Fatalf("Category which does not exist exists?")
	}
	if len(allOrigs) != 2 || allOrigs[0].U.Name != "NameOrigTest2" {
		t.Fatal("Names don't match, quitting!")
	}

	allConts, err := RetrieveAllContractEntities("contractor")
	if err != nil {
		t.Fatal(err)
	}
	if len(allConts) != 1 || allConts[0].U.Name != "NameConTest" {
		t.Fatal("Names don't match, quitting!")
	}
	_, err = newCE.ProposeContract("100 16x32 panels", 28000, "Puerto Rico", 6, "LEED+ Gold rated panels and this is random data out of nowhere and we supply our own devs and provide insurance guarantee as well. Dual audit maintenance upto 1 year. Returns capped as per defaults", recp.U.Index, 1)
	// 1 for retrieving martin as the recipient and 1 is the project Index
	if err != nil {
		log.Fatal(err)
	}
	_, err = newCE.ProposeContract("100 16x32 panels", 28000, "Puerto Rico", 6, "LEED+ Gold rated panels and this is random data out of nowhere and we supply our own devs and provide insurance guarantee as well. Dual audit maintenance upto 1 year. Returns capped as per defaults", 1000, 1)
	// 1 for retrieving martin as the recipient and 1 is the project Index
	if err == nil {
		t.Fatal("Able to retrieve non existent recipient, quitting!")
	}
	rOx, err := RetrieveProject(2)
	if err != nil {
		t.Fatal(err)
	}
	rOx.Params.ProjectRecipient = recp
	err = rOx.Save()
	if err != nil {
		t.Fatal(err)
	}

	allPCs, err := RetrieveProjectsR(ProposedProject, 6)
	if err != nil {
		t.Fatal(err)
	}
	if len(allPCs) != 1 { // add check for stuff here
		log.Println("LEN all proposed projects", len(allPCs))
	}
	rPC, err := FindInKey(2, allPCs)
	if err != nil {
		t.Fatal(err)
	}
	if rPC.Params.Index != rOx.Params.Index {
		t.Fatal("Indices don't match")
	}

	// now come the failure cases which should fail and we shall catch the case when they don't
	allPCs, err = RetrieveProjectsC(ProposedProject, 2)
	if len(allPCs) != 0 {
		log.Println("LEBNGRG: ", len(allPCs))
		t.Fatalf("Retrieving a missing contract succeeds, quitting!")
	}
	rPC, err = FindInKey(2, allPCs)
	if err == nil {
		t.Fatalf("Entity which should be missing exists!")
	}

	// newCE, err = NewEntity("ConTest", "pwd", "NameConTest", "123 ABC Street", "ConDescription", "contractor")
	rCE, err := SearchForEntity("ConTest", "9f88a8d40b90616715f868ed195d24e5df994f56bce34eddb022c213484eb0f220d8907e4ecd8f64ddd364cb30bb5758b32ee26541f340b930f7e5bf756907a4")
	if err != nil {
		t.Fatal(err)
	}
	if rCE.Contractor != true {
		log.Println("The role: ", rCE.Contractor, rCE)
		t.Fatal("Roles don't match!")
	}
	trC1, err := RetrieveEntity(7)
	if err != nil || trC1.U.Index == 0 {
		t.Fatal("SolarProject Entity lookup failed")
	}
	tmpx1, err := newCE2.OriginContract("100 16x24 panels on a solar rooftop", 14000, "Puerto Rico", 5, "ABC School in XYZ peninsula", recp.U.Index) // 1 is the index for martin
	if err != nil {
		log.Fatal(err)
	}
	_, err = newCE2.OriginContract("100 16x24 panels on a solar rooftop", 14000, "Puerto Rico", 5, "ABC School in XYZ peninsula", 1000) // 1 is the index for martin
	if err == nil {
		t.Fatalf("Not quitting for invalid recipient index")
	}
	tmpx1.Stage = 1
	err = tmpx1.Save()
	if err != nil {
		t.Fatal(err)
	}
	allOOs, err := RetrieveProjects(OriginProject) // this checks for stage 1 and not zero like the thing installed above
	if err != nil {
		t.Fatal(err)
	}
	if len(allOOs) != 1 {
		log.Println("Length of all Originated Projects: ", len(allOOs))
		t.Fatalf("Length of all orignated projects doesn't match")
	}
	// can't test the payback stuff since we need the recipient to have funds and stuff
	// maybe test the investor sending funds and stuff as well? but that's imported
	// by main and we could have problems here related to that
	project1.Index = 20
	nOP, err := NewOriginProject(project1, newCE2)
	if err != nil {
		t.Fatal(err)
	}
	err = PromoteStage0To1Project(20)
	if err == nil {
		t.Fatalf("Able to promote project which doesn't exist!")
	}
	err = nOP.SetPreOriginProject()
	if err != nil {
		t.Fatal(err)
	}
	if nOP.Stage != 0 {
		t.Fatalf("Stage doesn't match, quitting!")
	}

	err = nOP.SetLegalContractStage()
	if err != nil {
		t.Fatal(err)
	}
	if nOP.Stage != 0.5 {
		t.Fatalf("Stage doesn't match, quitting!")
	}

	err = nOP.SetOriginProject()
	if err != nil {
		t.Fatal(err)
	}
	if nOP.Stage != 1 {
		t.Fatalf("Stage doesn't match, quitting!")
	}

	err = nOP.SetOpenForMoneyStage()
	if err != nil {
		t.Fatal(err)
	}
	if nOP.Stage != 1.5 {
		t.Fatalf("Stage doesn't match, quitting!")
	}

	err = nOP.SetProposedProject()
	if err != nil {
		t.Fatal(err)
	}
	if nOP.Stage != 2 {
		t.Fatalf("Stage doesn't match, quitting!")
	}

	err = nOP.SetFinalizedProject()
	if err != nil {
		t.Fatal(err)
	}
	if nOP.Stage != 3 {
		t.Fatalf("Stage doesn't match, quitting!")
	}

	err = nOP.SetFundedProject()
	if err != nil {
		t.Fatal(err)
	}
	if nOP.Stage != 4 {
		t.Fatalf("Stage doesn't match, quitting!")
	}

	err = nOP.SetInstalledProjectStage()
	if err != nil {
		t.Fatal(err)
	}
	if nOP.Stage != 5 {
		t.Fatalf("Stage doesn't match, quitting!")
	}

	err = nOP.SetPowerGenerationStage()
	if err != nil {
		t.Fatal(err)
	}
	if nOP.Stage != 6 {
		t.Fatalf("Stage doesn't match, quitting!")
	}

	// cycle back to stage 0 and try using the other function to modify the stage
	err = nOP.SetPreOriginProject()
	if err != nil {
		t.Fatal(err)
	}
	if nOP.Stage != 0 {
		t.Fatalf("Stage doesn't match, quitting!")
	}

	err = PromoteStage0To1Project(20)
	if err == nil {
		t.Fatalf("Able to promote project which doesn't exist!")
	}
	allO, err := RetrieveProjectsO(OriginProject, newCE2.U.Index)
	if err != nil {
		t.Fatal(err)
	}
	if len(allO) != 1 {
		t.Fatalf("Multiple originated orders when there should be only one")
	}
	allProposedProjects, err := RetrieveProjects(ProposedProject)
	if err != nil {
		t.Fatal(err)
	}
	err = inv.AddVotingBalance(1000)
	if err != nil {
		t.Fatal(err)
	}
	err = nOP.SetProposedProject()
	if err != nil {
		t.Fatal(err)
	}
	for _, elem := range allProposedProjects {
		log.Println("SETHIS: ", elem.Params.Index)
	}
	err = VoteTowardsProposedProject(&inv, 100, 2)
	if err != nil {
		t.Fatal(err)
	}
	// func UpdateProjectSlice(a *database.Recipient, project SolarParams) error {
	recp.ReceivedSolarProjects = append(recp.ReceivedSolarProjects, project1.DEBAssetCode)
	// the above thing is to test the function itself and not the functionality since
	// DEBAssetCode for project1 should be empty
	err = recp.Save()
	if err != nil {
		t.Fatal(err)
	}
	err = UpdateProjectSlice(&recp, project1)
	if err != nil {
		t.Fatal(err)
	}
	err = nOP.RecipientAuthorizeContract(recp) // again a placeholder for testing this function, not the flow itself
	if err != nil {
		t.Fatal(err)
	}
	chk := nOP.CalculatePayback("100")
	if chk != "0.257143" {
		t.Fatalf("Balance doesn't match , quitting!")
	}
	var arr []SolarProject
	x, err := SelectContractByPrice(arr)
	if err == nil {
		t.Fatalf("Empty array returns choice")
	}
	y, err := SelectContractByTime(arr)
	if err == nil {
		t.Fatalf("Empty array returns choice")
	}
	arr = append(arr, nOP)
	x, err = SelectContractByPrice(arr)
	if err != nil {
		t.Fatal(err)
	}
	if x.Params.Index != nOP.Params.Index {
		t.Fatalf("Indices don't match, quitting!")
	}
	y, err = SelectContractByTime(arr)
	if err != nil {
		t.Fatal(err)
	}
	if y.Params.Index != nOP.Params.Index {
		t.Fatalf("Indices don't match, quitting!")
	}
	os.Remove(os.Getenv("HOME") + "/.openfinancing/database/" + "/yol.db")
}
