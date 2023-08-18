package getresidence

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/efficientgo/core/errors"

	"github.com/uptrace/bunrouter"
)

func nonce() string {
	var b [32]byte
	_, err := rand.Read(b[:])
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(b[:])
}

func writeDefaultHeaders(h http.Header, nonce string) {
	csp := "default-src 'none'; img-src data: imagedelivery.net; font-src use.typekit.net; connect-src 'self'"
	if nonce != "" {
		csp = fmt.Sprintf("%s; script-src 'nonce-%s' 'strict-dynamic'; style-src 'nonce-%s'", csp, nonce, nonce)
	}
	h.Set("X-Content-Type-Options", "nosniff")
	h.Set("Content-Type", "text/html; charset=utf-8")
	h.Set("Content-Security-Policy", csp)
	h.Set("Referrer-Policy", "no-referrer")
	h.Set("Permissions-Policy", "accelerometer=(), ambient-light-sensor=(), autoplay=(), battery=(), camera=(), display-capture=(), document-domain=(), encrypted-media=(), execution-while-not-rendered=(), execution-while-out-of-viewport=(), fullscreen=(), gamepad=(), geolocation=(), gyroscope=(), hid=(), identity-credentials-get=(), idle-detection=(), local-fonts=(), magnetometer=(), microphone=(), midi=(), otp-credentials=(), payment=(), picture-in-picture=(), publickey-credentials-create=(), publickey-credentials-get=(), screen-wake-lock=(), serial=(), speaker-selection=(), storage-access=(), usb=(), web-share=(), xr-spatial-tracking=()")
	h.Set("X-Frame-Options", "DENY")
}

func emojiToFavicon(emoji string) string {
	b := strings.Builder{}
	b.WriteString(`data:image/svg+xml,<svg xmlns=%22http://www.w3.org/2000/svg%22 viewBox=%220 0 100 100%22><text y=%22.9em%22 font-size=%2290%22>`)
	b.WriteString(emoji)
	b.WriteString(`</text></svg>`)
	return b.String()
}

func getSession(r bunrouter.Request) (string, error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return "", errors.Wrap(err, "session cookie")
	}

	return cookie.Value, nil
}
