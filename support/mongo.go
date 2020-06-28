package models
import ( 
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/myrachanto/allmicro/mongomicro/usermicroservice/model"
) 

type UserModel struct {
	Db *mgo.Database
	Collection string
}
func (usermodel UserModel) FindAll() (users []model.User, err error){
	err = usermodel.Db.C(usermodel.Collection).Find(bson.M{}).All(&users)
	return
}

func (usermodel UserModel) Find(id string) (user model.User, err error){
	err = usermodel.Db.C(usermodel.Collection).FindId(bson.ObjectIdHex(id)).One(&user)
	return
}
func (usermodel UserModel) FindEmail(email string) (user model.User, err error){
	err = usermodel.Db.C(usermodel.Collection).Find(bson.M{"email": email}).One(&user)
	return
}
func (usermodel UserModel) Create(user *model.User) error {
	err := usermodel.Db.C(usermodel.Collection).Insert(&user)
	return err
}
func (usermodel UserModel) CreateAuth(auth *model.Auth) error {
	err := usermodel.Db.C(usermodel.Collection).Insert(&auth)
	return err
}
func (usermodel UserModel) Findtoken(token string) (auth model.Auth, err error){
	err = usermodel.Db.C(usermodel.Collection).Find(bson.M{"token": token}).One(&auth)
	return
}

func (usermodel UserModel) Update(user *model.User) error {
	err := usermodel.Db.C(usermodel.Collection).UpdateId(user.Id, &user)
	return err
}
func (usermodel UserModel) Delete(user model.User) error {
	err := usermodel.Db.C(usermodel.Collection).Remove(user)
	return err
}
func (usermodel UserModel) DeleteAuth(auth model.Auth) error {
	err := usermodel.Db.C(usermodel.Collection).Remove(auth)
	return err
}