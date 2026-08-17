package handler

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/mail"
	"net/smtp"
	"strings"

	"daemontalk/internal/i18n"
)

func (h *Handler) Contact(w http.ResponseWriter, r *http.Request) {
	lang := langFromRequest(r)
	ui := i18n.Get(lang)

	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// Honeypot: bots fill this field
	if r.PostFormValue("website") != "" {
		fmt.Fprintf(w, `<p class="text-[var(--c-ok)]">%s</p>`, ui.Contact_Success)
		return
	}

	name := strings.TrimSpace(r.PostFormValue("name"))
	email := strings.TrimSpace(r.PostFormValue("email"))
	message := strings.TrimSpace(r.PostFormValue("message"))

	if name == "" || email == "" || message == "" || !validEmail(email) {
		fmt.Fprintf(w, `<p class="text-[var(--c-warn)]">%s</p>`, ui.Contact_Error)
		return
	}

	if h.SMTPHost != "" {
		if err := h.sendContactEmail(name, email, message); err != nil {
			slog.Error("send contact email failed", "error", err, "from_email", email)
			fmt.Fprintf(w, `<p class="text-[var(--c-warn)]">%s</p>`, ui.Contact_Error)
			return
		}
	} else {
		slog.Info("contact received [no SMTP configured]", "name", name, "email", email, "message", message)
	}

	fmt.Fprintf(w, `<p class="text-[var(--c-ok)]">%s</p>`, ui.Contact_Success)
}

// validEmail reports whether s parses as a single RFC 5322 address.
func validEmail(s string) bool {
	addr, err := mail.ParseAddress(s)
	if err != nil {
		return false
	}
	// Reject inputs that smuggle a display name (e.g. "x <a@b.c>") so the
	// stored/sent value is exactly the address the user typed.
	return addr.Address == s
}

func (h *Handler) sendContactEmail(fromName, fromEmail, body string) error {
	host := h.SMTPHost
	port := h.SMTPPort
	if port == "" {
		port = "587"
	}
	addr := net.JoinHostPort(host, port)

	to := h.SMTPTo
	if to == "" {
		to = h.SMTPUser
	}

	// Strip CR/LF from any value placed in a header to prevent header
	// injection. The message body is the only place untrusted multi-line
	// text is allowed.
	subject := stripCRLF(fmt.Sprintf("Portfolio contact: %s <%s>", fromName, fromEmail))
	msgBody := fmt.Sprintf("From: %s <%s>\r\n\r\n%s", stripCRLF(fromName), stripCRLF(fromEmail), body)

	msg := []byte("To: " + to + "\r\n" +
		"From: " + h.SMTPUser + "\r\n" +
		"Reply-To: " + stripCRLF(fromEmail) + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"\r\n" +
		msgBody)

	auth := smtp.PlainAuth("", h.SMTPUser, h.SMTPPass, host)
	return smtp.SendMail(addr, auth, h.SMTPUser, []string{to}, msg)
}

// stripCRLF removes carriage returns and newlines so untrusted input can't
// inject extra SMTP headers.
func stripCRLF(s string) string {
	return strings.NewReplacer("\r", "", "\n", "").Replace(s)
}
