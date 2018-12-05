# Sample Backend for Go (aka Golang)

This repository contains a sample backend code that demonstrates how to generate a Virgil JWT using the [Golang SDK](https://github.com/go-virgil/virgil)

## Prerequisites
- [Go aka Golang](https://golang.org/) 
- [Virgil Security GO SDK](https://github.com/go-virgil/virgil/tree/v5)

## Clone

Clone the repository from GitHub.

```
$ git clone https://github.com/VirgilSecurity/sample-backend-go.git
```
Go to virgil-crypto library directory and build it:
```
$ cd sample-backend-go/vendor/gopkg.in/virgilsecurity/virgil-crypto-go.v5
$ make
```

## Get Virgil Credentials

If you don't have an account yet, [sign up for one](https://dashboard.virgilsecurity.com/signup) using your e-mail.

To generate a JWT the following values are required:

| Variable Name                     | Description                    |
|-----------------------------------|--------------------------------|
| API_PRIVATE_KEY                  | Private key of your API key that is used to sign the JWTs. |
| API_KEY_ID               | ID of your API key. A unique string value that identifies your account in the Virgil Cloud. |
| APP_ID                   | ID of your Virgil Application. |

## Add Virgil Credentials to .env

- open the project folder
- create a `.env` file
- fill it with your account credentials (take a look at the `.env.example` file to find out how to setup your own `.env` file)
- save the `.env` file


## Install Dependencies and Run the Server
To run the server go to the server example directory and run
```
$ set -a && source '<yourenvfile>.env' && go run main.go
```
Now, use your client code to make a request to get a JWT from the sample backend that is working on http://localhost:3000.

## Specification

### /authenticate endpoint
This endpoint is an example of users authentication. It takes user `identity` and responds with unique token.

```http
POST https://localhost:3000/authenticate HTTP/1.1
Content-type: application/json;

{
    "identity": "string"
}

Response:

{
    "authToken": "string"
}
```

### /virgil-jwt endpoint
This endpoint checks whether a user is authorized by an authorization header. It takes user's `authToken`, finds related user identity and generates a `virgilToken` (which is [JSON Web Token](https://jwt.io/)) with this `identity` in a payload. Use this token to make authorized api calls to Virgil Cloud.

```http
GET https://localhost:3000/virgil-jwt HTTP/1.1
Content-type: application/json;
Authorization: Bearer <authToken>

Response:

{
    "virgilToken": "string"
}
```

## Virgil JWT Generation
To generate JWT, you need to use the `JwtGenerator` class from the SDK.

```go
import (
	virgilSDK "gopkg.in/virgil.v5/sdk"
	virgilCrypto "gopkg.in/virgilsecurity/virgil-crypto-go.v5"
)

	cryptoInstance := virgilCrypto.NewVirgilCrypto()
	cryptoPrivateKey, _ = cryptoInstance.ImportPrivateKey([]byte(os.Getenv("API_PRIVATE_KEY")), "")
	[...]
	tokenSigner := virgilCrypto.NewVirgilAccessTokenSigner()
	generator := virgilSDK.NewJwtGenerator(cryptoPrivateKey, os.Getenv("API_KEY_ID"), tokenSigner, os.Getenv("APP_ID"), time.Hour)
	jwtToken, err := generator.GenerateToken(identity, nil)

```
Then you need to provide an HTTP endpoint which will return the JWT with the user's identity as a JSON.

For more details take a look at the [main.go](main.go) file.



## License

This library is released under the [3-clause BSD License](LICENSE.md).

## Support
Our developer support team is here to help you. Find out more information on our [Help Center](https://help.virgilsecurity.com/).

You can find us on [Twitter](https://twitter.com/VirgilSecurity) or send us email support@VirgilSecurity.com.

Also, get extra help from our support team on [Slack](https://virgilsecurity.slack.com/join/shared_invite/enQtMjg4MDE4ODM3ODA4LTc2OWQwOTQ3YjNhNTQ0ZjJiZDc2NjkzYjYxNTI0YzhmNTY2ZDliMGJjYWQ5YmZiOGU5ZWEzNmJiMWZhYWVmYTM).
