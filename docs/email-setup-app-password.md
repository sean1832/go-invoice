# Email Setup: Google App Password

This guide explains how to set up email sending using a Google App Password. This method is simpler than OAuth2 but less secure because it uses a static password.

> [!WARNING]
> App Passwords provide full access to your Gmail account. If your `.env` file or server is compromised, attackers can access your email. Consider using [OAuth2 authentication](./email-setup-oauth2.md) for better security.

## Prerequisites

- A Google Account with Gmail
- 2-Step Verification enabled on your Google Account

## Step 1: Enable 2-Step Verification

App Passwords are only available if you have 2-Step Verification enabled.

1. Go to your [Google Account Security settings](https://myaccount.google.com/security)
2. Under "How you sign in to Google", click on **2-Step Verification**
3. Follow the prompts to enable 2-Step Verification if it's not already enabled
4. You'll need to verify your phone number and set up a second authentication method

## Step 2: Create an App Password

1. Go to [Google App Passwords](https://myaccount.google.com/apppasswords)
2. You may be prompted to sign in again
3. Under "App name", enter a descriptive name (e.g., "go-invoice")
4. Click **Create**
5. Google will display a 16-character password (it looks like `xxxx xxxx xxxx xxxx`)
6. **Copy this password immediately** - you won't be able to see it again

> [!IMPORTANT]
> Store the App Password securely. If you lose it, you'll need to delete it and create a new one.

## Step 3: Configure Your `.env` File

Create or edit the `.env` file in the **same directory as your `go-invoice` executable**:

```env
# SMTP Settings for Gmail
SMTP_FROM=your-email@gmail.com
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587

# App Password (remove spaces from the 16-character password Google provided)
SMTP_PASSWORD=xxxxxxxxxxxxxxxx
```

**Important notes:**

- Remove all spaces from the App Password (Google shows it as `xxxx xxxx xxxx xxxx`, but you should enter it as `xxxxxxxxxxxxxxxx`)
- Replace `your-email@gmail.com` with your actual Gmail address
- Do not use your regular Gmail password - it won't work

### Example `.env` File Location

Your directory structure should look like this:

```
your-app-folder/
├── go-invoice          # or go-invoice.exe on Windows
├── .env                # Your configuration file (create this)
└── db/                 # Auto-created on first run
    ├── clients/
    ├── invoices/
    ├── providers/
    └── email_templates/
```

## Step 4: Test Email Sending

1. Restart the `go-invoice` application to load the new configuration
2. Create or open an invoice in the web interface
3. Click the "Send Email" button
4. If configured correctly, the email should send successfully

## Troubleshooting

### "Invalid credentials" error

- Verify you're using the App Password, not your regular Gmail password
- Check that there are no spaces in the `SMTP_PASSWORD` value
- Ensure 2-Step Verification is still enabled on your Google Account

### "Authentication failed" error

- Verify `SMTP_FROM` matches your Gmail address exactly
- Check that `SMTP_HOST=smtp.gmail.com` and `SMTP_PORT=587`
- Try creating a new App Password and updating `.env`

### Email not sending

- Check the application logs for error messages
- Verify your `.env` file is in the same directory as the `go-invoice` executable
- Restart the application after making changes to `.env`

### App Password option not available

- Ensure 2-Step Verification is enabled
- You cannot create App Passwords if your account uses security keys only
- Advanced Protection users cannot use App Passwords

## Managing App Passwords

To view or revoke App Passwords:

1. Go to [Google App Passwords](https://myaccount.google.com/apppasswords)
2. You'll see a list of all App Passwords you've created
3. Click the delete icon (❌) next to any password to revoke it

> [!TIP]
> Rotate your App Password periodically for better security. Create a new one, update your `.env` file, restart the app, then delete the old App Password.

## Upgrading to OAuth2

For better security, consider upgrading to [OAuth2 authentication](./email-setup-oauth2.md). OAuth2 provides:

- Temporary access tokens instead of permanent passwords
- Ability to revoke access from your Google Account settings
- No credentials stored in your `.env` file
- Automatic token refresh

See the [OAuth2 setup guide](./email-setup-oauth2.md) for instructions.
