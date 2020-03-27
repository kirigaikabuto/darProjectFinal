package teachers

import (
	"DarProject-master/config"
	"DarProject-master/users"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)
var (
	collection *mongo.Collection
)
type teacherrepo struct{
	dbcon *mongo.Database
}

func NewTeacherRepository(config config.MongoConfig) (TeacherRepository,error){
	clientOptions:=options.Client().ApplyURI("mongodb://"+config.Host+":"+config.Port)
	client,err := mongo.Connect(context.TODO(),clientOptions)
	if err!=nil{
		return nil,err
	}
	err = client.Ping(context.TODO(),nil)
	if err!=nil{
		return nil,err
	}
	db:=client.Database(config.Database)
	collection=db.Collection("teachers")
	return &teacherrepo{dbcon:db,},nil
}
func(tr *teacherrepo) AddTeacher(t *Teacher) (*Teacher,error){
	teachers,err:=tr.GetTeachers()
	n:=len(teachers)
	hash,err:=getHashPassword(t.Password)
	if err!=nil{
		return nil,err
	}
	t.Password=hash
	if n!=0{
		last_teacher:=teachers[n-1]
		t.Id = last_teacher.Id+1
	}else{
		t.Id = 1
	}
	insertResult,err:=collection.InsertOne(context.TODO(),t)
	if err!=nil{
		return nil,err
	}
	fmt.Println("Inserted document",insertResult.InsertedID)
	return t,nil
}
func(tr *teacherrepo) GetTeachers() ([]*Teacher,error){
	findOptions:=options.Find()
	var teachers []*Teacher
	cur,err :=collection.Find(context.TODO(),bson.D{{}},findOptions)
	if err!=nil{
		return nil,err
	}
	for cur.Next(context.TODO()){
		var teacher Teacher
		err:=cur.Decode(&teacher)
		if err!=nil{
			return nil,err
		}
		teachers = append(teachers,&teacher)
	}
	if err:=cur.Err();err!=nil{
		return nil,err
	}
	cur.Close(context.TODO())
	return teachers,nil
}
func(tr *teacherrepo) GetTeacher(id int64) (*Teacher,error){
		filter:=bson.D{{"user.id",id}}
		teacher:=&Teacher{}
		err:=collection.FindOne(context.TODO(),filter).Decode(&teacher)
		if err!=nil{
			return nil,err
		}
		return teacher,nil
}

func (tr *teacherrepo) DeleteTeacher(t *Teacher) error{
	filter:=bson.D{{"user.id",t.Id}}
	_,err:=collection.DeleteOne(context.TODO(),filter)
	if err!=nil{
		return err
	}
	return nil
}
func (tr *teacherrepo) UpdateTeacher(t *Teacher)  (*Teacher,error){
	filter:=bson.D{{"user.id",t.Id}}
	hash,err := getHashPassword(t.Password)
	if err!=nil{
		return nil,err
	}
	t.Password = hash
	update:=bson.D{{"$set",bson.D{
		{"user.username",t.Username},
		{"user.password",t.Password},
		{"firstname",t.FirstName},
		{"lastname",t.LastName},
		{"rating",t.Rating},
	}}}
	_,err=collection.UpdateOne(context.TODO(),filter,update)
	if err!=nil{
		return nil,err
	}
	return t,nil
}
func (tr *teacherrepo) GetTeacherByUser(user *users.User)(*Teacher,error){
	filter:=bson.D{{"user.username",user.Username}}
	teacher:=&Teacher{}
	err:=collection.FindOne(context.TODO(),filter).Decode(&teacher)
	if err!=nil{
		return nil,err
	}
	return teacher,nil
}

func getHashPassword(password string) (string,error){
	hash,err:=bcrypt.GenerateFromPassword([]byte(password),2)
	if err!=nil{
		return "",err
	}
	return string(hash),nil
}
