package utils

import "strings"

func BuildSupplierInviteEmail(supplierName, signupLink, invitationMessage string) string {
	var b strings.Builder

	b.WriteString(`<!DOCTYPE html>
<html>
  <body style="font-family: Arial, sans-serif; background-color:#f9fafb; padding:20px;">
    <div style="max-width:600px; margin:0 auto; background:#ffffff; padding:24px; border-radius:6px; color:#111827;">
      
      <h2 style="margin-top:0;">You're invited to join Bambamload</h2>

      <p>Dear `)
	b.WriteString(supplierName)
	b.WriteString(`,</p>

      <p>
        You’ve been invited to join <strong>Bambamload</strong>, a platform that helps suppliers
        manage loads, transactions, and operations more efficiently.
      </p>

      <p>
        Click the button below to create your supplier account and get started.
      </p>

      <div style="margin:30px 0; text-align:center;">
        <a href="`)
	b.WriteString(signupLink)
	b.WriteString(`"
           style="
             background-color:#2563eb;
             color:#ffffff;
             padding:12px 24px;
             text-decoration:none;
             border-radius:4px;
             font-weight:bold;
             display:inline-block;
           ">
          Sign Up on Bambamload
        </a>
      </div>

      <p style="font-size:14px; color:#374151;">
        If the button doesn’t work, copy and paste this link into your browser:
      </p>

      <p style="font-size:14px; word-break:break-all;">
        `)
	b.WriteString(signupLink)

	b.WriteString(`
      </p>

      <h3 style="margin-top:32px;">Invitation Message</h3>

      <p style="white-space:pre-line;">
        `)
	b.WriteString(invitationMessage)

	b.WriteString(`
      </p>

      <p style="margin-top:30px;">
        We’re excited to have you on board.<br/>
        <strong>The Bambamload Team</strong>
      </p>

    </div>
  </body>
</html>`)

	return b.String()
}
