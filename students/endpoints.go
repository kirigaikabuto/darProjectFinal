package students

import (
	"DarProject-master/courses"
	"DarProject-master/redis_connect"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type Endpoints interface {
	GetStudents() func(w http.ResponseWriter,r *http.Request)
	AddStudent() func(w http.ResponseWriter,r *http.Request)
	GetStudent(idParam string) func(w http.ResponseWriter,r *http.Request)
	DeleteStudent(idParam string) func(w http.ResponseWriter,r *http.Request)
	UpdateStudent(idParam string) func(w http.ResponseWriter,r *http.Request)
	Register() func(w http.ResponseWriter,r *http.Request)
	Login() func(w http.ResponseWriter,r *http.Request)
	Profile() func (w http.ResponseWriter,r *http.Request)
}
type endpointsFactory struct {
	studentInter Repository
	redisClient redis_connect.RedisClient
	courseInter courses.CourseRepository
}
func NewEndpointsFactory(rep Repository,redis redis_connect.RedisClient,course courses.CourseRepository) Endpoints{
	return &endpointsFactory{
		studentInter: rep,
		redisClient:redis,
		courseInter:course,
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

func (ef *endpointsFactory) Profile() func(w http.ResponseWriter,r *http.Request){
	return func(w http.ResponseWriter,r *http.Request){
		tokenString := r.Header.Get("Authorization")
		student:=&Student{}
		err:=ef.redisClient.GetKey(tokenString,student)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}
		course,err:=ef.courseInter.GetCourse(student.CourseId)
		if err!=nil {
			respondJSON(w, http.StatusInternalServerError, "Курс не назначе")
			return
		}
		respondJSON(w,http.StatusOK,course)
	}
}
func (ef *endpointsFactory) Login() func(w http.ResponseWriter,r *http.Request){
	return func(w http.ResponseWriter,r *http.Request) {
		student:=Student{}
		body,err := ioutil.ReadAll(r.Body)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}
		err=json.Unmarshal(body,&student)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}

		result,err := ef.studentInter.GetStudentByUser(&student.User)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,"Invalid Username")
			return
		}
		err = bcrypt.CompareHashAndPassword([]byte(result.Password),[]byte(student.Password))
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,"Invalid password")
			return
		}
		token:=uuid.New()
		result.Token = token.String()
		fmt.Println(result)
		err = ef.redisClient.SetKey(token.String(),result,1*time.Minute)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}
		respondJSON(w,http.StatusOK,result)
	}
}
func (ef *endpointsFactory) Register() func(w http.ResponseWriter,r *http.Request){
	return func(w http.ResponseWriter,r *http.Request) {
		student:=Student{}
		body,err:=ioutil.ReadAll(r.Body)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,"Ошибка"+err.Error())
			return
		}
		err =json.Unmarshal(body,&student)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,"Ошибка"+err.Error())
			return
		}
		_,err=ef.studentInter.GetStudentByUser(&student.User)
		if err!=nil{
			st,err := ef.studentInter.AddStudent(&student)
			if err!=nil{
				respondJSON(w,http.StatusInternalServerError,"Ошибка"+err.Error())
				return
			}
			respondJSON(w,http.StatusOK,st)
			return
		}
		respondJSON(w,http.StatusInternalServerError,"User already exist")
	}
}
func (ef *endpointsFactory) GetStudents() func(w http.ResponseWriter,r *http.Request){
	return func(w http.ResponseWriter,r *http.Request){
		students, err := ef.studentInter.GetStudents()
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, "Ошибка"+err.Error())
			return
		}
		//value_session:=sessions.GetSession(r,"user")
		//student:=&Student{}
		//if err:=json.Unmarshal([]byte(value_session),&student);err!=nil{
		//	respondJSON(w,http.StatusBadRequest,err.Error())
		//	return
		//}
		//fmt.Println(student.Username)


		respondJSON(w, http.StatusOK, students)
	}
}

func (ef *endpointsFactory) AddStudent() func(w http.ResponseWriter,r *http.Request){
	return func(w http.ResponseWriter,r *http.Request){
		data,err:=ioutil.ReadAll(r.Body)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}
		student:=&Student{}
		if err:= json.Unmarshal(data,&student);err!=nil{
			respondJSON(w,http.StatusBadRequest,err.Error())
			return
		}
		st,err:=ef.studentInter.AddStudent(student)
		if err!=nil{
			respondJSON(w,http.StatusBadRequest,err.Error())
			return
		}
		respondJSON(w,http.StatusOK,st)
	}
}
func (ef *endpointsFactory) GetStudent(idParam string) func(w http.ResponseWriter,r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
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
		student,err:=ef.studentInter.GetStudent(id)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}
		respondJSON(w,http.StatusOK,student)
	}
}
func (ef *endpointsFactory) DeleteStudent(idParam string) func(w http.ResponseWriter,r *http.Request){
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
		student,err:=ef.studentInter.GetStudent(id)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}
		err=ef.studentInter.DeleteStudent(student)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}
		respondJSON(w,http.StatusOK,"Student was deleted")
	}
}
func (ef *endpointsFactory) UpdateStudent(idParam string) func(w http.ResponseWriter,r *http.Request){
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
		student,err:=ef.studentInter.GetStudent(id)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}
		data,err:=ioutil.ReadAll(r.Body)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}
		if err:=json.Unmarshal(data,&student);err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}
		updated_student,err:=ef.studentInter.UpdateStudent(student)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err)
			return
		}
		respondJSON(w,http.StatusOK,updated_student)
	}
}


