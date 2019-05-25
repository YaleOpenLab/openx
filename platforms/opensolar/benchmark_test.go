// +build all

package opensolar

import (
	"log"
	"os"
	"testing"

	database "github.com/YaleOpenLab/openx/database"
)

func CleanDB() {
	os.Remove(os.Getenv("HOME") + "/.openx/database/" + "/yol.db")
	os.MkdirAll(os.Getenv("HOME")+"/.openx/database", os.ModePerm)
}

func PopulateProjects() {
	var dummy Project
	dummy.PanelSize = "100 1000 sq.ft homes each with their own private spaces for luxury"
	dummy.TotalValue = 14000
	dummy.MoneyRaised = 0
	dummy.Metadata = "India Basin is an upcoming creative project based in San Francisco that seeks to invite innovators from all around to participate"
	dummy.InvestorAssetCode = ""
	dummy.DebtAssetCode = ""
	dummy.PaybackAssetCode = ""
	dummy.DateInitiated = ""
	dummy.Stage = 3

	for i := 0; i < 1000; i++ {
		dummy.Index = i
		dummy.Save()
	}

	aP, err := RetrieveAllProjects()
	if err != nil {
		log.Fatal(err)
	}

	if len(aP) != 999 {
		log.Println("Projects not inserted properly", len(aP))
		log.Fatal("No projects were inserted")
	}
}

func PopulateContractor() {
	for i := 1; i < 1000; i++ {
		_, _ = NewContractor("ContTest", "pwd", "blah", "NameContTest", "123 DEF Street", "ContDescription")
	}

	iA, err := RetrieveAllEntities("contractor")
	if err != nil {
		log.Fatal(err)
	}

	if len(iA) != 999 {
		log.Fatal("Couldn't populate db, quitting!")
	}
}

func PopulateDeveloper() {
	for i := 1; i < 1000; i++ {
		_, _ = NewDeveloper("DevTest", "pwd", "blah", "NameContTest", "123 DEF Street", "DevDescription")
	}

	iA, err := RetrieveAllEntities("developer")
	if err != nil {
		log.Fatal(err)
	}

	if len(iA) != 999 {
		log.Fatal("Couldn't populate db, quitting!")
	}
}

func PopulateOrig() {
	for i := 1; i < 1000; i++ {
		_, _ = NewOriginator("OrigTest", "pwd", "blah", "NameOrigTest", "123 ABC Street", "OrigDescription")
	}

	iA, err := RetrieveAllEntities("originator")
	if err != nil {
		log.Fatal(err)
	}

	if len(iA) != 999 {
		log.Fatal("Couldn't populate db, quitting!")
	}
}

func BenchmarkPopulateOrig(b *testing.B) {
	CleanDB()
	b.ResetTimer()
	PopulateOrig()
	b.StopTimer()
}

func BenchmarkPopulateCont(b *testing.B) {
	b.ResetTimer()
	PopulateContractor()
	b.StopTimer()
}

func BenchmarkPopulateProject(b *testing.B) {
	b.ResetTimer()
	PopulateProjects()
	b.StopTimer()
}

func BenchmarkPopulateDeveloper(b *testing.B) {
	b.ResetTimer()
	PopulateDeveloper()
	b.StopTimer()
}

func BenchmarkNewOrig(b *testing.B) {
	// populate the db wityh test values
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = newEntity("OrigTest", "pwd", "blah", "NameOrigTest", "123 ABC Street", "OrigDescription", "originator")
	}
}

func BenchmarkNewCont(b *testing.B) {
	// populate the db wityh test values
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = newEntity("ContTest", "pwd", "blah", "NameContTest", "123 DEF Street", "ContDescription", "contractor")
	}
}

func BenchmarNewProject(b *testing.B) {
	var dummy Project
	dummy.PanelSize = "100 1000 sq.ft homes each with their own private spaces for luxury"
	dummy.TotalValue = 14000
	dummy.MoneyRaised = 0
	dummy.Metadata = "India Basin is an upcoming creative project based in San Francisco that seeks to invite innovators from all around to participate"
	dummy.InvestorAssetCode = ""
	dummy.DebtAssetCode = ""
	dummy.PaybackAssetCode = ""
	dummy.DateInitiated = ""
	dummy.Stage = 3

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dummy.Index = i
		dummy.Save()
	}
}

func BenchmarkRetrieveEntity(b *testing.B) {
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		_, _ = RetrieveEntity(i)
	}
}

func BenchmarkValidateEntity(b *testing.B) {
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		_, _ = ValidateEntity("ContTest", "pwd")
	}
}

func BenchmarkRetrieveProject(b *testing.B) {
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		_, _ = RetrieveProject(i)
	}
}

func BenchmarkRetrieveAllProjects(b *testing.B) {
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		_, _ = RetrieveAllProjects()
	}
}

func BenchmarkRetrieveProjectAtStage(b *testing.B) {
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		_, _ = RetrieveProjectsAtStage(1)
	}
}

func BenchmarkProposeContract(b *testing.B) {
	x1, err := RetrieveEntity(1)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		_, _ = x1.Propose("100 16x32 panels", 28000, "Puerto Rico", 6, "LEED+ Gold rated panels and this is random data out of nowhere and we supply our own devs and provide insurance guarantee as well. Dual audit maintenance upto 1 year. Returns capped as per defaults", 1, 1, "blind")
	}
}

func BenchmarkOriginateContract(b *testing.B) {
	x1, err := RetrieveEntity(1)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		_, _ = x1.Originate("100 16x24 panels on a solar rooftop", 14000, "Puerto Rico", 5, "ABC School in XYZ peninsula", 1, "blind") // 1 is the index for martin
	}
}

func BenchmarkSetPO(b *testing.B) {
	project, err := RetrieveProject(1)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		_ = project.SetStage(0)
	}
}

func BenchmarkSetFP(b *testing.B) {
	project, err := RetrieveProject(1)
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		_ = project.SetStage(4)
	}
}

func BenchmarkDeleteKeyFromBucket(b *testing.B) {
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		_ = database.DeleteKeyFromBucket(i, database.ProjectsBucket)
	}
}

func BenchmarkGetTopReputationEWOR(b *testing.B) {
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		_, _ = TopReputationEntitiesWithoutRole()
	}
}

func BenchmarkGetTopReputationEWR(b *testing.B) {
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		_, _ = TopReputationEntities("contractor")
	}
}

func BenchmarkAddCollateral(b *testing.B) {
	contractor, err := NewContractor("ContTest", "pwd", "blah", "NameContTest", "123 DEF Street", "ContDescription")
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		_ = contractor.AddCollateral(100.5, "we're adding sample collateral")
	}
}

func BenchmarkSlashContractor(b *testing.B) {
	contractor, err := NewContractor("ContTest", "pwd", "blah", "NameContTest", "123 DEF Street", "ContDescription")
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		_ = contractor.Slash(100.5)
	}
}
