package teachers

import (
	"DarProject-master/courses"
	"DarProject-master/redis_connect"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"strconv"
	"time"
	"github.com/google/uuid"
)

type Endpoints interface {
	AddTeacher() func(w http.ResponseWriter,r *http.Request)
	GetTeachers() func(w http.ResponseWriter,r *http.Request)
	GetTeacher(idParam string) func(w http.ResponseWriter,r *http.Request)
	DeleteTeacher(idParam string) func(w http.ResponseWriter,r *http.Request)
	UpdateTeacher(idParam string) func(w http.ResponseWriter,r *http.Request)
	Login() func(w http.ResponseWriter,r *http.Request)
	Profile() func (w http.ResponseWriter,r *http.Request)
}
type endpointsFactory struct {
	teacherRep TeacherRepository
	courseRep courses.CourseRepository
	redisClient redis_connect.RedisClient
}
func NewEndpointsFactory(rep TeacherRepository,redis redis_connect.RedisClient,courseRep courses.CourseRepository) Endpoints{
	return &endpointsFactory{
		teacherRep: rep,
		redisClient:redis,
		courseRep:courseRep,
	}
}
func(ef *endpointsFactory) Profile() func(w http.ResponseWriter,r *http.Request){
	return func(w http.ResponseWriter,r *http.Request){
		tokenString := r.Header.Get("Authorization")
		teacher:=&Teacher{}
		err:=ef.redisClient.GetKey(tokenString,teacher)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}
		course,err:=ef.courseRep.GetCourseByTeacherId(teacher.Id)
		if err!=nil {
			respondJSON(w, http.StatusInternalServerError, "Нет курсов")
			return
		}
		respondJSON(w,http.StatusOK,course)
	}
}

func(ef *endpointsFactory) Login() func(w http.ResponseWriter,r *http.Request){
	return func(w http.ResponseWriter,r *http.Request){
		teacher:=Teacher{}
		body,err := ioutil.ReadAll(r.Body)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}
		err=json.Unmarshal(body,&teacher)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}

		result,err := ef.teacherRep.GetTeacherByUser(&teacher.User)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,"Invalid Username")
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(result.Password),[]byte(teacher.Password))
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,"Invalid password")
			return
		}
		token:=uuid.New()
		result.Token = token.String()
		err = ef.redisClient.SetKey(token.String(),result,1*time.Minute)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}
		respondJSON(w,http.StatusOK,result)
	}
}
func (ef *endpointsFactory) AddTeacher() func(w http.ResponseWriter,r *http.Request){
	return func(w http.ResponseWriter,r *http.Request){
		data,err:=ioutil.ReadAll(r.Body)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}

		teacher:=&Teacher{}
		if err:= json.Unmarshal(data,&teacher);err!=nil{
			respondJSON(w,http.StatusBadRequest,err.Error())
			return
		}


		teacher,err=ef.teacherRep.AddTeacher(teacher)
		if err!=nil{
			respondJSON(w,http.StatusBadRequest,err.Error())
			return
		}
		respondJSON(w,http.StatusOK,teacher)
	}
}
func (ef *endpointsFactory) GetTeachers() func(w http.ResponseWriter,r *http.Request){
	return func(w http.ResponseWriter,r *http.Request){
		teachers, err := ef.teacherRep.GetTeachers()
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, "Ошибка"+err.Error())
			return
		}

		respondJSON(w, http.StatusOK, teachers)
	}
}
func (ef *endpointsFactory) GetTeacher(idParam string) func(w http.ResponseWriter,r *http.Request){
	return func(w http.ResponseWriter,r *http.Request){
		vars:=mux.Vars(r)
		paramid,paramerr:=vars[idParam]
		if !paramerr{
			respondJSON(w,http.StatusBadRequest,"Не был передан аргумент")
			return
		}
		id,err:=strconv.ParseInt(paramid,10,10)
		if err!=nil{
			respondJSON(w,http.StatusBadRequest,err.Error())
			return
		}
		teacher,err:=ef.teacherRep.GetTeacher(id)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}
		respondJSON(w,http.StatusOK,teacher)
	}
}
func (ef *endpointsFactory) DeleteTeacher(idParam string) func(w http.ResponseWriter,r *http.Request){
	return func(w http.ResponseWriter,r *http.Request){
		vars:=mux.Vars(r)
		paramid,paramerr:=vars[idParam]
		if !paramerr{
			respondJSON(w,http.StatusBadRequest,"Не был передан аргумент")
			return
		}
		id,err:=strconv.ParseInt(paramid,10,10)
		if err!=nil{
			respondJSON(w,http.StatusBadRequest,err.Error())
			return
		}
		teacher,err:=ef.teacherRep.GetTeacher(id)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}
		err=ef.teacherRep.DeleteTeacher(teacher)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}
		respondJSON(w,http.StatusOK,"Teacher was deleted")
	}
}
func (ef *endpointsFactory) UpdateTeacher(idParam string) func(w http.ResponseWriter,r *http.Request){
	return func(w http.ResponseWriter,r *http.Request){
		vars:=mux.Vars(r)
		paramid,paramerr:=vars[idParam]
		if !paramerr{
			respondJSON(w,http.StatusBadRequest,"Не был передан аргумент")
			return
		}
		id,err:=strconv.ParseInt(paramid,10,10)
		if err!=nil{
			respondJSON(w,http.StatusBadRequest,err.Error())
			return
		}
		teacher,err:=ef.teacherRep.GetTeacher(id)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}
		data,err:=ioutil.ReadAll(r.Body)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}
		if err:=json.Unmarshal(data,&teacher);err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}
		updated_teacher,err:=ef.teacherRep.UpdateTeacher(teacher)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err)
			return
		}
		respondJSON(w,http.StatusOK,updated_teacher)
	}
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(response))
}
