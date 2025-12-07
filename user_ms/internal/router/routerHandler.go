package router

import (
	"fmt"
	"net/http"
	"time"
	dbcon "userms/internal/dbCon"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserRouterHandler struct {
	localCon *dbcon.SqlCon
	opts     *Options
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type loginResponse struct {
	AccessToken         string `json:"token"`
	ExpirationTimestamp int64  `json:"expiration"`
}

func pong(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "user pong",
	})
}

func NewRouterHandler(optsFuncs ...OptsFunc) *UserRouterHandler {
	opts := defaultOptions()
	for _, optsFunc := range optsFuncs {
		optsFunc(&opts)
	}
	sqlCon, err := dbcon.NewSqlCon(&opts.DbConfig)
	if err != nil {
		fmt.Println("Error when router handler", err)
		return nil
	}
	sqlCon.InitializeSchema()
	return &UserRouterHandler{
		localCon: sqlCon,
		opts:     &opts,
	}
}

func (u *UserRouterHandler) login(c *gin.Context) {
	var loginReq loginRequest

	// bind context request to variable that expected username and password
	err := c.ShouldBindJSON(&loginReq)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	// check that username and password is matched
	userId, err := u.localCon.UserAuthentication(loginReq.Username, loginReq.Password)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	// mock user id
	loginTime := time.Now()
	expirationTime := loginTime.Add(time.Duration(u.opts.TokenPeriodTimestamp) * time.Second)

	// generate jwt token
	claims := jwt.MapClaims{
		"userId":      userId,
		"expiredUnix": expirationTime.Unix(),
		"authUnix":    loginTime.Unix(),
		"iss":         "user-ms",
	}

	// generate new token with claims map and method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// signing token with secret
	signingToken, err := token.SignedString(u.opts.JwtSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, loginResponse{
		AccessToken:         signingToken,
		ExpirationTimestamp: u.opts.TokenPeriodTimestamp,
	})
}

func (u *UserRouterHandler) getUserInfo(c *gin.Context) {
	userId := c.Request.Header.Get("userId")

	// query user from local
	userResult, err := u.localCon.GetUserById(userId)

	if err != nil {
		c.AbortWithError(http.StatusNotFound, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": gin.H{
			"name": userResult.Name,
			"id":   userResult.Id,
		},
	})
}

func (u *UserRouterHandler) createUser(c *gin.Context) {

	userInfo := dbcon.UserInfo{
		Id: uuid.NewString(),
	}

	c.ShouldBindJSON(&userInfo)
	err := u.localCon.CreateUser(userInfo)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": gin.H{
			"id": userInfo.Id,
		},
	})
}
