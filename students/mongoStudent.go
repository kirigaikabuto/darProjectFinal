package students

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

type repo struct{
	dbcon *mongo.Database
}

func NewStudentRepository(config config.MongoConfig) (Repository,error){
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
	collection=db.Collection("students")
	return &repo{dbcon:db,},nil
}

func(mpro *repo) GetStudents() ([]*Student,error){
	findOptions:=options.Find()
	var students []*Student
	cur,err :=collection.Find(context.TODO(),bson.D{{}},findOptions)
	if err!=nil{
		return nil,err
	}
	for cur.Next(context.TODO()){
		var student Student
		err:=cur.Decode(&student)
		if err!=nil{
			return nil,err
		}
		students = append(students,&student)
	}
	if err:=cur.Err();err!=nil{
		return nil,err
	}
	cur.Close(context.TODO())
	return students,nil
}
func (mpro *repo) AddStudent(st *Student) (*Student,error){
	students,err:=mpro.GetStudents()
	n:=len(students)
	hash,err := getHashPassword(st.Password)
	if err!=nil{
		return nil,err
	}
	st.Password = hash
	if n!=0{
		student:=students[n-1]
		st.User.Id = student.User.Id+1
	}else{
		st.User.Id = 1
	}

	insertResult,err:=collection.InsertOne(context.TODO(),st)
	if err!=nil{
		return nil,err
	}
	fmt.Println("Inserted document",insertResult.InsertedID)
	return st,nil
}
func (mpro *repo) GetStudent(id int64) (*Student,error){
	filter:=bson.D{{"user.id",id}}
	student:=&Student{}
	err:=collection.FindOne(context.TODO(),filter).Decode(&student)
	if err!=nil{
		return nil,err
	}
	return student,nil
}
func (mpro *repo) DeleteStudent(st *Student) error{
	filter:=bson.D{{"user.id",st.Id}}
	_,err:=collection.DeleteOne(context.TODO(),filter)
	if err!=nil{
		return err
	}
	return nil
}
func (mpro *repo) UpdateStudent(st *Student)  (*Student,error){
	filter:=bson.D{{"user.id",st.Id}}
	hash,err := getHashPassword(st.Password)
	if err!=nil{
		return nil,err
	}
	st.Password = hash
	update:=bson.D{{"$set",bson.D{
		{"user.username",st.Username},
		{"user.password",st.Password},
		{"firstname",st.FirstName},
		{"lastname",st.LastName},
		{"courseid",st.CourseId},
	}}}
	_,err=collection.UpdateOne(context.TODO(),filter,update)
	if err!=nil{
		return nil,err
	}
	return st,nil
}
func (mpro *repo) GetStudentByUser(user *users.User)(*Student,error){
	filter:=bson.D{{"user.username",user.Username}}
	student:=&Student{}
	err:=collection.FindOne(context.TODO(),filter).Decode(&student)
	if err!=nil{
		return nil,err
	}
	return student,nil
}
func getHashPassword(password string) (string,error){
	hash,err:=bcrypt.GenerateFromPassword([]byte(password),5)
	if err!=nil{
		return "",err
	}
	return string(hash),nil
}
