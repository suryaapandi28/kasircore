package email

import (
	"fmt"
	"time"

	"github.com/suryaapandi28/kasircore/internal/entity"

	"gopkg.in/gomail.v2"
)

type EmailSender struct {
	Config *entity.Config
}

func NewEmailSender(config *entity.Config) *EmailSender {
	return &EmailSender{Config: config}
}

func (e *EmailSender) SendEmail(to []string, subject, body string) error {
	from := "no-reply@apsmdev.com"
	password := e.Config.SMTP.Password
	smtpHost := e.Config.SMTP.Host

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", from)
	mailer.SetHeader("To", to...)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", body)

	dialer := gomail.NewDialer(smtpHost, 587, from, password)
	err := dialer.DialAndSend(mailer)
	if err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	return nil
}

//	func (e *EmailSender) SendWelcomeEmail(to, name string) error {
//		subject := "Welcome Email | Depublic"
//		body := fmt.Sprintf("Dear %s,\nThis is a welcome email message from depublic\n\nDepublic Team", name)
//		return e.SendEmail([]string{to}, subject, body)
//	}
func (e *EmailSender) SendWelcomeEmail(to, name string, createTime time.Time) error {
	subject := "Selamat Datang di APSM Indonesia Global 🎉"

	createdAt := createTime.Format("02 January 2006 15:04")

	body := fmt.Sprintf(`
	<html>
	<body style="font-family: Arial, sans-serif; background-color:#f9f9f9; padding:20px;">
		<table width="100%%" cellspacing="0" cellpadding="0">
			<tr>
				<td align="center">
					<table width="600" style="background:#ffffff; padding:30px; border-radius:8px; box-shadow:0 2px 6px rgba(0,0,0,0.1);">
						
						<tr>
							<td align="center" style="font-size:22px; font-weight:bold; color:#333;">
								Selamat Datang di APSM Indonesia Global 🎉
							</td>
						</tr>
						
						<tr>
							<td style="padding:20px 0; font-size:16px; color:#555;">
								Dear <strong>%s</strong>,<br><br>
								Terima kasih telah mendaftar di <b>APSM Indonesia Global</b>!<br>
								Akun Anda berhasil dibuat dengan detail berikut:<br><br>

								<table cellpadding="6" cellspacing="0" style="border:1px solid #ddd; border-radius:6px; font-size:15px; color:#333;">
									<tr>
										<td style="background:#f4f4f4; font-weight:bold;">Email</td>
										<td>%s</td>
									</tr>
									<tr>
										<td style="background:#f4f4f4; font-weight:bold;">Status</td>
										<td>Belum Terverifikasi</td>
									</tr>
									<tr>
										<td style="background:#f4f4f4; font-weight:bold;">Tanggal Pendaftaran</td>
										<td>%s</td>
									</tr>
								</table>

								<br>
								Untuk memudahkan akses dan menjaga keamanan akun Anda, 
								silakan lakukan proses <b>verifikasi akun</b> terlebih dahulu. 
								Setelah verifikasi selesai, Anda dapat menggunakan seluruh layanan APSM tanpa batasan.
							</td>
						</tr>

						<tr>
							<td style="padding-top:30px; font-size:14px; color:#555;">
								Hormat kami,<br>
								<b>APSM Indonesia Global</b>
							</td>
						</tr>

					</table>
				</td>
			</tr>
		</table>
	</body>
	</html>
	`, name, to, createdAt)

	return e.SendEmail([]string{to}, subject, body)
}

func (e *EmailSender) SendResetPasswordEmail(to, name, resetCode string) error {
	subject := "Reset Password | Depublic"
	body := fmt.Sprintf("Dear %s,\nPlease use the following code to reset your password: %s\n\nDepublic Team", name, resetCode)
	return e.SendEmail([]string{to}, subject, body)
}

// func (e *EmailSender) SendVerificationEmail(to, name, code string) error {
// 	subject := "Verify Your Email | Depublic"
// 	body := fmt.Sprintf("Dear %s,\nPlease use the following code to verify your email: %s\n\nDepublic Team", name, code)
// 	return e.SendEmail([]string{to}, subject, body)
// }

func (e *EmailSender) SendVerificationEmail(to, name, code string) error {
	subject := "Verify Your Email | APSM Indonesia Global"

	body := fmt.Sprintf(`
    <html>
    <body style="font-family: Arial, sans-serif; background-color:#f9f9f9; padding:20px;">
        <table width="100%%" cellspacing="0" cellpadding="0">
            <tr>
                <td align="center">
                    <table width="600" style="background:#ffffff; padding:30px; border-radius:8px; box-shadow:0 2px 6px rgba(0,0,0,0.1);">
                        <tr>
                            <td align="center" style="font-size:20px; font-weight:bold; color:#333;">
                                Email Verification
                            </td>
                        </tr>
                        <tr>
                            <td style="padding:20px 0; font-size:16px; color:#555;">
                                Dear <strong>%s</strong>,<br><br>
                                Thank you for registering with <b>APSM Indonesia Global</b>.<br>
                                Please use the following code to verify your email:
                            </td>
                        </tr>
                        <tr>
                            <td align="center" style="padding:20px 0;">
                                <div style="font-size:24px; font-weight:bold; color:#2b6cb0; letter-spacing:4px;">
                                    %s
                                </div>
                            </td>
                        </tr>
                        <tr>
                            <td style="padding-top:20px; font-size:14px; color:#777;">
                                This code will expire in 5 minutes.<br>
                                If you did not request this, please ignore this email.
                            </td>
                        </tr>
                        <tr>
                            <td style="padding-top:30px; font-size:14px; color:#555;">
                                Regards,<br>
                                <b>APSM Indonesia Global</b>
                            </td>
                        </tr>
                    </table>
                </td>
            </tr>
        </table>
    </body>
    </html>
    `, name, code)

	return e.SendEmail([]string{to}, subject, body)
}

func (e *EmailSender) SendTransactionInfo(to, Transactions_id, Cart_id, User_id,
	Fullname_user, Trx_date, Payment, Payment_url, Amount string) error {
	subject := "Transaction Info | Depublic"
	body := fmt.Sprintf("Dear %s,\nThis is your transaction info from Depublic:\n\nTransaction ID: %s\n\nCart ID: %s\n\nUser ID: %s\n\nFullname: %s\n\nTransaction Date: %s\n\nPayment Type: %s\n\nURL Payment: %s\n\nTotal Amount: %s\n\n\nDepublic Team ",
		Fullname_user, Transactions_id, Cart_id, User_id, Fullname_user, Trx_date, Payment, Payment_url, Amount)
	return e.SendEmail([]string{to}, subject, body)
}
