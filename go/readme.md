## Hello JWT(GO)
- fork from the article: https://www.nexmo.com/blog/2020/03/13/using-jwt-for-authentication-in-a-golang-application-dr

```sh
go run main.go
```

```sh
▶ curl :8080/login -d'{
"username":"username",
"password":"password"
}'

TOKEN=...
curl :8080/todo -H 'authorization:bearer $TOKEN' -d'{
"title":"My first todo"
}'
{"user_id":1,"title":"My first todo"}
```

