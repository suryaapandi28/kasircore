package whatsapp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/suryaapandi28/kasircore/internal/entity"
)

type WhatsappSender struct {
	config *entity.Config
}

func NewWhatsappSender(cfg *entity.Config) *WhatsappSender {
	return &WhatsappSender{config: cfg}
}

type fonntePayload struct {
	Target      string `json:"target"`
	Message     string `json:"message"`
	CountryCode string `json:"countryCode,omitempty"`
}

func (w *WhatsappSender) SendMessage(phone, message string) error {
	if phone == "" {
		return errors.New("nomor whatsapp kosong")
	}

	if w.config == nil {
		return errors.New("config whatsapp kosong")
	}

	payload := fonntePayload{
		Target:      phone,
		Message:     message,
		CountryCode: "62",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s/send", "https://api.fonnte.com"),
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "!yyK17TPNb9kTcLgaXYv")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("gagal mengirim whatsapp, status: %d", resp.StatusCode)
	}

	return nil
}
