package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Name        string             `json:"name,,omitempty" bson:"name,omitempty"`
	Email       string             `json:"email,omitempty" bson:"email,omitempty"`
	Password    string             `json:"-" bson:"password,omitempty"`
	PhoneNumber string             `json:"phone_number,omitempty" bson:"phone_number,omitempty"`
	Address     string             `json:"address,omitempty" bson:"address,omitempty"`
	Nationality string             `json:"nationality,omitempty" bson:"nationality,omitempty"`
	Position    string             `json:"position,omitempty" bson:"position,omitempty"`
	Verified    bool               `json:"verified" bson:"verified"`
	RequestId   string             `json:"-" bson:"request_id"`
}

type UserJson struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
	Name        string             `json:"name,,omitempty" bson:"name,omitempty"`
	Email       string             `json:"email,omitempty" bson:"email,omitempty"`
	Password    string             `json:"password,omitempty" bson:"password,omitempty"`
	PhoneNumber string             `json:"phone_number,omitempty" bson:"phone_number,omitempty"`
	Address     string             `json:"address,omitempty" bson:"address,omitempty"`
	Nationality string             `json:"nationality,omitempty" bson:"nationality,omitempty"`
	Position    string             `json:"position,omitempty" bson:"position,omitempty"`
	Verified    bool               `json:"verified" bson:"verified"`
	RequestId   string             `json:"request_id" bson:"request_id"`
}

type Verification struct {
	Otp         string `json:"otp"`
	PhoneNumber string `json:"phone_number"`
}

func (user User) HashPassword(password string) (string, error) {
	hash, errs := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hash), errs
}
