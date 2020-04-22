package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/feuyeux/hello-jwt/env"
	"github.com/feuyeux/hello-jwt/store"
	"github.com/feuyeux/hello-jwt/stru"
	"github.com/gin-gonic/gin"
	"github.com/twinj/uuid"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	router = gin.Default()
	user   = stru.User{
		ID:       1,
		Username: "username",
		Password: "password",
	}
)

func main() {
	router.POST("/login", Login)
	router.POST("/api", API)
	router.POST("/logout", Logout)
	router.POST("/refresh", Refresh)
	log.Fatal(router.Run(":8080"))
}

func Login(c *gin.Context) {
	var u stru.User
	log.Print("[Login] 1 Receive and unmarshall the user’s request into the User struct")
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}
	//compare the user from the request, with the one we defined:

	log.Print("[Login] 2 Authenticate")
	if user.Username != u.Username || user.Password != u.Password {
		c.JSON(http.StatusUnauthorized, "Please provide valid login details")
		return
	}
	log.Print("[Login] 3 Create AT&RT")
	ts, err := signToken(user.ID)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	log.Print("[Login] 4 Store AT&RT")
	saveErr := storeToken(user.ID, ts)
	if saveErr != nil {
		c.JSON(http.StatusUnprocessableEntity, saveErr.Error())
	}
	log.Print("[Login] 5 Return AT&RT")
	tokens := map[string]string{
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
	}
	c.JSON(http.StatusOK, tokens)
}
func API(c *gin.Context) {
	log.Print("[API] 1 Verify Token")
	token, err := verifyToken(c.Request)
	log.Print("[API] 2 ParseTo Auth")
	tokenAuth, err := parseAuth(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	log.Print("[API] 3 Authorize")
	userId, err := authorize(tokenAuth)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	log.Print("[API] 4 Resource")
	var td *stru.ResSt
	if err := c.ShouldBindJSON(&td); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "invalid json")
		return
	}
	td.UserID = userId
	c.JSON(http.StatusCreated, td)
}
func Logout(c *gin.Context) {
	token, _ := verifyToken(c.Request)
	au, err := parseAuth(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	deleted, delErr := revokeAuth(au.AccessUuid)
	if delErr != nil || deleted == 0 { //if any goes wrong
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	c.JSON(http.StatusOK, "Successfully logged out")
}

func Refresh(c *gin.Context) {
	mapToken := map[string]string{}
	if err := c.ShouldBindJSON(&mapToken); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	refreshToken := mapToken["refresh_token"]

	log.Print("[Refresh] 1 Verify RT: refreshToken=", refreshToken)
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(env.Get("REFRESH_SECRET")), nil
	})
	//if there is an error, the token must have expired
	if err != nil {
		c.JSON(http.StatusUnauthorized, "Refresh token expired")
		return
	}
	//is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		c.JSON(http.StatusUnauthorized, err)
		return
	}
	log.Print("[Refresh] 2 Parse Auth")
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			c.JSON(http.StatusUnprocessableEntity, err)
			return
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, "Error occurred")
			return
		}
		log.Print("[Refresh] 3 Revoke Refresh Token")
		deleted, delErr := revokeAuth(refreshUuid)
		if delErr != nil || deleted == 0 { //if any goes wrong
			c.JSON(http.StatusUnauthorized, "unauthorized")
			return
		}
		log.Print("[Refresh] 4 Create AT&RT")
		ts, createErr := signToken(userId)
		if createErr != nil {
			c.JSON(http.StatusForbidden, createErr.Error())
			return
		}
		log.Print("[Refresh] 5 Store AT&RT")
		saveErr := storeToken(userId, ts)
		if saveErr != nil {
			c.JSON(http.StatusForbidden, saveErr.Error())
			return
		}
		log.Print("[Login] 6 Return AT&RT")
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}
		c.JSON(http.StatusCreated, tokens)
	} else {
		c.JSON(http.StatusUnauthorized, "refresh expired")
	}
}
func signToken(userid uint64) (*stru.TokenDetails, error) {
	td := &stru.TokenDetails{}
	log.Print("[signToken] 1.1 Set the AccessToken Expires/to be valid for 15 minutes")
	td.AtExpires = time.Now().Add(time.Minute * 15).Unix()
	log.Print("[signToken] 1.2 Set the RefreshToken Expires/to be valid for 7 days")
	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	log.Print("[signToken] 1.3 Generate UUID for AccessToken and RefreshToken")
	td.AccessUuid = uuid.NewV4().String()
	td.RefreshUuid = uuid.NewV4().String()

	var err error
	log.Print("[signToken] 2 Sign Access Token Claims with HS256")
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = userid
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(env.Get("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	log.Print("[signToken] 3 Sign Refresh Token Claims with HS256")
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = userid
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(env.Get("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	return td, nil
}

func storeToken(userid uint64, td *stru.TokenDetails) error {
	//converting Unix to UTC(to Time object)
	at := time.Unix(td.AtExpires, 0)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	atExpired := at.Sub(now)
	log.Print("[storeToken] 1 Save Access Token to redis")
	errAccess := store.Client.Set(td.AccessUuid, strconv.Itoa(int(userid)), atExpired).Err()
	if errAccess != nil {
		return errAccess
	}
	rtExpired := rt.Sub(now)
	log.Print("[storeToken] 2 Save Refresh Token to redis")
	errRefresh := store.Client.Set(td.RefreshUuid, strconv.Itoa(int(userid)), rtExpired).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

func parseAuth(token *jwt.Token) (*stru.AccessDetails, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, nil
		}
		userId, err := strconv.ParseUint(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, nil
		}
		return &stru.AccessDetails{
			AccessUuid: accessUuid,
			UserId:     userId,
		}, nil
	}
	return nil, nil
}

func verifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := extractToken(r)
	log.Print("[verifyToken] extractToken:", tokenString)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(env.Get("ACCESS_SECRET")), nil
	})
	log.Print("[verifyToken] token=", token)
	if err != nil {
		return nil, err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return nil, err
	}
	return token, nil
}

func extractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}
func authorize(authD *stru.AccessDetails) (uint64, error) {
	userid, err := store.Client.Get(authD.AccessUuid).Result()
	if err != nil {
		return 0, err
	}
	userID, _ := strconv.ParseUint(userid, 10, 64)
	return userID, nil
}

func revokeAuth(givenUuid string) (int64, error) {
	deleted, err := store.Client.Del(givenUuid).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}
