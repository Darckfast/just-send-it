# ![logo](.github/images/logo.webp)

## 📨 Just send it

Small serverless project, using Vercel free tier in addition of Resend.com, to send emails

## Requirements

This project, use of the following services

* Vercel serverless platform
* Vercel KV Storage
* Resend API

### 🔑 Auth

The Vercel KV Storage is in use, to keep a record of users current live session

JWT are the default format for all sessions

To authenticate, call the `/api/auth` endpoint, passing the api_key in the body

```sh
curl --request POST \
  --url http://localhost:3000/api/auth \
  --header 'Content-Type: application/json' \
  --data '{
 "api-key": "super-secret-api"
}'
```

The normal response is a 200 status response containing a session token, this token MUST be set on every other request, in the `Authorization` header

It is also possible to delete a session, to achieve that just do a
`DELETE` request with your session in the headers

All session have a duration of 1 hour, including their `time to live` in the Vercel's KV Storage
