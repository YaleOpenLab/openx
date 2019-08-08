package main

/*
// newLivingUnitCoop creates a new living unit coop
func newLivingUnitCoop(mdate string, mrights string, stype string, intrate float64, rating string,
	bIssuer string, uWriter string, totalAmount float64, typeOfUnit string, monthlyPayment float64,
	title string, location string, description string) (opzones.LivingUnitCoop, error) {
	var coop opzones.LivingUnitCoop
	coop.MaturationDate = mdate
	coop.MemberRights = mrights
	coop.SecurityType = stype
	coop.InterestRate = intrate
	coop.Rating = rating
	coop.BondIssuer = bIssuer
	coop.Underwriter = uWriter
	coop.Title = title
	coop.Location = location
	coop.Description = description
	coop.DateInitiated = utils.Timestamp()

	x, err := opzones.RetrieveAllLivingUnitCoops()
	if err != nil {
		return coop, errors.Wrap(err, "could not retrieve all living unit coops")
	}
	coop.Index = len(x) + 1
	coop.UnitsSold = 0
	coop.Amount = totalAmount
	coop.TypeOfUnit = typeOfUnit
	coop.MonthlyPayment = monthlyPayment
	err = coop.Save()
	return coop, err
}

// newConstructionBond returns a New Construction Bond and automatically stores it in the db
func newConstructionBond(mdate string, stype string, intrate float64, rating string,
	bIssuer string, uWriter string, unitCost float64, itype string, nUnits int, tax string, recIndex int,
	title string, location string, description string) (opzones.ConstructionBond, error) {
	var cBond opzones.ConstructionBond
	cBond.MaturationDate = mdate
	cBond.SecurityType = stype
	cBond.InterestRate = intrate
	cBond.Rating = rating
	cBond.BondIssuer = bIssuer
	cBond.Underwriter = uWriter
	cBond.Title = title
	cBond.Location = location
	cBond.Description = description
	cBond.DateInitiated = utils.Timestamp()

	x, err := opzones.RetrieveAllConstructionBonds()
	if err != nil {
		return cBond, errors.Wrap(err, "could not retrieve all living unit coops")
	}

	cBond.Index = len(x) + 1
	cBond.CostOfUnit = unitCost
	cBond.InstrumentType = itype
	cBond.NoOfUnits = nUnits
	cBond.Tax = tax
	cBond.RecipientIndex = recIndex
	err = cBond.Save()
	return cBond, err
}

/*
	_, err = newConstructionBond("Dec 21 2021", "Security Type 1", 5.4, "AAA", "Moody's Investments", "Wells Fargo",
		200000, "Opportunity Zone Construction", 200, "5% tax for 10 years", 1, "India Basin Project", "San Francisco", "India Basin is an upcoming creative project based in San Francisco that seeks to host innovators from all around the world")
	if err != nil {
		return err
	}

	_, err = newConstructionBond("Apr 2 2025", "Security Type 2", 3.6, "AA", "Ant Financial", "People's Bank of China",
		50000, "Opportunity Zone Construction", 400, "No tax for 20 years", 1, "Shenzhen SEZ Development", "Shenzhen", "Shenzhen SEZ Development seeks to develop a SEZ in Shenzhen to foster creation of manufacturing jobs.")
	if err != nil {
		return err
	}

	_, err = newConstructionBond("Jul 9 2029", "Security Type 3", 4.2, "BAA", "Softbank Corp.", "Bank of Japan",
		150000, "Opportunity Zone Construction", 100, "3% Tax for 5 Years", 1, "Osaka Development Project", "Osaka", "This Project seeks to develop cutting edge technologies in Osaka")
	if err != nil {
		return err
	}

	_, err = newLivingUnitCoop("Dec 21 2021", "Member Rights Link", "Security Type 1", 5.4, "AAA", "Moody's Investments", "Wells Fargo",
		200000, "Coop Model", 4000, "India Basin Project", "San Francisco", "India Basin is an upcoming creative project based in San Francisco that seeks to host innovators from all around the world")
	if err != nil {
		return err
	}

	_, err = newLivingUnitCoop("Apr 2 2025", "Member Rights Link", "Security Type 2", 3.6, "AA", "Ant Financial", "People's Bank of China",
		50000, "Coop Model", 1000, "Shenzhen SEZ Development", "Shenzhen", "Shenzhen SEZ Development seeks to develop a SEZ in Shenzhen to foster creation of manufacturing jobs.")
	if err != nil {
		return err
	}

	_, err = newLivingUnitCoop("Jul 9 2029", "Member Rights Link", "Security Type 3", 4.2, "BAA", "Softbank Corp.", "Bank of Japan",
		150000, "Coop Model", 2000, "Osaka Development Project", "Osaka", "ODP seeks to develop cutting edge technologies in Osaka and invites investors all around the world to be a part of this new age")
	if err != nil {
		return err
	}
*/
