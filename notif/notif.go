package notif

import (
	"log"
	"net/smtp"

	utils "github.com/OpenFinancing/openfinancing/utils"
	"github.com/spf13/viper"
)

// package notif is used to send out notifications regarding important events that take
// place with respect to a specific project / investment

// MWTODO: Get comments on general text here
// TODO: maybe encrypt the config file so a person with access to this file cannot
// steal all credentials
// footerString is a common footer string that is used by all emails
var footerString = "Have a nice day!\n\nWarm Regards, \nThe OpenSolar Team\n\n\n\n" +
	"You're receiving this email because your contact was given" +
	" on the opensolar platform for receiving notifications on orders in which you're a party.\n\n\n"

// sendMail is a handler for sending out an email to an entity, reading required params
// from the config file
func sendMail(body string, to string) error {
	var err error
	// read from config.yaml in the working directory
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		log.Println("Error while reading email values from config file")
		return err
	}
	log.Println("SENDING EMAIL: ", viper.Get("email"), viper.Get("password"))
	from := viper.Get("email").(string)    // interface to string
	pass := viper.Get("password").(string) // interface to string
	auth := smtp.PlainAuth("", from, pass, "smtp.gmail.com")
	// to can also be an array of addresses if needed
	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: OpenSolar Notification\n\n" + body

	err = smtp.SendMail("smtp.gmail.com:587", auth, from, []string{to}, []byte(msg))
	if err != nil {
		log.Printf("smtp error: %s", err)
		return err
	}
	return nil
}

// SendInvestmentNotifToRecipient sends a notification to the recipient when an investor
// invests in an order he's the recipient of
func SendInvestmentNotifToRecipient(projIndex int, to string, recpPbTrustHash string, recpAssetHash string, recpDebtTrustHash string, recpDebtAssetHash string) error {
	// this is sent to the recipient on investment from an investor
	body := "Greetings from the opensolar platform! \n\n" +
		"We're writing to let you know that project number: " + utils.ItoS(projIndex) + " has been invested in.\n\n" +
		"Your proofs of payment are attached below and may be used as future reference in case of discrepancies:  \n\n" +
		"Your payback trusted asset hash is: https://testnet.steexp.com/tx/" + recpPbTrustHash + "\n" +
		"Your payback asset hash is: https://testnet.steexp.com/tx/" + recpAssetHash + "\n" +
		"Your debt trusted asset hash is: https://testnet.steexp.com/tx/" + recpDebtTrustHash + "\n" +
		"Your debt asset hash is: https://testnet.steexp.com/tx/" + recpDebtAssetHash + "\n\n\n" +
		footerString
	return sendMail(body, to)
}

// SendInvestmentNotifToInvestor sends a notification to the investor when he invests
// in a particular project
func SendInvestmentNotifToInvestor(projIndex int, to string, stableHash string, trustHash string, assetHash string) error {
	// this is sent to the investor on investment
	// this should ideally contain all the information he needs for a concise proof of
	// investment
	body := "Greetings from the opensolar platform! \n\n" +
		"We're writing to let you know have invested in project number: " + utils.ItoS(projIndex) + "\n\n" +
		"Your proofs of payment are attached below and may be used as future reference in case of discrepancies:  \n\n" +
		"Your stablecoin payment hash is: https://testnet.steexp.com/tx/" + stableHash + "\n" +
		"Your trusted asset hash is: https://testnet.steexp.com/tx/" + trustHash + "\n" +
		"Your investment asset hash is: https://testnet.steexp.com/tx/" + assetHash + "\n\n\n" +
		footerString
	return sendMail(body, to)
}

func SendSeedInvestmentNotifToInvestor(projIndex int, to string, stableHash string, trustHash string, assetHash string) error {
	// this is sent to the investor on investment
	// this should ideally contain all the information he needs for a concise proof of
	// investment
	body := "Greetings from the opensolar platform! \n\n" +
		"We're writing to let you know have invested in the seed round of project: " + utils.ItoS(projIndex) + "\n\n" +
		"Your proofs of payment are attached below and may be used as future reference in case of discrepancies:  \n\n" +
		"Your stablecoin payment hash is: https://testnet.steexp.com/tx/" + stableHash + "\n" +
		"Your trusted asset hash is: https://testnet.steexp.com/tx/" + trustHash + "\n" +
		"Your investment asset hash is: https://testnet.steexp.com/tx/" + assetHash + "\n\n\n" +
		footerString
	return sendMail(body, to)
}

