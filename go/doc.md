# Using JWT for Authentication in a Golang Application

fork from https://www.nexmo.com/blog/2020/03/13/using-jwt-for-authentication-in-a-golang-application-dr

## Introduction

A JSON Web Token (JWT) is a compact and self-contained way for securely transmitting information between parties as a JSON object, and they are commonly used by developers in their APIs. JWTs are popular because:

1. A JWT is stateless. That is, it does not need to be stored in a database (persistence layer), unlike opaque tokens.
1. The signature of a JWT is never decoded once formed, thereby ensuring that the token is safe and secure.
1. A JWT can be set to be invalid after a certain period of time. This helps minimize or totally eliminate any damage that can be done by a hacker, in the event that the token is hijacked.


## What Makes Up a JWT?
A JWT is comprised of three parts:

- **Header**: the type of token and the signing algorithm used.
The type of token can be “JWT” while the Signing Algorithm can either be HMAC or SHA256.
- **Payload**: the second part of the token which contains the claims. These claims include application specific data(e.g, user id, username), token expiration time(exp), issuer(iss), subject(sub), and so on.
- **Signature**: the encoded header, encoded payload, and a secret you provide are used to create the signature.
Let’s use a simple token to understand the above concepts.


```sh
Token = eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdXRoX3V1aWQiOiIxZGQ5MDEwYy00MzI4LTRmZjMtYjllNi05NDRkODQ4ZTkzNzUiLCJhdXRob3JpemVkIjp0cnVlLCJ1c2VyX2lkIjo3fQ.Qy8l-9GUFsXQm4jqgswAYTAX9F4cngrl28WJVYNDwtM
```

You can navigate to **https://jwt.io/** and test the token signature if it is verified or not. Use “HS512” as the algorithm. You will get the message “Signature Verified”.

To make the signature, your application will need to provide a key. This key enables the signature to remain secure—even when the JWT is decoded the signature remains encrypted. It is highly recommended to always use a secret when creating a JWT.

## Token Types
Since a JWT can be set to expire (be invalidated) after a particular period of time, two tokens will be considered:

- **Access Token**: An access token is used for requests that require authentication. It is normally added in the header of the request. It is recommended that an access token have a short lifespan, say ***15 minutes***. Giving an access token a short time span can prevent any serious damage if a user’s token is tampered with, in the event that the token is hijacked. The hacker only has 15 minutes or less to carry out his operations before the token is invalidated.
- **Refresh Token**: A refresh token has a longer lifespan, usually ***7 days***. This token is used to generate new access and refresh tokens. In the event that the access token expires, new sets of access and refresh tokens are created when the refresh token route is hit (from our application).

## Where to Store a JWT
For a production grade application, it is highly recommended to **store JWTs in an HttpOnly cookie**. To achieve this, while sending the cookie generated from the backend to the frontend (client), a HttpOnly flag is sent along the cookie, instructing the browser not to display the cookie through the client-side scripts. Doing this can prevent XSS (Cross Site Scripting) attacks.
JWT can also be stored in browser local storage or session storage. Storing a JWT this way can expose it to several attacks such as XSS mentioned above, so it is generally less secure when compared to using `HttpOnly cookie technique.

## The Application
We will consider a simple todo restful API.Create a directory called jwt-todo, then initialize go.mod for dependency management. go.mod is initialized using:

```sh
go mod init jwt-todo
```

You can install **gin**, if you have not already, using:

```sh
go get github.com/gin-gonic
```

## Login Request
When a user’s details have been verified, they are logged in and a JWT is generated on their behalf. We will achieve this in the Login() function defined below:

func Login(c *gin.Context) {
...
}

We received the user’s request, then unmarshalled it into the User struct. We then compared the input user with the one we defined we defined in memory. If we were using a database, we would have compared it with a record in the database.

So as not to make the Login function bloated, the logic to generate a JWT is handled by CreateToken. Observe that the user id is passed to this function. It is used as a claim when generating the JWT.

Let’s define the CreateToken function:

func CreateToken(userid uint64) (string, error) {
...
  return token, nil
}

We set the token to be valid only for 15 minutes, after which, it is invalid and cannot be used for any authenticated request. Also observe that we signed the JWT using a secret(ACCESS_SECRET) obtained from our environmental variable. It is highly recommended that this secret is not exposed in your codebase, but rather called from the environment just like we did above. You can save it in a .env, .yml or whatever works for you.

We can now run the application:

```sh
go run main.go

