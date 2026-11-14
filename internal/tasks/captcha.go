package tasks

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Секрет hCaptcha хранится только на сервере (HCAPTCHA_SECRET).
// Если он не задан, проверка пропускается (режим разработки) и
// в лог пишется предупреждение один раз.
var captchaSecretWarned bool

// hcaptchaVerifyTimeout ограничивает время запроса к siteverify.
var hcaptchaVerifyTimeout = 10 * time.Second

type hcaptchaSiteverifyResponse struct {
	Success     bool     `json:"success"`
	ChallengeTS string   `json:"challenge_ts"`
	Hostname    string   `json:"hostname"`
	ErrorCodes  []string `json:"error-codes"`
}

// VerifyCaptchaToken проверяет токен капчи запросом
// POST https://api.hcaptcha.com/siteverify (secret, response, remoteip).
// Токен приходит с фронтенда вместе с данными формы.
func VerifyCaptchaToken(ctx context.Context, token string, remoteip string) error {
	secret := runtimeSettings.CaptchaSecret
	if secret == "" {
		if !captchaSecretWarned {
			captchaSecretWarned = true
			log.Print("HCAPTCHA_SECRET не задан — проверка hCaptcha пропускается (режим разработки)")
		}
		return nil
	}

	if strings.TrimSpace(token) == "" {
		return ErrCaptchaRequired
	}

	form := url.Values{
		"secret":   {secret},
		"response": {token},
	}
	if remoteip != "" {
		form.Set("remoteip", remoteip)
	}

	ctx, cancel := context.WithTimeout(ctx, hcaptchaVerifyTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.hcaptcha.com/siteverify", strings.NewReader(form.Encode()))
	if err != nil {
		return ErrCaptchaFailed
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: hcaptchaVerifyTimeout}
	resp, err := client.Do(req)
	if err != nil {
		// Ошибка сети не должна пропускать запрос без проверки
		return ErrCaptchaFailed
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return ErrCaptchaFailed
	}

	var result hcaptchaSiteverifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ErrCaptchaFailed
	}

	if !result.Success {
		return ErrCaptchaFailed
	}
	return nil
}
