package lessons

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
type lessonrepo struct{
	dbcon *mongo.Database
}
func NewLessonRepository(config config.MongoConfig) (LessonRepository,error){
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
	collection=db.Collection("lessons")
	return &lessonrepo{dbcon:db},nil
}
func(ls *lessonrepo) Get()([]*Lesson,error){
	findOptions:=options.Find()
	var mylessons []*Lesson
	cur,err :=collection.Find(context.TODO(),bson.D{{}},findOptions)
	if err!=nil{
		return nil,err
	}
	for cur.Next(context.TODO()){
		var lesson Lesson
		err:=cur.Decode(&lesson)
		if err!=nil{
			return nil,err
		}
		mylessons = append(mylessons,&lesson)
	}
	if err:=cur.Err();err!=nil{
		return nil,err
	}
	cur.Close(context.TODO())
	return mylessons,nil
}
func(ls *lessonrepo) Add(l *Lesson) (*Lesson,error){
	mylessons,err:=ls.Get()
	n:=len(mylessons)
	if n!=0{
		last_course:=mylessons[n-1]
		l.Id = last_course.Id+1
	}else{
		l.Id = 1
	}
	_,err=collection.InsertOne(context.TODO(),l)
	if err!=nil{
		return nil,err
	}
	return l,nil
}
func (ls *lessonrepo) GetById(id int64) (*Lesson,error){
	filter:=bson.D{{"id",id}}
	lesson:=&Lesson{}
	err:=collection.FindOne(context.TODO(),filter).Decode(&lesson)
	if err!=nil{
		return nil,err
	}
	return lesson,nil
}
func (ls *lessonrepo) Update(l *Lesson)  (*Lesson,error){
	filter:=bson.D{{"id",l.Id}}
	update:=bson.D{{"$set",bson.D{
		{"name",l.Name},
		{"description",l.Description},
		{"videolink",l.VideoLink},
		{"resource",l.Resource},
		{"courseid",l.CourseId},
	}}}
	_,err:=collection.UpdateOne(context.TODO(),filter,update)
	if err!=nil{
		return nil,err
	}
	return l,nil
}
func (ls *lessonrepo) Remove(l *Lesson) error{
	filter:=bson.D{{"id",l.Id}}
	_,err:=collection.DeleteOne(context.TODO(),filter)
	if err!=nil{
		return err
	}
	return nil
}
