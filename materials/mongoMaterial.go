package materials


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
type materialrepo struct{
	dbcon *mongo.Database
}
func NewMaterialRepository(config config.MongoConfig) (MaterialRepository,error){
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
	collection=db.Collection("materials")
	return &materialrepo{dbcon:db},nil
}
func(mat *materialrepo) Get()([]*Material,error){
	findOptions:=options.Find()
	var mymaterials []*Material
	cur,err :=collection.Find(context.TODO(),bson.D{{}},findOptions)
	if err!=nil{
		return nil,err
	}
	for cur.Next(context.TODO()){
		var material Material
		err:=cur.Decode(&material)
		if err!=nil{
			return nil,err
		}
		mymaterials = append(mymaterials,&material)
	}
	if err:=cur.Err();err!=nil{
		return nil,err
	}
	cur.Close(context.TODO())
	return mymaterials,nil
}
func(mat *materialrepo) Add(m *Material) (*Material,error){
	mymaterials,err:=mat.Get()
	n:=len(mymaterials)
	if n!=0{
		last_material:=mymaterials[n-1]
		m.Id = last_material.Id+1
	}else{
		m.Id = 1
	}
	_,err=collection.InsertOne(context.TODO(),m)
	if err!=nil{
		return nil,err
	}
	return m,nil
}
func (mat *materialrepo) GetById(id int64) (*Material,error){
	filter:=bson.D{{"id",id}}
	material:=&Material{}
	err:=collection.FindOne(context.TODO(),filter).Decode(&material)
	if err!=nil{
		return nil,err
	}
	return material,nil
}
func (mat *materialrepo) Update(m *Material)  (*Material,error){
	filter:=bson.D{{"id",m.Id}}
	update:=bson.D{{"$set",bson.D{
		{"codelink",m.CodeLink},
		{"video",m.Video},
		{"description",m.Description},
		{"lessonid",m.LessonId},
	}}}
	_,err:=collection.UpdateOne(context.TODO(),filter,update)
	if err!=nil{
		return nil,err
	}
	return m,nil
}
func (mat *materialrepo) Remove(m *Material) error{
	filter:=bson.D{{"id",m.Id}}
	_,err:=collection.DeleteOne(context.TODO(),filter)
	if err!=nil{
		return err
	}
	return nil
}