// SendPaybackNotifToInvestor sends a notification email to the recipient when he
// pays back towards a particular order
func SendPaybackNotifToRecipient(projIndex int, to string, stableUSDHash string, debtPaybackHash string) error {
	// this is sent to the recipient
	body := "Greetings from the opensolar platform! \n\n" +
		"We're writing to let you know have paid back towards project number: " + utils.ItoS(projIndex) + "\n\n" +
		"Your proofs of payment are attached below and may be used as future reference in case of discrepancies:  \n\n" +
		"Stablecoin payment hash is: https://testnet.steexp.com/tx/" + stableUSDHash + "\n" +
		"Debt asset hash is: https://testnet.steexp.com/tx/" + debtPaybackHash + "\n\n\n" +
		footerString
	return sendMail(body, to)
}

// SendPaybackNotifToInvestor sends a notification email to the investor when the recipient
// pays back towards a particular order
func SendPaybackNotifToInvestor(projIndex int, to string, stableUSDHash string, debtPaybackHash string) error {
	// this is sent to the investor on payback from an investor
	body := "Greetings from the opensolar platform! \n\n" +
		"We're writing to let you know that the recipient has paid back towards project number: " + utils.ItoS(projIndex) + "\n\n" +
		"The recipient's proofs of payment are attached below and may be used as future reference in case of discrepancies:  \n\n" +
		"Stablecoin payment hash is: https://testnet.steexp.com/tx/" + stableUSDHash + "\n" +
		"Debt asset hash is: https://testnet.steexp.com/tx/" + debtPaybackHash + "\n\n\n" +
		footerString
	return sendMail(body, to)
}

// SendPaybackNotifToInvestor sends a notification email to the investor when the recipient
// pays back towards a particular order
func SendUnlockNotifToRecipient(projIndex int, to string) error {
	// this is sent to the investor on payback from an investor
	body := "Greetings from the opensolar platform! \n\n" +
		"We're writing to let you know that project number: " + utils.ItoS(projIndex) + " has been invested in\n\n" +
		"You are required to logon to the platform within a period of 3(THREE) days in order to accept the investment\n\n" +
		"If you choose to not accept the given investment in your project, please be warned that your reputation score " +
		"will be adjusted accordingly and this may affect any future proposal that you seek funding for on the platform\n\n" +
		footerString
	return sendMail(body, to)
}

func SendEmail(message string, to string, name string) error {
	// we can't send emails directly sicne we would need their gmail usernames and password for that
	startString := "Greetings from the opensolar platform! \n\n" +
		"We're writing to let you know that " + name + " has sent you a message. The message contents follow: \n\n"
	body := startString + message + "\n\n\n" + footerString
	return sendMail(body, to)
}

func SendAlertEmail(message string, to string) error {
	startString := "Greetings from the opensolar platform! \n\n" +
		"We're writing to let you know that you have received a message from the platform: \n\n\n" + message
	body := startString + "\n\n\n" + footerString
	return sendMail(body, to)
}

// we don't know if the user has paid
func SendPaybackAlertEmail(projIndex int, to string) error {
	startString := "Greetings from the opensolar platform! \n\n" +
		"This is a kind reminder to let you know that your payment is due this period for project numbered: " + utils.ItoS(projIndex) +
		"\n\n If you have already paid or have received a donation towards this month, please ignore this alert."
	body := startString + "\n\n\n" + footerString
	return sendMail(body, to)
}

// we know that the user is due in their payment
func SendNicePaybackAlertEmail(projIndex int, to string) error {
	startString := "Greetings from the opensolar platform! \n\n" +
		"This is a kind reminder to let you know that your payment is due this period for project numbered: " + utils.ItoS(projIndex) +
		"\n\n Please payback at the earliest."
	body := startString + "\n\n\n" + footerString
	return sendMail(body, to)
}

func SendSternPaybackAlertEmail(projIndex int, to string) error {
	startString := "Greetings from the opensolar platform! \n\n" +
		"We're writing to let you know that your payment is due this period for project numbered: " + utils.ItoS(projIndex) +
		"\n\n Please payback within two payback cycles to avoid re-routing of power services."
	body := startString + "\n\n\n" + footerString
	return sendMail(body, to)
}

