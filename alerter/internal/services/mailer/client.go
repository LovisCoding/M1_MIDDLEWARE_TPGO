package mailer

import (
	"bytes"
	"encoding/json"
	"io"
	"middleware/alerter/internal/models"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

var (
	mailerURL  = "https://mail-api.edu.forestier.re/mail"
	apiToken   = "MmknIoaZRvjmxSMQVbrEwdMTXxMYOHdacTvPFWwy"
	senderMail = "no-reply@etu.uca.fr"
)

func init() {
	if token := os.Getenv("MAILER_API_TOKEN"); token != "" {
		apiToken = token
	}
	if sender := os.Getenv("SENDER_MAIL"); sender != "" {
		senderMail = sender
	}
}

func SendEmail(to []string, subject, body string) error {
	client := &http.Client{}

	for _, recipient := range to {
		reqBody := models.EmailRequest{
			Recipient: recipient,
			Subject:   subject,
			Content:   body,
		}

		jsonData, err := json.Marshal(reqBody)
		if err != nil {
			return err
		}

		req, err := http.NewRequest("POST", mailerURL, bytes.NewBuffer(jsonData))
		if err != nil {
			return err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", apiToken)

		resp, err := client.Do(req)
		if err != nil {
			logrus.Errorf("Failed to send email to %s: %v", recipient, err)
			continue
		}

		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
			logrus.Errorf("Failed to send email to %s, status: %d, body: %s", recipient, resp.StatusCode, string(respBody))
			continue
		}

		logrus.Infof("Email sent successfully to %s", recipient)
	}

	return nil
}
