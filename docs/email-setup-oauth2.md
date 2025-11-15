# Email Setup: Google OAuth2

This guide explains how to set up email sending using Google OAuth2 authentication. This method is more secure than App Passwords because it uses temporary access tokens and doesn't require storing your password.

## Why Use OAuth2?

- **Better Security**: No passwords stored in configuration files
- **Granular Permissions**: Grant only email sending access
- **Easy Revocation**: Revoke access anytime from your Google Account
- **Automatic Refresh**: Access tokens refresh automatically
- **Audit Trail**: See which apps have access to your account

## Prerequisites

- A Google Account with Gmail
- Access to [Google Cloud Console](https://console.cloud.google.com/)
- Your `go-invoice` deployment URL (e.g., `http://localhost:8080` or `https://yourdomain.com`)

## Overview

Setting up OAuth2 involves:

1. Creating a Google Cloud project
2. Configuring the OAuth consent screen
3. Enabling the Gmail API
4. Creating OAuth2 credentials
5. Configuring redirect URIs
6. Adding credentials to your `.env` file
7. Authenticating in the app

## Step 1: Create a Google Cloud Project

1. Go to [Google Cloud Console](https://console.cloud.google.com/)
2. Click the project dropdown at the top and select **New Project**
3. Enter a project name (e.g., "go-invoice")
4. Click **Create**
5. Wait for the project to be created, then select it from the project dropdown

## Step 2: Configure OAuth Consent Screen

Before creating credentials, you must configure the OAuth consent screen.

1. In the Google Cloud Console, navigate to **APIs & Services** > **OAuth consent screen**
2. Select **External** as the user type (unless you have a Google Workspace account)
3. Click **Create**

### Fill in the required information:

**App information:**

- **App name**: `go-invoice` (or your preferred name)
- **User support email**: Your email address
- **App logo**: (Optional) Upload a logo

**App domain** (Optional for testing):

- **Application home page**: Your app URL (e.g., `http://localhost:8080`)
- **Application privacy policy link**: (Leave blank for testing)
- **Application terms of service link**: (Leave blank for testing)

**Developer contact information:**

- **Email addresses**: Your email address

4. Click **Save and Continue**

### Add Scopes:

1. Click **Add or Remove Scopes**
2. Filter or search for these scopes and select them:
   - `https://mail.google.com/` - Full Gmail access (required for sending emails)
   - `.../auth/userinfo.email` - See your email address (automatically added)
   - `.../auth/userinfo.profile` - See your personal info (automatically added)
3. Click **Update**
4. Click **Save and Continue**

### Test Users (for External apps in Testing mode):

1. Click **Add Users**
2. Add your Gmail address (and any other users who should have access)
3. Click **Add**
4. Click **Save and Continue**
5. Click **Back to Dashboard**

> [!NOTE]
> Apps in "Testing" mode can have up to 100 test users. For production use with unlimited users, you'll need to publish your app (requires Google verification for sensitive scopes like Gmail).

## Step 3: Enable Gmail API

1. In Google Cloud Console, go to **APIs & Services** > **Library**
2. Search for "Gmail API"
3. Click on **Gmail API**
4. Click **Enable**
5. Wait for the API to be enabled

## Step 4: Create OAuth2 Credentials

1. Go to **APIs & Services** > **Credentials**
2. Click **Create Credentials** at the top
3. Select **OAuth client ID**
4. If prompted to configure the consent screen, follow Step 2 above

### Configure the OAuth client:

1. **Application type**: Select **Web application**
2. **Name**: Enter a name (e.g., "go-invoice web client")

### Add Authorized Redirect URIs:

This is critical - the redirect URI must match exactly where your app is deployed.

**For local development:**

```
http://localhost:8080/api/v1/mailer/auth/google/callback
```

**For production (replace with your domain):**

```
https://yourdomain.com/api/v1/mailer/auth/google/callback
```

**For custom ports (e.g., port 3000):**

```
http://localhost:3000/api/v1/mailer/auth/google/callback
```

> [!IMPORTANT]
> The redirect URI format is always:
>
> ```
> {PUBLIC_URL}/api/v1/mailer/auth/google/callback
> ```
>
> Where `PUBLIC_URL` is your server's public URL including protocol and port (if not 80/443).

### Add multiple redirect URIs (recommended):

You can add multiple URIs for different environments:

1. Click **+ Add URI**
2. Add localhost URI for development
3. Click **+ Add URI** again
4. Add production URI

Example:

```
http://localhost:8080/api/v1/mailer/auth/google/callback
https://invoice.yourdomain.com/api/v1/mailer/auth/google/callback
```

5. Click **Create**

### Save Your Credentials:

After clicking Create, a dialog will appear with your credentials:

- **Client ID**: Looks like `123456789-abcdefghijklmnop.apps.googleusercontent.com`
- **Client Secret**: Looks like `GOCSPX-abc123def456ghi789jkl`

> [!WARNING]
> **Copy both values immediately!** You can retrieve the Client ID later, but the Client Secret is only shown once. If you lose it, you'll need to create a new client secret.

## Step 5: Configure Your `.env` File

Create or edit the `.env` file in the **same directory as your `go-invoice` executable**:

```env
# SMTP Settings (host and port are still required for OAuth2)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587

# Google OAuth2 Credentials
GOOGLE_OAUTH_CLIENT_ID=123456789-abcdefghijklmnop.apps.googleusercontent.com
GOOGLE_OAUTH_CLIENT_SECRET=GOCSPX-abc123def456ghi789jkl

# Public URL (must match the redirect URI you configured)
PUBLIC_URL=http://localhost:8080

# For production, set these:
# PUBLIC_URL=https://yourdomain.com
# IS_PROD=true
```

**Important notes:**

- `PUBLIC_URL` must match the redirect URI exactly (including protocol and port)
- Do not add `SMTP_PASSWORD` when using OAuth2
- For production, use HTTPS and set `IS_PROD=true`

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

## Step 6: Authenticate in the Application

1. **Restart** the `go-invoice` application to load the new configuration
2. Open your browser and go to your app (e.g., `http://localhost:8080`)
3. Navigate to **Settings** in the app
4. Look for the "Email Configuration" section
5. Click **Sign in with Google**
6. A popup window will open with Google's authentication page
7. **Select your Google Account** (must be one of the test users you added)
8. **Review the permissions** - the app is requesting access to:
   - See your email address
   - See your personal info
   - Read, compose, send, and permanently delete all your email from Gmail
9. Click **Continue** (or **Allow**)
10. If you see a warning "Google hasn't verified this app", click **Continue** (this is normal for apps in Testing mode)
11. The popup will close automatically after successful authentication
12. You should see "Authenticated as: your-email@gmail.com" in the Settings page

## Step 7: Test Email Sending

1. Create or open an invoice in the web interface
2. Click the "Send Email" button
3. If configured correctly, the email should send successfully using OAuth2

## Understanding OAuth2 Scopes

The app requests these OAuth2 scopes:

| Scope                                              | Purpose                                                       |
| -------------------------------------------------- | ------------------------------------------------------------- |
| `https://mail.google.com/`                         | Full Gmail access - required to send emails via SMTP          |
| `https://www.googleapis.com/auth/userinfo.email`   | Read your email address - used to display who's authenticated |
| `https://www.googleapis.com/auth/userinfo.profile` | Read your basic profile info - used to display your name      |

## Token Refresh

OAuth2 access tokens expire after a short time (typically 1 hour). The app automatically:

1. Stores a refresh token in your session
2. Checks token expiry before sending emails
3. Requests a new access token using the refresh token
4. Updates the session with the new token

You don't need to re-authenticate unless:

- You revoke access from your Google Account
- Your session expires (default: 30 days)
- You clear your browser cookies
- The refresh token becomes invalid

## Production Deployment

For production deployments:

### 1. Use HTTPS

OAuth2 requires HTTPS for production. Update your `.env`:

```env
PUBLIC_URL=https://yourdomain.com
IS_PROD=true
```

### 2. Update Redirect URI in Google Cloud Console

1. Go to **APIs & Services** > **Credentials**
2. Click on your OAuth client
3. Add your production redirect URI:
   ```
   https://yourdomain.com/api/v1/mailer/auth/google/callback
   ```
4. Click **Save**

### 3. Publish Your App (Optional)

If you need more than 100 users or don't want the "unverified app" warning:

1. Go to **OAuth consent screen**
2. Click **Publish App**
3. Click **Confirm**

> [!WARNING]
> Publishing an app with sensitive scopes (like `https://mail.google.com/`) requires Google's verification process, which can take weeks. For personal or small business use, keeping the app in Testing mode is usually sufficient.

## Troubleshooting

### "redirect_uri_mismatch" error

This is the most common OAuth2 error. It means the redirect URI doesn't match exactly.

**Check these:**

1. `PUBLIC_URL` in `.env` matches what you configured in Google Cloud Console
2. Protocol matches (http vs https)
3. Port number is included if not 80/443 (e.g., `:8080`)
4. No trailing slash in `PUBLIC_URL`
5. The callback path is exactly `/api/v1/mailer/auth/google/callback`

**Example of correct match:**

- `.env`: `PUBLIC_URL=http://localhost:8080`
- Google Console Redirect URI: `http://localhost:8080/api/v1/mailer/auth/google/callback`

### "Access blocked: This app's request is invalid"

Your OAuth consent screen configuration is incomplete:

1. Go to **OAuth consent screen**
2. Verify all required fields are filled
3. Ensure your email is added as a test user (for Testing mode)
4. Check that required scopes are added

### "Email sending failed" with OAuth2

1. Check the browser console for error messages
2. Verify you've authenticated in Settings
3. Try re-authenticating (sign out and sign in again)
4. Check that `SMTP_HOST=smtp.gmail.com` and `SMTP_PORT=587`
5. Ensure you don't have `SMTP_PASSWORD` set (it conflicts with OAuth2)

### "Token has been expired or revoked"

1. Go to Settings in the app
2. Click **Sign out**
3. Click **Sign in with Google** again
4. Re-authenticate

If this persists:

1. Go to [Google Account Permissions](https://myaccount.google.com/permissions)
2. Find your app and remove access
3. Re-authenticate in the app

### Can't add test users

For apps in Testing mode with External user type:

1. Ensure the email addresses have Gmail accounts
2. You can add up to 100 test users
3. Users must be added before they can authenticate

### "Access_denied" error

The user declined the permission request. They need to:

1. Click "Sign in with Google" again
2. Review and accept the requested permissions

## Managing OAuth2 Access

### View Connected Apps

See which apps have access to your Google Account:

1. Go to [Google Account Permissions](https://myaccount.google.com/permissions)
2. Find "go-invoice" (or your app name)
3. Click to see details or remove access

### Revoke Access

To revoke access:

1. Go to [Google Account Permissions](https://myaccount.google.com/permissions)
2. Find your app
3. Click **Remove Access**
4. Confirm

After revoking, you'll need to re-authenticate in the app Settings.

### Rotate Client Secret

If your client secret is compromised:

1. Go to **APIs & Services** > **Credentials**
2. Click on your OAuth client
3. Click **Add Secret** under "Client secrets"
4. Copy the new secret
5. Update your `.env` file with the new secret
6. Restart the app
7. After verifying it works, delete the old secret

## Security Best Practices

1. **Never commit `.env` to version control** - Add `.env` to `.gitignore`
2. **Use HTTPS in production** - Set `IS_PROD=true` and `PUBLIC_URL=https://...`
3. **Set a persistent SESSION_SECRET** - Generate with `openssl rand -base64 32`
4. **Limit test users** - Only add email addresses that need access
5. **Rotate secrets periodically** - Create new client secrets every 6-12 months
6. **Monitor access** - Regularly check connected apps in Google Account settings
7. **Keep the app in Testing mode** - Unless you need 100+ users (avoids verification process)

## Switching from App Password to OAuth2

If you're currently using an App Password:

1. Follow all steps above to set up OAuth2
2. Remove or comment out `SMTP_PASSWORD` from `.env`:
   ```env
   # SMTP_PASSWORD=xxxxxxxxxxxxxxxx  # Disabled - using OAuth2 instead
   ```
3. Restart the application
4. Authenticate in Settings
5. Test email sending
6. Once confirmed working, revoke the App Password:
   - Go to [Google App Passwords](https://myaccount.google.com/apppasswords)
   - Delete the app password you were using

## Additional Resources

- [Google OAuth2 Documentation](https://developers.google.com/identity/protocols/oauth2)
- [Gmail API Scopes](https://developers.google.com/gmail/api/auth/scopes)
- [OAuth2 Playground](https://developers.google.com/oauthplayground/) - Test OAuth2 flows
- [Google Cloud Console](https://console.cloud.google.com/)
