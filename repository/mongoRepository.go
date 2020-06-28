package repository

import (
	"fmt"
	"log"
	"os"
	"gopkg.in/mgo.v2/bson"
	"github.com/joho/godotenv"
	jwt "github.com/dgrijalva/jwt-go"
	 mgo "gopkg.in/mgo.v2"
	"github.com/myrachanto/allmicro/mongomicro/usermicroservice/httperors"
	"github.com/myrachanto/allmicro/mongomicro/usermicroservice/model"
	s "github.com/myrachanto/allmicro/mongomicro/usermicroservice/support"
)

var (
	Mongorepository mongorepository = mongorepository{}
	
)

///curtesy to gorm
type mongorepository struct{}

func GetMongoDB() (*mgo.Database, error){
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	//Mongo := os.Getenv("MongoDb")
	host := os.Getenv("Mongohost")
	dbName := os.Getenv("MongodbName")

	session, err := mgo.Dial(host)
	if err != nil{
		return nil, err
	}
	db := session.DB(dbName)
	return db, nil
}
func (repository mongorepository) Create(user *model.User) (*model.User, *httperors.HttpError) {
	db, err := GetMongoDB()
	if err != nil {
		return nil, httperors.NewBadRequestError("Mongo db connection failed")
	}else{ 
		UserModel := s.UserModel{
			Db: db,
			Collection: "users",
		}
		if err := user.Validate(); err != nil {
			return nil, err
		}		
	ok, err1 := user.ValidatePassword(user.Password)
	if !ok {
		return nil, err1
	} 
	ok = user.ValidateEmail(user.Email)
	if !ok {
		return nil, httperors.NewNotFoundError("Your email format is wrong!")
	}
	UserModel = s.UserModel{
		Db: db,
		Collection: "users",
	}
	a, err := UserModel.FindEmail(user.Email)
	if err == nil {
		fmt.Println(a)
		return nil, httperors.NewNotFoundError("Your email already exists!")
	}
	
	hashpassword, err2 := user.HashPassword(user.Password)
	if err2 != nil {
		return nil, err2
	}
	user.Password = hashpassword
		user.Id = bson.NewObjectId()
		err3 := UserModel.Create(user)
		if err3 != nil {
			return nil, httperors.NewBadRequestError(fmt.Sprintf("Create user Failed, %d", err3))
		} else {
			return user, nil
		}
	}
}
func (repository mongorepository) Login(user *model.LoginUser) (*model.Auth, *httperors.HttpError) {
	db, err := GetMongoDB()
	if err != nil {
		return nil, httperors.NewBadRequestError("Mongo db connection failed")
	} else {

		UserModel := s.UserModel{
			Db: db,
			Collection: "user",
		}
		if err := user.Validate(); err != nil {
			return nil, err
		}
	
	fmt.Println(user.Email)
	e := user.Email
	UserModel = s.UserModel{
		Db: db,
		Collection: "users",
	}
	auser, err1 := UserModel.FindEmail(e)
	if err1 != nil {
	fmt.Println(err1)
		return nil, httperors.NewNotFoundError("Your email does not exists!")
	}
	ok := user.Compare(user.Password, auser.Password)
	if !ok {
		return nil, httperors.NewNotFoundError("wrong email password combo!")
	}
	tk := &model.Token{
		UserID: auser.ID,
		UserName:   auser.UName,
		Email:  auser.Email,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: model.ExpiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tk)

	encyKey := Enkey()
	tokenString, error := token.SignedString([]byte(encyKey))
	if error != nil {
		fmt.Println(error)
	}
	
	UserModel = s.UserModel{
		Db: db,
		Collection: "auth",
	}
	auth := &model.Auth{UserID:auser.Id, Token:tokenString}
	UserModel.CreateAuth(auth)

	return auth, nil
}
}
func (repository mongorepository) Logout(token string) (*httperors.HttpError) {
	db, err := GetMongoDB()
	if err != nil {
		return httperors.NewBadRequestError("Mongo db connection failed")
	}else{
		UserModel := s.UserModel{
			Db: db,
			Collection: "auth",
		}
		auth, err4 := UserModel.Findtoken(token)
		if err4 != nil {
			return httperors.NewNotFoundError("could not find users with that id" )
		}
		err2 := UserModel.DeleteAuth(auth)
		if err2 != nil {
			return httperors.NewNotFoundError(fmt.Sprintf("deletion of %d failed", err2))
		}else{
			return httperors.NewNotFoundError("deleted successfully")
		}
	}
}
func (repository mongorepository) GetOne(id string) (*model.User, *httperors.HttpError) {
	Colls := Colls()
	db, err := GetMongoDB()
	if err != nil {
		return nil, httperors.NewBadRequestError("Mongo db connection failed")
	}else{
		UserModel := s.UserModel{
			Db: db,
			Collection: Colls,
		}
		user, err2 := UserModel.Find(id)
		if err2 != nil {
			return nil,	httperors.NewNotFoundError("something went wrong")
		}
		return &user, nil
	}
}

func (repository mongorepository) GetAll(users []model.User) ([]model.User, *httperors.HttpError) {
	Colls := Colls()
	db, err := GetMongoDB()
	if err != nil {
		return nil, httperors.NewNotFoundError("Mongo db connection failed")
	}
	UserModel := s.UserModel{
		Db: db,
		Collection: Colls,
	}
	users, err2 := UserModel.FindAll()
	if err2 != nil {
		return nil,	httperors.NewNotFoundError("no results found")
	}
	return users, nil

}

func (repository mongorepository) Update(id string, user *model.User) (*model.User, *httperors.HttpError) {
	Colls := Colls()
	db, err := GetMongoDB()
	if err != nil {
		return nil, httperors.NewBadRequestError("Mongo db connection failed")
	}else{
		UserModel := s.UserModel{
			Db: db,
			Collection: Colls,
		}
		var user model.User
		uuser, err2 := UserModel.Find(id)
		if err2 != nil {
			return nil, httperors.NewBadRequestError("No user found with that id")
		}
		if user.FName  == "" {
			user.FName = uuser.FName
		}
		if user.LName  == "" {
			user.LName = uuser.LName
		}
		if user.UName  == "" {
			user.UName = uuser.UName
		}
		if user.Phone  == "" {
			user.Phone = uuser.Phone
		}
		if user.Address  == "" {
			user.Address = uuser.Address
		}
		if user.Picture  == "" {
			user.Picture = uuser.Picture
		}
		if user.Email  == "" {
			user.Email = uuser.Email
		}
		if user.Password  == "" {
			user.Password = uuser.Password
		}
		err3 := UserModel.Update(&user)
		if err3 != nil {
			return nil, httperors.NewBadRequestError("update failed for the userwith that id")
		} else {
			return &user, nil
		}
	}
}
func (repository mongorepository) Delete(id string) (*httperors.HttpSuccess, *httperors.HttpError) {
	Colls := Colls()
	db, err := GetMongoDB()
	if err != nil {
		return nil, httperors.NewBadRequestError("Mongo db connection failed")
	}else{
		UserModel := s.UserModel{
			Db: db,
			Collection: Colls,
		}
		user, err4 := UserModel.Find(id)
		if err4 != nil {
			return nil, httperors.NewNotFoundError("could not find users with that id" )
		}
		if user.Id == "" {
			return nil, httperors.NewNotFoundError("No results found")
		}
		err2 := UserModel.Delete(user)
		if err2 != nil {
			return nil, httperors.NewNotFoundError(fmt.Sprintf("deletion of %d failed", err2))
		}else{
			return httperors.NewSuccessMessage("deleted successfully"), nil
		}
	}
}

