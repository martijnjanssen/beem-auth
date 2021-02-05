# BEEM AUTH

Install dependencies
```sh
make dependencies
```

Generate messages and services from proto:
```sh
make generate
```

Start the database
``` sh
docker run --name beem-auth-postgres -p 5432:5432 -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=beem-auth -d 
```

For testing
``` sh
make test
```

Start with docker compose

``` sh
docker-compose up
```


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
