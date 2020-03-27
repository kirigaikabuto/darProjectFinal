package schedule

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
type scheduleRepo struct{
	dbcon *mongo.Database
}
func NewScheduleRepository(config config.MongoConfig) (ScheduleRepo,error){
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
	collection=db.Collection("schedule")
	return &scheduleRepo{dbcon:db},nil
}
func(ls *scheduleRepo) Get()([]*Schedule,error){
	findOptions:=options.Find()
	var mysch []*Schedule
	cur,err :=collection.Find(context.TODO(),bson.D{{}},findOptions)
	if err!=nil{
		return nil,err
	}
	for cur.Next(context.TODO()){
		var sch Schedule
		err:=cur.Decode(&sch)
		if err!=nil{
			return nil,err
		}
		mysch = append(mysch,&sch)
	}
	if err:=cur.Err();err!=nil{
		return nil,err
	}
	cur.Close(context.TODO())
	return mysch,nil
}
func(ls *scheduleRepo) Add(sch *Schedule) (*Schedule,error){
	schs,err:=ls.Get()
	n:=len(schs)
	if n!=0{
		last_sch:=schs[n-1]
		sch.Id = last_sch.Id+1
	}else{
		sch.Id = 1
	}
	_,err=collection.InsertOne(context.TODO(),sch)
	if err!=nil{
		return nil,err
	}
	return sch,nil
}
func (ls *scheduleRepo) GetById(id int64) (*Schedule,error){
	filter:=bson.D{{"id",id}}
	sh:=&Schedule{}
	err:=collection.FindOne(context.TODO(),filter).Decode(&sh)
	if err!=nil{
		return nil,err
	}
	return sh,nil
}
func (ls *scheduleRepo) Update(sh *Schedule)  (*Schedule,error){
	filter:=bson.D{{"id",sh.Id}}
	update:=bson.D{{"$set",bson.D{
		{"courseid",sh.CourseId},
	}}}
	_,err:=collection.UpdateOne(context.TODO(),filter,update)
	if err!=nil{
		return nil,err
	}
	return sh,nil
}

func (ls *scheduleRepo) Remove(l *Schedule) error{
	filter:=bson.D{{"id",l.Id}}
	_,err:=collection.DeleteOne(context.TODO(),filter)
	if err!=nil{
		return err
	}
	return nil
}
