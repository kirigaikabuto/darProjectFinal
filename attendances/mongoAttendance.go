package attendances

import (
	"DarProject-master/config"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var(
	collection *mongo.Collection
)
type attendancerepo struct{
	dbcon *mongo.Database
}
func NewAttendanceRepository(config config.MongoConfig) (AttendanceRepository,error){
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
	collection=db.Collection("attendance")
	return &attendancerepo{dbcon:db},nil
}
func(as *attendancerepo) Get()([]*Attendance,error){
	findOptions:=options.Find()
	var myattendances []*Attendance
	cur,err :=collection.Find(context.TODO(),bson.D{{}},findOptions)
	if err!=nil{
		return nil,err
	}
	for cur.Next(context.TODO()){
		var attendance Attendance
		err:=cur.Decode(&attendance)
		if err!=nil{
			return nil,err
		}
		myattendances = append(myattendances,&attendance)
	}
	if err:=cur.Err();err!=nil{
		return nil,err
	}
	cur.Close(context.TODO())
	return myattendances,nil
}
func(as *attendancerepo) Add(a *Attendance) (*Attendance,error){
	myattendances,err:=as.Get()
	n:=len(myattendances)
	if n!=0{
		last_attendance:=myattendances[n-1]
		a.Id = last_attendance.Id+1
	}else{
		a.Id = 1
	}
	_,err=collection.InsertOne(context.TODO(),a)
	if err!=nil{
		return nil,err
	}
	return a,nil
}
func (as *attendancerepo) GetById(id int64) (*Attendance,error){
	filter:=bson.D{{"id",id}}
	attendance:=&Attendance{}
	err:=collection.FindOne(context.TODO(),filter).Decode(&attendance)
	if err!=nil{
		return nil,err
	}
	return attendance,nil
}
func (as *attendancerepo) Update(a *Attendance)  (*Attendance,error){
	filter:=bson.D{{"id",a.Id}}
	update:=bson.D{{"$set",bson.D{
		{"lessonid",a.LessonId},
		{"studentid",a.StudentId},
		{"exist",a.Exist},
	}}}
	_,err:=collection.UpdateOne(context.TODO(),filter,update)
	if err!=nil{
		return nil,err
	}
	return a,nil
}
func (as *attendancerepo) Remove(a *Attendance) error{
	filter:=bson.D{{"id",a.Id}}
	_,err:=collection.DeleteOne(context.TODO(),filter)
	if err!=nil{
		return err
	}
	return nil
}
