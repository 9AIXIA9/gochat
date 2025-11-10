package gomail

const userCreatedEmailFormat = `<!DOCTYPE html>
<html lang="en" style="font-family:Segoe UI,Arial,sans-serif;background:#f5f7fa;margin:0;padding:0;">
  <head>
    <meta charset="UTF-8"/>
    <title>Welcome to GoChat</title>
  </head>
  <body style="background:#f5f7fa;margin:0;padding:32px;">
    <table role="presentation" width="100%%" style="max-width:600px;margin:0 auto;background:#ffffff;border-radius:12px;overflow:hidden;box-shadow:0 4px 12px rgba(0,0,0,0.08);">
      <tr>
        <td style="background:#4e46e5;padding:24px;">
          <h1 style="color:#ffffff;font-size:24px;margin:0;">Welcome to GoChat</h1>
        </td>
      </tr>
      <tr>
        <td style="padding:32px;">
          <p style="font-size:15px;color:#333;line-height:1.6;margin:0 0 16px;">
            Hi, your account has been successfully created.
          </p>
          <p style="font-size:14px;color:#555;line-height:1.6;margin:0 0 24px;">
            Your user number is <strong style="color:#4e46e5;">%s</strong>. Please keep it safe.
          </p>
          <a href="#login" style="display:inline-block;background:#4e46e5;color:#ffffff;text-decoration:none;padding:12px 24px;border-radius:6px;font-weight:600;font-size:14px;">
            Get Started
          </a>
          <p style="font-size:12px;color:#999;margin:32px 0 0;">
            If you did not sign up, please ignore this email.
          </p>
        </td>
      </tr>
      <tr>
        <td style="background:#f0f2f6;padding:16px;text-align:center;">
          <p style="font-size:11px;color:#7a8899;margin:0;">© GoChat. All rights reserved.</p>
        </td>
      </tr>
    </table>
  </body>
</html>`
