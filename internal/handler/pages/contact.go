package pages

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/mail"
	"net/smtp"
	"strings"

	"daemontalk/internal/handler/common"
	"daemontalk/internal/i18n"
)

func (h *Handler) Contact(w http.ResponseWriter, r *http.Request) {
	lang := common.LangFromRequest(r)
	ui := i18n.Get(lang)

	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	if r.PostFormValue("website") != "" {
		fmt.Fprintf(w, `<p class="text-[var(--c-ok)]">%s</p>`, ui.Contact_Success)
		return
	}

	name := strings.TrimSpace(r.PostFormValue("name"))
	email := strings.TrimSpace(r.PostFormValue("email"))
	message := strings.TrimSpace(r.PostFormValue("message"))

	if name == "" || email == "" || message == "" || !ValidEmail(email) {
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

func ValidEmail(s string) bool {
	addr, err := mail.ParseAddress(s)
	if err != nil {
		return false
	}
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

	subject := common.StripCRLF(fmt.Sprintf("Portfolio contact: %s <%s>", fromName, fromEmail))
	msgBody := fmt.Sprintf("From: %s <%s>\r\n\r\n%s", common.StripCRLF(fromName), common.StripCRLF(fromEmail), body)

	msg := []byte("To: " + to + "\r\n" +
		"From: " + h.SMTPUser + "\r\n" +
		"Reply-To: " + common.StripCRLF(fromEmail) + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"Content-Type: text/plain; charset=UTF-8\r\n" +
		"\r\n" +
		msgBody)

	auth := smtp.PlainAuth("", h.SMTPUser, h.SMTPPass, host)
	return smtp.SendMail(addr, auth, h.SMTPUser, []string{to}, msg)
}
