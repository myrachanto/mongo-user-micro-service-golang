package model

import (
	"regexp"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/jinzhu/gorm"
	"github.com/myrachanto/allmicro/mongomicro/usermicroservice/httperors"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

var ExpiresAt = time.Now().Add(time.Minute * 100000).Unix()

type User struct {
	Id       bson.ObjectId `json:"_id" bson:"_id"`
	FName    string        `bson:"fname"`
	LName    string        `bson:"lname"`
	UName    string        `bson:"uname"`
	Phone    string        `bson:"phone"`
	Address  string        `bson:"address"`
	Dob      *time.Time
	Picture  string `bson:"picture"`
	Email    string `bson:"email"`
	Password string `bson:"password"`
	gorm.Model
}
type Auth struct {
	//User User `gorm:"foreignKey:UserID; not null"`
	UserID string   `json:"userid" bson:"userid"`
	Token  string `bson:"token"`
	gorm.Model
}
type LoginUser struct {
	Email    string `bson:"email"`
	Password string `bson:"password"`
}

//Token struct declaration
type Token struct {
	UserID   string
	UserName string
	Email    string
	*jwt.StandardClaims
}

func (user User) ValidateEmail(email string) (matchedString bool) {
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&amp;'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	matchedString = re.MatchString(email)
	return matchedString
}
func (user User) ValidatePassword(password string) (bool, *httperors.HttpError) {
	if len(password) < 5 {
		return false, httperors.NewBadRequestError("your password need more characters!")
	} else if len(password) > 32 {
		return false, httperors.NewBadRequestError("your password is way too long!")
	}
	return true, nil
}
func (user User) HashPassword(password string) (string, *httperors.HttpError) {
	pass, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return "", httperors.NewNotFoundError("type a stronger password!")
	}
	return string(pass), nil

}
func (user LoginUser) Compare(p1, p2 string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p2), []byte(p1))
	if err != nil {
		return false
	}
	return true
}
func (loginuser LoginUser) Validate() *httperors.HttpError {
	if loginuser.Email == "" {
		return httperors.NewNotFoundError("Invalid Email")
	}
	if loginuser.Password == "" {
		return httperors.NewNotFoundError("Invalid password")
	}
	return nil
}
func (user User) Validate() *httperors.HttpError {
	if user.FName == "" {
		return httperors.NewNotFoundError("Invalid first Name")
	}
	if user.LName == "" {
		return httperors.NewNotFoundError("Invalid last name")
	}
	if user.UName == "" {
		return httperors.NewNotFoundError("Invalid username")
	}
	if user.Phone == "" {
		return httperors.NewNotFoundError("Invalid phone number")
	}
	if user.Email == "" {
		return httperors.NewNotFoundError("Invalid Email")
	}
	if user.Address == "" {
		return httperors.NewNotFoundError("Invalid Address")
	}
	if user.Password == "" {
		return httperors.NewNotFoundError("Invalid password")
	}
	// if user.Picture == "" {
	// 	return httperors.NewNotFoundError("Invalid picture")
	// }
	if user.Email == "" {
		return httperors.NewNotFoundError("Invalid picture")
	}
	return nil
}
