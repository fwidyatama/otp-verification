# OTP-Verification

> This is mini project that used SMS OTP from vonage. Developed using Golang using MongoDB as database and Vonage for SMS verification.

## Usage

Run the server in project root

```
go run main.go
```

Make sure you already have vonage account, you can register [here](https://www.vonage.id/). You will get €2 credits to use vonage services.

Don't forget to change .env.example into .env

| KEY        | Value             |
| ---------- | ----------------- |
| API_KEY    | VONAGE_API_KEY    |
| API_SECRET | VONAGE_API_SECRET |

## List of endpoints

| URL                                      | Method | Description                                                      |
| ---------------------------------------- | ------ | ---------------------------------------------------------------- |
| /api/users                               | POST   | Register a new user                                              |
| /api/users                               | GET    | Show all user                                                    |
| /api/users?verified=value&position=value | GET    | Show all filtered user. Can choose only 1 filter or combine them |
| /api/users/verify                        | POST   | Verify user using otp                                            |
