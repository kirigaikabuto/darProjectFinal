package courses

import (
	"DarProject-master/config"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)
var (
	collection *mongo.Collection
)
type courserepo struct{
	dbcon *mongo.Database
}

func NewCourseRepository(config config.MongoConfig) (CourseRepository,error){
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
	collection=db.Collection("courses")
	return &courserepo{dbcon:db,},nil
}

func(cr *courserepo) AddCourse(course *Course) (*Course,error){
	courses,err:=cr.GetCourses()
	n:=len(courses)
	if n!=0{
		last_course:=courses[n-1]
		course.Id = last_course.Id+1
	}else{
		course.Id = 1
	}
	insertResult,err:=collection.InsertOne(context.TODO(),course)
	if err!=nil{
		return nil,err
	}
	fmt.Println("Inserted document",insertResult.InsertedID)
	return course,nil
}

func(cr *courserepo) GetCourses() ([]*Course,error){
	findOptions:=options.Find()
	var courses []*Course
	cur,err :=collection.Find(context.TODO(),bson.D{{}},findOptions)
	if err!=nil{
		return nil,err
	}
	for cur.Next(context.TODO()){
		var course Course
		err:=cur.Decode(&course)
		if err!=nil{
			return nil,err
		}
		courses = append(courses,&course)
	}
	if err:=cur.Err();err!=nil{
		return nil,err
	}
	cur.Close(context.TODO())
	return courses,nil
}
func (cr *courserepo) GetCourse(id int64) (*Course,error){
	filter:=bson.D{{"id",id}}
	course:=&Course{}
	err:=collection.FindOne(context.TODO(),filter).Decode(&course)
	if err!=nil{
		return nil,err
	}
	return course,nil
}
func (cr *courserepo) DeleteCourse(course *Course) error{
	filter:=bson.D{{"id",course.Id}}
	_,err:=collection.DeleteOne(context.TODO(),filter)
	if err!=nil{
		return err
	}
	return nil
}
func (cr *courserepo) GetCourseByTeacherId(id int64)(*Course,error){
	filter:=bson.D{{"teacherid",id}}
	course:=&Course{}
	err:=collection.FindOne(context.TODO(),filter).Decode(&course)
	if err!=nil{
		return nil,err
	}
	return course,nil
}
func (cr *courserepo) UpdateCourse(course *Course)  (*Course,error){
	filter:=bson.D{{"id",course.Id}}
	update:=bson.D{{"$set",bson.D{
		{"name",course.Name},
		{"datestart",course.DateStart},
		{"dateend",course.DateEnd},
		{"description",course.Description},
		{"teacherid",course.TeacherId},
	}}}
	_,err:=collection.UpdateOne(context.TODO(),filter,update)
	if err!=nil{
		return nil,err
	}
	return course,nil
}
