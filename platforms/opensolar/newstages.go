package opensolar

import ()

// newstages represents the final stage definition for the opensolar project. These are
// the stages that need to be enforced before moving on to subsequent stages. On the frontend,
// we shall be having a checklist for each stage and the stages themselves will be on a
// circle scale timeline where a person who's visiting the site can see how many stages are
// present for the particular platform (opensolar in this case)

// TODO: transition to this new stage definition from the older one. Initial step would be
// to do what the other stages do and then slowly migrate to this construction. But this
// would require knowledge on whether these stages are actually immutable since if they aren't,
// we require db handlers and stuff. MWTODO: are the stages immutable?

// StageNew is the evolution of the erstwhile stage construction
type StageNew struct {
	Number          int
	FriendlyName    string   // the informal name that one can use while referring to the stage (nice for UI as well)
	Name            string   // this is a more formal name to give to the given stage
	Activities      []string // the activities that are covered in this particular stage and need to be fulfilled in order to move to the next stage.
	StateTrigger    []string // trigger state change from n to n+1
	BreachCondition []string // define breach conditions for a particular stage
}

var Stage0 = &StageNew{
	Number:       0,
	FriendlyName: "Handshake",
	Name:         "Idea Consolidation",
	Activities: []string{
		"[Originator] proposes project and either secures or agrees to serve as [Solar Developer]",
		"NOTE: Originator is the community leader or catalyst for the project, they may opt to serve as the solar developer themselves, or pass that responsibility off, going forward we will use solar developer to represent the interest of both.",
		"[Solar Developer] creates general estimation of project (eg. with an automatic calculation through Google Project Sunroof, PV) ",
		"If [Originator]/[Solar Developer] is not landowner [Host] states legal ownership of site (hard proof is optional at this stage)",
	},
	StateTrigger: []string{
		"Matching of originator with receiver, and mutual approval/intention of interest.",
	},
}

var Stage1 = &StageNew{
	Number:       1,
	FriendlyName: "Engagement",
	Name:         "RFP Development",
	Activities: []string{
		"[Solar Developer] Analyse parameters, create financial model (proforma)",
		"[Host] & [Solar Developer] engage [Legal] & begin scoping site for planning constraints and opportunities (viability analysis)",
		"[Solar Developer] Create RFP (‘Request For Proposal’)",
		"Simple: Automatic calculation (eg. Sunroof style)",
		"Complex: Public project with 3rd party RFP consultant (independent engineer)",
		"[Originator][Solar Developer][Offtaker] Post project for RFP",
		"[Beneficiary/Host] Define and select RFP developer.",
		"[Investor] First angel investment option (high risk)",
		"Allow ‘time banking’ as sweat equity, monetized as tokenized capital or shadow stock",
	},
	StateTrigger: []string{
		"Issue an RFP",
		"Letter of Intent or MOU between originator and developer",
	},
}

var Stage2 = &StageNew{
	Number:       2,
	FriendlyName: "Quotes",
	Name:         "Actions",
	Activities: []string{
		"[Solar Developer][Beneficiary/Offtaker][Legal] PPA model negotiation.",
		"[Originator][Beneficiary]  Compare quotes from bidders: ",
		"[Engineering Procurement and Construction] (labor)",
		"[Vendors] (Hardware)",
		"[Insurers]",
		"[Issuer]",
		"[Intermediary Portal]",
		"[Originator/Receiver] Begin negotiation with [Utility]",
		"[Solar Developer] checks whether site upgrades are necessary.",
		"[Solar Developer][Host] Prepare submission for permitting and planning",
		"[Investor] Angel incorporation (less risk)",
	},
	StateTrigger: []string{
		"Selection of quotes and vendors",
		"Necessary identification of entities: Installers and offtaker",
	},
}

var Stage3 = &StageNew{
	Number:       3,
	FriendlyName: "Signing",
	Name:         "Contract Execution",
	Activities: []string{
		"[Solar Developer] pays [Legal] for PPA finalization.",
		"[Solar Developer][Host] Signs site Lease with landowner.",
		"[Solar Developer] OR [Issuer] signs Offering Agreement with [Intermediary Portal].",
		"[Solar Developer][Beneficiary] selects and signs contracts with: ",
		"[Engineering Procurement and Construction] (labor)",
		"[Vendors] (Hardware)",
		"[Insurers]",
		"[Issuer] OR [Intermediary Portal]",
		"[Offtaker] OR [Solar Developer][Engineering, Procurement and Construction] sign vendor/developer EPC Contracts",
		"[Solar Developer][Offtaker] signs PPA/Offtake Agreement",
		"[Investor] 2nd stage of eligible funding",
		"[Solar Developer][Beneficiary] makes downpayment to [Engineering Procurement and Construction] (labor)",
		"[Investor] Profile with risk ",
	},
	StateTrigger: []string{
		"Execution of contracts - Sign!",
	},
}

