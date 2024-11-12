# ovh-oauth2cli

SSO Client for OVH

## Usage

OVH SDKs are able to read the `ovh.conf` file and fetch the credentials from it.  
It also includes the Terraform provider. For instance, you can set up the Terraform provider this way:
```hcl
provider "ovh" {
  endpoint           = "ovh-eu"
}
```

Then, run `ovh-oauth2cli` in the Terraform root directory. You should be able to execute a `terraform plan`.\
⚠️ Don't forget to add `ovh.conf` to your `.gitignore` file.

## Installation

```
go install github.com/bmassemin/ovh-oauth2cli@v0.0.4
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

## Credits

- https://github.com/int128/oauth2cli
