# BEEM AUTH

## First start

Install dependencies
```sh
make dependencies
```

Generate messages and services from proto:
```sh
make generate
```

## Testing

For testing
``` sh
make test
```

For testing with coverage

``` sh
make cover
make coverreport
```

## Running
Copy the `dev.env.sample` to `dev.env` and fill in the required values.
Start with docker compose

``` sh
docker-compose up
```

After code changes remember to run `docker-compose build`, otherwise code changes are not visible.

# Making requests
For testing purposes, there is a `make grpcui` command that is available which allows you to make postman-esque calls to the endpoints the server exposes.


# email/password account creation

1. create email,password,accountvalid=false in db

2. send confirmation email to email

​

# email validation

validate jwt token or random string from the email callback

just use a mailing service like sendinblue/sendgrid, 100 emails per day is plenty for now, and testing with the testing apikey is unlimited as far as I know.

​

# password reset

send password reset email with token that expires

​

# password reset

takes new password+token that expires

​

# oauth/openid login

​

### Sidenote

Token validity could be skipped for now if we just use 30 char random strings. We could keep them in memory for 1 day and then remove them from memory with a cron job?