var Stage4 = &StageNew{
	Number:       4,
	FriendlyName: "The Raise",
	Name:         "Finance and Capitalization",
	Activities: []string{
		"[Issuer] engages [Intermediary Portal] to develop Form C or prospectus",
		"[Intermediary Portal] lists [Issuer] project",
		"[Originator][Solar Developer][Offtaker] market the crowdfunded offering",
		"[Investors] Commit capital to the project",
		"[Intermediary Portal] closes offering and disburses capital from Escrow account to [Issuers]",
		"If [Issuer] is not also [Solar Developer] then [Issuer] passes funds to [Solar Developer] ",
	},
	StateTrigger: []string{
		"Project account receives funds that cover the raise amount. Raise amount: normally includes both project capital expenditure (i.e. hardware and labor) and ongoing Operation & Management costs",
	},
}

var Stage5 = &StageNew{
	Number:       5,
	FriendlyName: "Construction",
	Name:         "Payments and Construction",
	Activities: []string{
		"[Solar Developer] coordinates installation dates and arrangements with [Host][Off-takers]",
		"[Solar Developer] OR [Engineering, Procurement and Construction] take delivery of equipment from [Vendor]",
		"[Utility] issues conditional interconnection",
		"[Solar Developer] schedules installation with [Engineering, Procurement and Construction]",
		"[Engineering, Procurement and Construction] completes installation.",
		"[Solar Developer] pays [Engineering, Procurement and Construction] for substantial completion of the project.",
		"[Insurers] verifies policy, [Solar Developer] pays [Insurers]",
		"[Investor] role?",
	},
	StateTrigger: []string{
		"Installation reaches substantial completion",
		"IoT devices detect energy generation",
	},
}

var Stage6 = &StageNew{
	Number:       6,
	FriendlyName: "Interconnection",
	Name:         "Contract Execution",
	Activities: []string{
		"[Solar Developer] coordinates with [Engineering Procurement and Construction] to schedule interconnection dates with [Utility] ",
		"[Engineering, Procurement and Construction] submits ‘as-built’ drawings to City/County Inspectors and schedules interconnection with [Utility]",
		"[Solar Developer] schedules City/County Building Inspector visit",
		"[Utility] visits site for witness test",
		"[Utility] places project in service ",
	},
	StateTrigger: []string{
		"[Utility] places project in service",
	},
}

var Stage7 = &StageNew{
	Number:       7,
	FriendlyName: "Legacy",
	Name:         "Operation and Management",
	Activities: []string{
		"[Solar Developer] hires OR becomes [Manager]",
		"[Manager] hires [Operations & Maintenance] provider",
		"[Manager] sets up billing system and issues monthly bills to [Offtaker] and collects payment on bills",
		"[Manager] monitors for breaches of payment or contract, other indentures, force majeure or adverse conditions [see below for Breach Conditions]",
		"[Manager] files annual taxes",
		"[Manager] handles annual true-up on net-metering payments",
		"[Manager] makes annual cash distributions and issues 1099-DIV to [Investors] or coordinates share repurchase from [Investors]",
		"If applicable, [Manager] executes flip between [Solar Developer] ownership interest and [Tax equity investor]",
		"[Manager] OR [Operations & Maintenance] monitors system performance and coordinates with [Off-takers] to schedule routine maintenance",
		"[Manager] OR [Operations & Maintenance] coordinates with [Engineering, Procurement and Construction] to change inverters or purchase replacements from [Vendors] as needed.",
		"[Investors] can engage in secondary market (i.e. re-selling its securities). ",
	},
	StateTrigger: []string{
		"[Investors] reach preferred return rate, or Power Purchase Agreement stipulates ownership flip date or conditions ",
	},
	BreachCondition: []string{
		"[Offtaker] fails to make $/kWh payments after X period of time due. ",
	},
}

var Stage8 = &StageNew{
	Number:       8,
	FriendlyName: "Handoff",
	Name:         "Ownership Flip",
	Activities: []string{
		"[Beneficiary/Offtakers] Payments accrue to cover the [Investor] principle (i.e. total raised amount)",
		"Escrow account (eg. capital account) pays off principle to [Investor]",
	},
	StateTrigger: []string{
		"[Beneficiary] (eg. Host, Holding)  becomes full legal owner of physical assets",
		"[Investors] exit the project",
	},
}

var Stage9 = &StageNew{
	Number:       9,
	FriendlyName: "End of Life",
	Name:         "Disposal",
	Activities: []string{
		"[IoT] Solar equipment is generating below a productivity threshold, or shows general malfunction",
		"[Beneficiaries][Developers]  dispose of the equipment to a recycling program",
		"[Developer/Recycler] Certifies equipment is received",
	},
	StateTrigger: []string{
		"Project termination",
		"Wallet terminations",
	},
}
