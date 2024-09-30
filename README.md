# ovh-oauth2cli

SSO Client for OVH

## Installation

```
go install github.com/bmassemin/ovh-oauth2cli@latest
```

## OVH Configuration

To use this tool, you first need to create an OAuth2 client through the OVH API.
https://eu.api.ovh.com/console/?section=%2Fme&branch=v1#post-/me/api/oauth2/client
The payload should look like this:
```json
{
    "callbackUrls": [
        "https://localhost:8080"
    ],
    "description": "description",
    "flow": "AUTHORIZATION_CODE",
    "name": "name"
}
```
⚠️ `callbackUrls` MUST at least contain `https://localhost:8080` (The port can be changed)\
⚠️ `flow` MUST be set to `AUTHORIZATION_CODE`

The response will provide you with a clientId and a clientSecret, which will be required each time you need to log in.

## Generate a `ovh.conf` file with an access token

Simply run the following command: 
```
ovh-oauth2cli -client-id <client_id> -client-secret <client_secret>
```

### Automatically open the browser with WSL

- Install https://github.com/wslutilities/wslu (already installed on Ubuntu)
- Add the following line to your .bashrc or .zshrc file:
```
export BROWSER=explorer.exe
```

### Generate your own SSL certificates

In the `cert` folder, you can generate your own certificates, and add them to your certificate store.\
Then run `ovh-oauth2cli` with additionnal parameters:
```
./ovh-oauth2cli -client-id <client_id> -client-secret <client_secret> -server-crt <path_to_server_crt> -server-key <path_to_server_key>
```