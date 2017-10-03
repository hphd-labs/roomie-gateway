# roomie-gateway
The API gateway service for the roomie app

## auth
`/auth` routes are reverse proxied to `$AUTH_UPSTREAM`

## config
`PORT` port to listen on  
`AUTH_UPSTREAM` upstream gateway for `/auth` routes