func SendDisconnectionEmail(projIndex int, to string) error {
	startString := "Greetings from the opensolar platform! \n\n" +
		"We're writing to let you know that electricity produced from your project numbered: " + utils.ItoS(projIndex) +
		"\n\n Had been redirected towards the main power grid. Please contact your guarantor to resume services"
	body := startString + "\n\n\n" + footerString
	return sendMail(body, to)
}

func SendDisconnectionEmailI(projIndex int, to string) error {
	startString := "Greetings from the opensolar platform! \n\n" +
		"We're writing to let you know that electricity produced from your project numbered: " + utils.ItoS(projIndex) +
		"\n\n Had been redirected towards the main power grid due to irregualr payments by the recipient involved.\n\n" +
		"We are constantly monitoring this situation and will be continuing to send you emails on the same.\n\n" +
		"Please feel free to write to support with your queries in the meantime."
	body := startString + "\n\n\n" + footerString
	return sendMail(body, to)
}

func SendSternPaybackAlertEmailI(projIndex int, to string) error {
	startString := "Greetings from the opensolar platform! \n\n" +
		"We're writing to let you know that we are aware that payments towards the project: " + utils.ItoS(projIndex) +
		"\n\n haven't been made and we have reached out to the project recipient on the same. If this situation continues for " +
		"two more payment periods, we will be redirecting power towards the general grid and you would receive payments " +
		"for all periods where they were due. \n\n" +
		"We are constantly monitoring this situation and will be continuing to send you emails on the same.\n\n" +
		"Please feel free to write to support with your queries in the meantime."
	body := startString + "\n\n\n" + footerString
	return sendMail(body, to)
}

func SendSternPaybackAlertEmailG(projIndex int, to string) error {
	startString := "Greetings from the opensolar platform! \n\n" +
		"We're writing to let you know that we are aware that payments towards the project: " + utils.ItoS(projIndex) +
		"\n\n haven't been made and have reached out to the project recipient on the same. If this situation continues for " +
		"two more payment periods, we will be redirecting power towards the general grid and contact you for further" +
		"information on how the gurantee towards the project would be realized to investors.\n\n" +
		"We are constantly monitoring this situation and will be continuing to send you emails on the same.\n\n" +
		"Please feel free to write to support with your queries in the meantime."
	body := startString + "\n\n\n" + footerString
	return sendMail(body, to)
}

func SendDisconnectionEmailG(projIndex int, to string) error {
	startString := "Greetings from the opensolar platform! \n\n" +
		"We're writing to let you know that electricity produced from your project numbered: " + utils.ItoS(projIndex) +
		"\n\n Had been redirected towards the main power grid due to irregular payments by the recipient involved.\n\n" +
		"We will be reaching out to you in the coming days on how to proceed with realizing the guarantee towards this " +
		"project in order to safeguard investors. We will also be contacting the recipient involved to update them on the" +
		"situation and will make efforts to alleviate this problem as soon as possible." +
		"We are constantly monitoring this situation and will be continuing to send you emails on the same.\n\n" +
		"Please feel free to write to support with your queries in the meantime."
	body := startString + "\n\n\n" + footerString
	return sendMail(body, to)
}

func SendContractNotification(Hash1 string, Hash2 string, Hash3 string, Hash4 string, Hash5 string, to string) error {
	body := "Greetings from the opensolar platform! \n\n" +
		"We're writing to let you know that you have signed a contract\n\n" +
		"Your proofs of signing are attached below and may be used as future reference in case of discrepancies:  \n\n" +
		"Your first hash is: https://testnet.steexp.com/tx/" + Hash1 + "\n" +
		"Your second hash is: https://testnet.steexp.com/tx/" + Hash2 + "\n" +
		"Your third hash is: https://testnet.steexp.com/tx/" + Hash3 + "\n" +
		"Your fourth hash is: https://testnet.steexp.com/tx/" + Hash4 + "\n" +
		"Your fifth hash is: https://testnet.steexp.com/tx/" + Hash5 + "\n\n\n" +
		footerString
	return sendMail(body, to)
}