package repository
import (
	"log"
	"os"
	mgo "gopkg.in/mgo.v2"
	"github.com/joho/godotenv"
	"github.com/myrachanto/allmicro/mongomicro/usermicroservice/httperors"
	"github.com/myrachanto/allmicro/mongomicro/usermicroservice/model"
)

type Redirectrepository interface{
	Create(user *model.User) (*model.User, *httperors.HttpError)
	Login(user *model.LoginUser) (*model.Auth, *httperors.HttpError)
	Logout(token string) (*httperors.HttpError)
	GetOne(id string) (*model.User, *httperors.HttpError)
	GetAll(users []model.User) ([]model.User, *httperors.HttpError)
	Update(id string, user *model.User) (*model.User, *httperors.HttpError)
	Delete(id string) (*httperors.HttpSuccess, *httperors.HttpError)
}

/////////////////////////////////////////////////////////////////////////////////////
////////////////figure how to switch repositories automatically//////////////////////////////////
func ChooseRepo() (repository Redirectrepository) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	switch os.Getenv("DbType2") {
	case "mongo":
		_, err1 := NewMongoRepository()
		if err1 != nil {
			log.Fatal(err1)
		}
		repository = Mongorepository
		return repository
	// case "mysql":
	// 	_, err1 := NewGormRepository()
	// 	if err1 != nil {
	// 		log.Fatal(err1)
	// 	}
	// 	repository = Sqlrepository
	// 	// model.CheckMongo(gorm)
	// 	return repository
	
	// case "postgress":
	// 	repository, err1 := NewMongoRepository()
	// 	if err1 != nil {
	// 		log.Fatal(err1)
	// 	}
	// 	return repository
	// case "redis":
	// 	repository, err1 := NewMongoRepository()
	// 	if err1 != nil {
	// 		log.Fatal(err1)
	// 	}
	}
	return
	
}
func NewMongoRepository()(Redirectrepository, error){
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	//Mongo := os.Getenv("MongoDb")
	host := os.Getenv("Mongohost")

	_, err = mgo.Dial(host)
	if err != nil{
		return nil, err
	}
	return Mongorepository, nil
}
func Enkey()string{
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	key := os.Getenv("EncryptionKey")
	return key
}
func Colls()string{
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	colls := os.Getenv("Collection")
	return colls
}