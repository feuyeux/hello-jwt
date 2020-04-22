## Hello JWT(GO)
- fork from the article: https://www.nexmo.com/blog/2020/03/13/using-jwt-for-authentication-in-a-golang-application-dr

<img src="jwt-workflow.png" style="width:500px;" />

### RUN & TEST
```sh
go run main.go
```

#### Login
```sh
▶ curl -s :8080/login -d'{
"username":"username",
"password":"password"
}'| jq
```

#### Using jq to assign multiple output variables
```sh
▶ AT=$(curl -s :8080/login -d'{
"username":"username",
"password":"password"
}' | jq -r ".access_token")
echo "Access Token=$AT"

▶ read RT AT < <(echo $(curl -s :8080/login -d'{
"username":"username",
"password":"password"
}' | jq -r ".refresh_token,.access_token"))
echo "Access Token=$AT"
echo "Refresh Token=$RT"
```

#### Request api
```sh
▶ curl :8080/api -H "authorization:bearer $AT" -d'{"title":"Hello JWT!"}'
```

#### Logout $ Request api
```sh
▶ curl :8080/logout -H "authorization:bearer $AT" -XPOST

▶ curl :8080/api -H "authorization:bearer $AT" -d'{"title":"Hello JWT!"}'
"unauthorized"
```
#### Refresh
```sh
▶ curl -s :8080/refresh -d'{"refresh_token":"'$RT'"}' | jq
```
#### All in one
```sh
▶ bash curl_test.sh

:::: 1 Test Login ::::
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6IjEwYzRlZGVjLTJhMDEtNDdhZi1hNWZiLWFkNjhjY2EyMGFhZiIsImF1dGhvcml6ZWQiOnRydWUsImV4cCI6MTU4NzU2MTQ5MSwidXNlcl9pZCI6MX0.sctYPeJqYIlLnhP3mZm22Mm8iL08wQeaXdTu6aqj_DY",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODgxNjUzOTEsInJlZnJlc2hfdXVpZCI6ImY3MDEzYWY0LThiYTUtNDQ0Mi04ODhhLThmODkzOGQzZDQzNSIsInVzZXJfaWQiOjF9.BLv-6vFXKiMvvDPeD1gZ_N8j3uymKBq4pxVWKUok7Ak"
}

Access Token  :eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6ImRhZWQ3ZmE2LTFhMjktNGMxYi04OGM4LTRjMTQ4MGY5YzZhYSIsImF1dGhvcml6ZWQiOnRydWUsImV4cCI6MTU4NzU2MTQ5MSwidXNlcl9pZCI6MX0.1UpXTjQu7t7yAIMC7lTIMwziO_SLKTd9DjLHJSI44jA
Refresh Token :eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODgxNjUzOTEsInJlZnJlc2hfdXVpZCI6IjUyNTRiMjViLWFjNWEtNDc0Ny05NzZmLWEwMTMxZjQ4NTAyNSIsInVzZXJfaWQiOjF9.bwjfSVf-NkifR0g2uX0_G-se0hyt2RYV28GBFhs4OVU

:::: 2 Test Request api ::::
{"user_id":1,"title":"Hello JWT!"}

:::: 3 Test Logout $ Request api ::::
"Successfully logged out"

"unauthorized"

:::: 4 Test Refresh ::::
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfdXVpZCI6ImYzNTA5NzQwLTI0YWUtNGRjOC1hOGUwLWVjZDRkODEwNzZjYSIsImF1dGhvcml6ZWQiOnRydWUsImV4cCI6MTU4NzU2MTQ5MiwidXNlcl9pZCI6MX0.9-v8Myjty7OS7umzbU9sygRYhWHHLcKiv9YXNc0dYlM",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1ODgxNjUzOTIsInJlZnJlc2hfdXVpZCI6ImU2ZWIwODk3LTM3YzQtNGJlZS04NWNjLTc5YWZlYjM5NzBmMCIsInVzZXJfaWQiOjF9.9ucOsSzmCxEb7cCbyVgIVZMIqIqdNam7iPtkfE1mj14"
}
Done!
```
