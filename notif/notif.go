package notif

import (
	email "github.com/Varunram/essentials/email"
)

var footerString = "Have a nice day!\n\nWarm Regards, \nThe Openx Team\n\n\n\n" +
	"You're receiving this email because your contact was given" +
	" on the opensolar platform for receiving notifications on orders in which you're a party.\n\n\n"

// SendSecretsEmail is an email to trusted social contacts notifying that a user has shared a secret with them
func SendSecretsEmail(userEmail string, email1 string, email2 string, email3 string, secret1 string, secret2 string, secret3 string) error {
	bodyBase := "Greetings from the opensolar platform! \n\nWe're writing to let you know that user with email: " + userEmail +
		" has designated you as a trusted entity. Towards this, we request that you keep the attached secret in a safe and secure place and provide " +
		"it to the above user in case they request for it. \n\n" + "SECRET:\n\n"
	body1 := bodyBase + secret1 + "\n\n\n" + footerString
	err := email.SendMail(body1, email1)
	if err != nil {
		return err
	}

	body2 := bodyBase + secret2 + "\n\n\n" + footerString
	err = email.SendMail(body2, email2)
	if err != nil {
		return err
	}

	body3 := bodyBase + secret3 + "\n\n\n" + footerString
	err = email.SendMail(body3, email3)
	if err != nil {
		return err
	}

	return nil
}

// SendPasswordResetEmail sends a password reset email to the email address of the user
func SendPasswordResetEmail(to string, vCode string) error {
	body := "Greetings from the opensolar platform! \n\nWe're writing to let you know that you requested a password reset recently\n\n" +
		"Please use this given code along with the link attached in order to reset your password\n\n" +
		"VERIFICATION CODE: " + vCode + "\n\n\n" + footerString

	return email.SendMail(body, to)
}

// SendUserConfEmail sends a registration confirmation email to the email address of the user
func SendUserConfEmail(to string, code string) error {
	body := "Greetings from the opensolar platform! \n\nWe're writing to let you know that you requested a new account recently\n\n" +
		"Please input this code into the confirmation dialogue displayed on your screen to confirm your registration on openx\n\n" +
		"CONFIRMATION CODE: " + code + "\n\n\n" + footerString

	return email.SendMail(body, to)
}
