package controller

import (
	"context"
	"fmt"
	"github.com/fwidyatama/e-recruitment/app/model"
	"github.com/fwidyatama/e-recruitment/util"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/vonage/vonage-go-sdk"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"os"
	"time"
)

var userCollection = util.InitDatabase().Database("e-recruitment").Collection("users")

func LoadEnv() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("failed load env file")
		panic(err)
	}
}

func SendVerify(number string) string {
	LoadEnv()
	auth := vonage.CreateAuthFromKeySecret(os.Getenv("API_KEY"), os.Getenv("API_SECRET"))
	verifyClient := vonage.NewVerifyClient(auth)
	response, errResp, err := verifyClient.Request(number, "Recruitment", vonage.VerifyOpts{CodeLength: 4, Lg: "id-id", WorkflowID: 1})
	if err != nil {
		fmt.Printf("%#v\n", err)
	} else if response.Status != "0" {
		fmt.Println("Error status " + errResp.Status + ": " + errResp.ErrorText)
	} else {
		fmt.Println("Request started: " + response.RequestId)
	}
	return response.RequestId
}

func VerifyOtp(c *gin.Context) {
	var user model.UserJson
	var verification model.Verification
	_ = c.ShouldBindJSON(&verification)

	ctx, _ := context.WithTimeout(context.TODO(), 10*time.Second)
	_ = userCollection.FindOne(ctx, bson.M{"phone_number": verification.PhoneNumber}).Decode(&user)
	LoadEnv()
	auth := vonage.CreateAuthFromKeySecret(os.Getenv("API_KEY"), os.Getenv("API_SECRET"))
	verifyClient := vonage.NewVerifyClient(auth)
	response, errResp, err := verifyClient.Check(user.RequestId, verification.Otp)

	if err != nil {
		fmt.Printf("%#v\n", err)
	} else if response.Status != "0" {
		fmt.Println("Error status " + errResp.Status + ": " + errResp.ErrorText)
	} else {
		fmt.Println("Request complete: " + response.RequestId)
	}
	if response.RequestId != "" {
		update := bson.M{"$set": bson.M{
			"verified":   true,
			"request_id": "-",
		}}
		userCollection.FindOneAndUpdate(ctx, bson.M{"phone_number": verification.PhoneNumber}, update)
	}
	c.JSON(http.StatusOK, gin.H{"Message": "Success verify identity."})
}

func RegisterUser(c *gin.Context) {
	type ErrorMessage struct {
		Message []string
	}
	var user model.User
	var duplicateEmail model.User
	var duplicatePhone model.User
	var errors ErrorMessage

	_ = c.BindJSON(&user)
	ctx, _ := context.WithTimeout(context.TODO(), 10*time.Second)
	_ = userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&duplicateEmail)
	if user.Email == duplicateEmail.Email {
		errors.Message = append(errors.Message, "Email already registered")
	}
	_ = userCollection.FindOne(ctx, bson.M{"phone_number": user.PhoneNumber}).Decode(&duplicatePhone)
	if user.PhoneNumber == duplicatePhone.PhoneNumber {
		errors.Message = append(errors.Message, "Phone number already registered")
	}

	if errors.Message != nil {
		c.JSON(http.StatusBadRequest, errors)
		return
	}
	//send verify code
	user.RequestId = SendVerify(user.PhoneNumber)
	user.Password, _ = user.HashPassword(user.Password)
	_, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		panic(err.Error())
	}
	c.JSON(http.StatusOK, gin.H{"Message": "Success register. Don't forget to verify your identity."})
}

func GetAllUser(c *gin.Context) {
	var users []model.User
	var user model.User
	var cursor *mongo.Cursor
	var err error

	verified := c.Query("verified")
	ctx, _ := context.WithTimeout(context.TODO(), 5*time.Second)

	switch verified {
	case "true":
		cursor, err = userCollection.Find(context.TODO(), bson.M{"verified": true})
	case "false":
		cursor, err = userCollection.Find(context.TODO(), bson.M{"verified": false})
	default:
		cursor, err = userCollection.Find(context.TODO(), bson.M{})
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Message": err.Error(),
		})
	}
	if cursor == nil {
		panic("cursor nil")
	}
	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		err := cursor.Decode(&user)
		if err != nil {
			log.Fatal(err.Error())
		}
		users = append(users, user)
	}
	if len(users) == 0 {
		c.JSON(http.StatusOK, gin.H{"Message": "No data available"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Data": users})
}