▶ curl -s :8080/login -d'{
"username":"username",
"password":"password"
}'| jq
```


As seen above, we have generated a JWT that will last for 15 minutes.

## Implementation Loopholes
Yes we can login a user a generate a JWT, but there is a lot wrong with the above implementation:

1. The JWT can only be invalidated when it expires. A major limitation to this is: a user can login, then decide to **logout** immediately, but the user’s JWT remains valid until the expiration time is reached.
1. The JWT might be hijacked and used by a hacker without the user doing anything about it, until the token expires.
1. The user will need to **relogin** after the token expires, thereby leading to poor user experience.

We can address the problems stated above in two ways:

1. Using a persistence storage layer to **store JWT metadata**. This will enable us to invalidate a JWT the very second a the user logs out, thereby improving security.
1. Using the concept of **refresh token** to generate a new access token, in the event that the access token expired, thereby improving the user experience.

## Using Redis to Store JWT Metadata
One of the solutions we proffered above is saving a JWT metadata in a persistence layer. This can be done in any persistence layer of choice, but redis is highly recommended. Since the JWTs we generate have expiry time, redis has a feature that automatically deletes data whose expiration time has reached. Redis can also handle a lot of writes and can scale horizontally.

Since redis is a key-value storage, its keys need to be unique, to achieve this, we will use uuid as the key and use the user id as the value.

So let’s install two packages to use:

```
go get github.com/go-redis/redis/v7
go get github.com/twinj/uuid
```

Note: It is expected that you have redis installed in your local machine. If not, you can pause and do that, before continuing.

```
docker run --name my-redis -p 6379:6379 --restart always --detach redis
```

The redis client is initialized in the init() function. This ensures that each time we run, redis is automatically connected.

When we create a token from this point forward, we will generate a uuid that will be used as one of the token claims, just as we used the user id as a claim in the previous implementation.

## Define the Metadata
In our proposed solution, instead of just creating one token, we will need to create two JWTs:

1. The Access Token
1. The Refresh Token

To achieve this, we will need to define a struct that house these tokens definitions, their expiration periods and uuids:
...

The expiration period and the uuids are very handy because they will be used when saving token metadata in redis.

Now, let’s update the CreateToken function to look like this: 
...

In the above function, the Access Token expires after 15 minutes and the Refresh Token expires after 7 days. You can also observe we added a uuid as a claim to each token.
Since the uuid is unique each time it is created, a user can create more than one token. This happens when a user is logged in on different devices. The user can also logout from any of the devices without them being logged out from all devices. How cool!

## Saving JWTs metadata
Let’s now wire up the function that will be used to save the JWTs metadata:

func CreateAuth(userid uint64, td *TokenDetails) error {
   ...
    return nil
}

We passed in the TokenDetails which have information about the expiration time of the JWTs and the uuids used when creating the JWTs. If the expiration time is reached for either the refresh token or the access token, the JWT is automatically deleted from redis.


Excellent! We have both the access_token and the refresh_token, and also have token metadata persisted in redis.

## Creating a RESTful api
We can now proceed to make requests that require authentication using JWT.

One of the unauthenticated requests in this API is the creation of todo request.

When performing any authenticated request, we need to validate the token passed in the authentication header to see if it is valid. We need to define some helper functions that help with these.

First we will need to extract the token from the request header using the ExtractToken function:

```
func ExtractToken(r *http.Request) string {
  ...
  return ""
}
```

Then we will verify the token:

```
func VerifyToken(r *http.Request) (*jwt.Token, error) {
  ...
  return token, nil
}
```

We called ExtractToken inside the VerifyToken function to get the token string, then proceeded to check the signing method.

Then, we will check the validity of this token, whether it is still useful or it has expired, using the TokenValid function:

```
func TokenValid(r *http.Request) error {
  ...
  return nil
}
```

We will also extract the token metadata that will lookup in our redis store we set up earlier. To extract the token, we define the ExtractTokenMetadata function:

```
func ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
  ...
  return nil, err
}
```

The ExtractTokenMetadata function returns an AccessDetails (which is a struct). This struct contains the metadata (access_uuid and user_id) that we will need to make a lookup in redis. If there is any reason we could not get the metadata from this token, the request is halted with an error message.

The AccessDetails struct mentioned above looks like this:

We also mentioned looking up the token metadata in redis. Let’s define a function that will enable us to do that:

```
func FetchAuth(authD *AccessDetails) (uint64, error) {
  ...
}
```

FetchAuth() accepts the AccessDetails from the ExtractTokenMetadata function, then looks it up in redis. If the record is not found, it may mean the token has expired, hence an error is thrown.

Let’s finally wire up the CreateTodo function to better understand the implementation of the above functions:

```
func CreateTodo(c *gin.Context) {
...
  c.JSON(http.StatusCreated, td)
}
```

As seen, we called the ExtractTokenMetadata to extract the JWT metadata which is used in FetchAuth to check if the metadata still exists in our redis store. If everything is good, the Todo can then be saved to the database, but we chose to return it to the caller.

Let’s update main() to include the CreateTodo function:

```
func main() {
  router.POST("/login", Login)
  router.POST("/todo", CreateTodo)

  log.Fatal(router.Run(":8080"))
}
```

To test CreateTodo, login and copy the access_token and add it to the Authorization Bearer Token field like this:

```
TOKEN=...
curl :8080/todo -H 'authorization:bearer $TOKEN' -d'{
"title":"My first todo"
}'
{"user_id":1,"title":"My first todo"}
```


Then add a title to the request body to create a todo and make a POST request to the /todo endpoint, as shown below:

Attempting to create a todo without an access_token will be unauthorized:

## Logout Request
Thus far, we have seen how a JWT is used to make an authenticated request. When a user logs out, we will instantly revoke/invalidate their JWT. This is achieved by deleting the JWT metadata from our redis store.

We will now define a function that enables us delete a JWT metadata from redis:

```
func DeleteAuth(givenUuid string) (int64,error) {
  deleted, err := client.Del(givenUuid).Result()
  if err != nil {
     return 0, err
  }
  return deleted, nil
}
```

The function above will delete the record in redis that corresponds with the uuid passed as a parameter.

The Logout function looks like this:

```
func Logout(c *gin.Context) {
  ...
  c.JSON(http.StatusOK, "Successfully logged out")
}
```

In the Logout function, we first extracted the JWT metadata. If successful, we then proceed with deleting that metadata, thereby rendering the JWT invalid immediately.

Before testing, update the main.go file to include the logout endpoint like this:
...

Provide a valid access_token associated with a user, then logout the user. Remember to add the access_token to the Authorization Bearer Token, then hit the logout endpoint:

Now the user is logged out, and no further request can be performed with that JWT again as it is immediately invalidated. This implementation is more secure than waiting for a JWT to expire after a user logs out.

## Securing Authenticated Routes
We have two routes that require authentication: /login and /logout. Right now, with or without authentication, anybody can access these routes. Let’s change that.

We will need to define the TokenAuthMiddleware() function to secure these routes:

```
func TokenAuthMiddleware() gin.HandlerFunc {
  return func(c *gin.Context) {
     err := TokenValid(c.Request)
     if err != nil {
        c.JSON(http.StatusUnauthorized, err.Error())
        c.Abort()
        return
     }
     c.Next()
  }
}
```

As seen above, we called the TokenValid() function (defined earlier) to check if the token is still valid or has expired. The function will be used in the authenticated routes to secure them.
Let’s now update main.go to include this middleware:

```
func main() {
  router.POST("/login", Login)
  router.POST("/todo", TokenAuthMiddleware(), CreateTodo)
  router.POST("/logout", TokenAuthMiddleware(), Logout)

  log.Fatal(router.Run(":8080"))
}
```

## Refreshing Tokens
Thus far, we can create, use and revoke JWTs. In an application that will involve a user interface, what happens if the access token expires and the user needs to make an authenticated request? Will the user be unauthorized, and be made to login again? Unfortunately, that will be the case. But this can be averted using the concept of a refresh token. The user does not need to relogin.
The refresh token created alongside the access token will be used to create new pairs of access and refresh tokens.

Using JavaScript to consume our API endpoints, we can refresh the JWTs like a breeze using https://github.com/axios/axios. In our API, we will need to send a POST request with a refresh_token as the body to the /refresh endpoint.

Let’s first create the Refresh() function:

```
func Refresh(c *gin.Context) {
 ...
     c.JSON(http.StatusUnauthorized, "refresh expired")
  }
}
```

While a lot is going on in that function, let’s try and understand the flow.

- We first took the refresh_token from the request body.
- We then verified the signing method of the token.
- Next, check if the token is still valid.
- The refresh_uuid and the user_id are then extracted, which are metadata used as claims when creating the refresh token.
- We then search for the metadata in redis store and delete it using the refresh_uuid as key.
- We then create a new pair of access and refresh tokens that will now be used for future requests.
- The metadata of the access and refresh tokens are saved in redis.
- The created tokens are returned to the caller.
- In the else statement, if the refresh token is not valid, the user will not be allowed to create a new pair of tokens. We will need to relogin to get new tokens.

Next, add the refresh token route in the main() function:

```
  router.POST("/refresh", Refresh)
```

Testing the endpoint with a valid refresh_token:

```
curl :8080/refresh
```

And we have successfully created new token pairs. Great😎.

## Conclusion
You have seen how you can create and invalidate a JWT. You also saw how you can integrate the Nexmo Messages API in your Golang application to send notifications. You can extend this application and use a real database to persist users and todos, and you can also use a React or VueJS to build a frontend. That is where you will really appreciate the Refresh Token feature with the help of Axios Interceptors.