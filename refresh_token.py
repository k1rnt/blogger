from google_auth_oauthlib.flow import InstalledAppFlow
flow = InstalledAppFlow.from_client_secrets_file(
    "client_secret.json",
    scopes=["https://www.googleapis.com/auth/blogger"]
)
creds = flow.run_local_server(port=0)
print("refresh_token =", creds.refresh_token)