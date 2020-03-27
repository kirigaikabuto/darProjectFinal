package courses

import (
	"DarProject-master/lessons"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"encoding/json"
	"strconv"
)
type Endpoints interface {
	AddCourse() func(w http.ResponseWriter,r *http.Request)
	GetCourses() func(w http.ResponseWriter,r *http.Request)
	GetCourse(idParam string) func(w http.ResponseWriter,r *http.Request)
	DeleteCourse(idParam string) func(w http.ResponseWriter,r *http.Request)
	UpdateCourse(idParam string) func(w http.ResponseWriter,r *http.Request)
	GetLessons(idParam string) func(w http.ResponseWriter,r *http.Request)
}
type endpointsFactory struct {
	courseRep CourseRepository
	lessonRep lessons.LessonRepository
}
func NewEndpointsFactory(rep CourseRepository,les lessons.LessonRepository) Endpoints{
	return &endpointsFactory{
		courseRep: rep,
		lessonRep:les,
	}
}
func (ef *endpointsFactory) GetLessons(idParam string) func(w http.ResponseWriter,r *http.Request){
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
		course,err:=ef.courseRep.GetCourse(id)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}
		lessons,err:=ef.lessonRep.GetLessonsByCourseId(course.Id)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
		}
		respondJSON(w,http.StatusOK,lessons)
	}
}
func (ef *endpointsFactory) AddCourse() func(w http.ResponseWriter,r *http.Request){
	return func(w http.ResponseWriter,r *http.Request){
		data,err:=ioutil.ReadAll(r.Body)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}

		course:=&Course{}
		if err:= json.Unmarshal(data,&course);err!=nil{
			respondJSON(w,http.StatusBadRequest,err.Error())
			return
		}


		course,err=ef.courseRep.AddCourse(course)
		if err!=nil{
			respondJSON(w,http.StatusBadRequest,err.Error())
			return
		}
		respondJSON(w,http.StatusOK,course)
	}
}
func (ef *endpointsFactory) GetCourses() func(w http.ResponseWriter,r *http.Request){
	return func(w http.ResponseWriter,r *http.Request){
		courses, err := ef.courseRep.GetCourses()
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, "Ошибка"+err.Error())
			return
		}

		respondJSON(w, http.StatusOK, courses)
	}
}
func (ef *endpointsFactory) GetCourse(idParam string) func(w http.ResponseWriter,r *http.Request) {
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
		course,err:=ef.courseRep.GetCourse(id)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}
		respondJSON(w,http.StatusOK,course)
	}
}
func (ef *endpointsFactory) DeleteCourse(idParam string) func(w http.ResponseWriter,r *http.Request){
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
		course,err:=ef.courseRep.GetCourse(id)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}
		err=ef.courseRep.DeleteCourse(course)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}
		respondJSON(w,http.StatusOK,"Course was deleted")
	}
}
func (ef *endpointsFactory) UpdateCourse(idParam string) func(w http.ResponseWriter,r *http.Request){
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
		course,err:=ef.courseRep.GetCourse(id)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}
		data,err:=ioutil.ReadAll(r.Body)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}
		if err:=json.Unmarshal(data,&course);err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}
		updated_student,err:=ef.courseRep.UpdateCourse(course)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err)
			return
		}
		respondJSON(w,http.StatusOK,updated_student)
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
