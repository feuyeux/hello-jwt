## Hello JWT
- fork from the article: https://www.javainuse.com/spring/boot-jwt
- support webflux
- upgrade springboot to 2.3.0 to support rsocket-security

#### Spring Security and JWT Configuration
We will be configuring Spring Security and JWT for performing 2 operations
- Generating JWT - Expose a POST API with mapping /authenticate. On passing correct username and password it will generate a JSON Web Token(JWT)
- Validating JWT - If user tries to access GET API with mapping /hello. It will allow access only if request has a valid JSON Web Token(JWT)
![](doc/1.jpeg)

#### Generating JWT
![](doc/2.jpeg)
![](doc/3.jpeg)

#### Validating JWT
![](doc/4.jpeg)

##### jwt.secret
The secret key is combined with the header and the payload to create a unique hash. We are only able to verify this hash if you have the secret key.

src/main/resources/application.properties
```properties
jwt.secret=hello_man
```

##### JwtTokenUtil 
src/main/java/org/feuyeux/jwt/config/JwtTokenUtil.java

The JwtTokenUtil is responsible for performing JWT operations like creation and validation.It makes use of the io.jsonwebtoken.Jwts for achieving this.

##### JWTUserDetailsService
JWTUserDetailsService implements the Spring Security UserDetailsService interface. It overrides the loadUserByUsername for fetching user details from the database using the username. The Spring Security Authentication Manager calls this method for getting the user details from the database when authenticating the user details provided by the user. Here we are getting the user details from a hardcoded User List.

src/main/java/org/feuyeux/jwt/service/JwtUserDetailsService.java
 
##### JwtAuthenticationController
Expose a POST API /authenticate using the JwtAuthenticationController. The POST API gets username and password in the body- Using Spring Authentication Manager we authenticate the username and password.If the credentials are valid, a JWT token is created using the JWTTokenUtil and provided to the client.

src/main/java/org/feuyeux/jwt/controller/JwtAuthenticationController.java

##### JwtRequest/JwtResponse
- src/main/java/org/feuyeux/jwt/model/JwtRequest.java
- src/main/java/org/feuyeux/jwt/model/JwtResponse.java

##### JwtRequestFilter
The JwtRequestFilter extends the Spring Web Filter OncePerRequestFilter class. For any incoming request this Filter class gets executed. It checks if the request has a valid JWT token. If it has a valid JWT Token then it sets the Authentication in the context, to specify that the current user is authenticated.

- src/main/java/org/feuyeux/jwt/config/JwtRequestFilter.java

##### JwtAuthenticationEntryPoint
This class will extend Spring's AuthenticationEntryPoint class and override its method commence. It rejects every unauthenticated request and send error code 401

src/main/java/org/feuyeux/jwt/config/JwtAuthenticationEntryPoint.java

##### WebSecurityConfig
This class extends the WebSecurityConfigurerAdapter is a convenience class that allows customization to both WebSecurity and HttpSecurity.

src/main/java/org/feuyeux/jwt/config/WebSecurityConfig.java

#### Testing
```sh
▶ mvn spring-boot:run
```

```sh
▶ ./curl_test.sh
Generate a JSON Web Token:
TOKEN=eyJhbGciOiJIUzUxMiJ9.eyJzdWIiOiJoZWxsb19tYW4iLCJleHAiOjE1ODY3OTg0MDIsImlhdCI6MTU4Njc4MDQwMn0.yXTZN8SD4JCXdi1vBlhPJRgiDHmL2lRyAikrZ3KrM08t3D4-PXCldW89Z9YuT4bhLnYezLykuGxc_l65DXroCg

Validate the JSON Web Token(Web):
localhost:8080/hello:
Hello World

Validate the JSON Web Token(Webflux):
localhost:8080/hello-stream:
data:0

data:1

data:4
```