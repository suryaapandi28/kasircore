package whatsapp

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type WhatsappSender struct {
	Token  string
	Client *http.Client
}

func NewWhatsappSender(token string) *WhatsappSender {
	return &WhatsappSender{
		Token: token,
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (f *WhatsappSender) SendOtp(phone, name, otp string) error {
	message := "Halo " + name + ",\n\n" +
		"Kode OTP Anda adalah:\n*" + otp + "*\n\n" +
		"Berlaku selama 5 menit.\n" +
		"Jangan bagikan kode ini kepada siapa pun."

	data := url.Values{}
	data.Set("target", phone) // contoh: 08123456789
	data.Set("message", message)
	data.Set("countryCode", "62") // opsional

	req, err := http.NewRequest(
		"POST",
		"https://api.fonnte.com/send",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", f.Token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := f.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("gagal mengirim whatsapp OTP via fonnte")
	}

	return nil
}
