package schedule



import (
"encoding/json"
"github.com/gorilla/mux"
"io/ioutil"
"net/http"
"strconv"
)

type Endpoints interface {
	AddSchedule() func(w http.ResponseWriter,r *http.Request)
	GetSchedules() func(w http.ResponseWriter,r *http.Request)
	GetSchedule(idParam string) func(w http.ResponseWriter,r *http.Request)
	DeleteSchedule(idParam string) func(w http.ResponseWriter,r *http.Request)
	UpdateSchedule(idParam string) func(w http.ResponseWriter,r *http.Request)
}
type endpointsFactory struct {
	scheduleRep ScheduleRepo
}
func NewEndpointsFactory(scheduleRep ScheduleRepo) Endpoints{
	return &endpointsFactory{
		scheduleRep:scheduleRep,
	}
}
func(ef *endpointsFactory) AddSchedule() func(w http.ResponseWriter,r *http.Request){
	return func(w http.ResponseWriter,r *http.Request){
		data,err:=ioutil.ReadAll(r.Body)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}
		schedule:=&Schedule{}
		if err:= json.Unmarshal(data,&schedule);err!=nil{
			respondJSON(w,http.StatusBadRequest,err.Error())
			return
		}
		st,err:=ef.scheduleRep.Add(schedule)
		if err!=nil{
			respondJSON(w,http.StatusBadRequest,err.Error())
			return
		}
		respondJSON(w,http.StatusOK,st)
	}
}
func(ef *endpointsFactory) GetSchedules() func(w http.ResponseWriter,r *http.Request){
	return func(w http.ResponseWriter,r *http.Request) {
		myschedules, err := ef.scheduleRep.Get()
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, "Ошибка"+err.Error())
			return
		}
		respondJSON(w, http.StatusOK, myschedules)
	}
}
func (ef *endpointsFactory) GetSchedule(idParam string) func(w http.ResponseWriter,r *http.Request) {
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
		schedule,err:=ef.scheduleRep.GetById(id)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}
		respondJSON(w,http.StatusOK,schedule)
	}
}
func (ef *endpointsFactory) DeleteSchedule(idParam string) func(w http.ResponseWriter,r *http.Request){
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
		schedule,err:=ef.scheduleRep.GetById(id)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}
		err=ef.scheduleRep.Remove(schedule)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}
		respondJSON(w,http.StatusOK,"Schedule was deleted")
	}
}
func(ef *endpointsFactory) UpdateSchedule(idParam string) func(w http.ResponseWriter,r *http.Request){
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
		schedule,err:=ef.scheduleRep.GetById(id)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}
		data,err:=ioutil.ReadAll(r.Body)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}
		if err:=json.Unmarshal(data,&schedule);err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}
		updatedSchedule,err:=ef.scheduleRep.Update(schedule)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err)
			return
		}
		respondJSON(w,http.StatusOK,updatedSchedule)
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

