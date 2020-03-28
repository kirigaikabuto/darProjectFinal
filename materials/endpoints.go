package materials

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)
var (
	fileuploaddir = "./assets/media/"
)
type Endpoints interface {
	AddMaterial() func(w http.ResponseWriter,r *http.Request)
	GetMaterials() func(w http.ResponseWriter,r *http.Request)
	GetMaterial(idParam string) func(w http.ResponseWriter,r *http.Request)
	DeleteMaterial(idParam string) func(w http.ResponseWriter,r *http.Request)
	UpdateMaterial(idParam string) func(w http.ResponseWriter,r *http.Request)
	StreamHandler(idParam string,segParam string) func(w http.ResponseWriter,r *http.Request)
}
type endpointsFactory struct {
	materialRep MaterialRepository
}
func NewEndpointsFactory(materialrep MaterialRepository) Endpoints{
	return &endpointsFactory{
		materialRep:materialrep,
	}
}
func(ef *endpointsFactory) AddMaterial() func(w http.ResponseWriter,r *http.Request){
	return func(w http.ResponseWriter,r *http.Request){
		enableCors(&w)
		lessonName:="lesson_naruto"
		file, handle, err := r.FormFile("video")
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}
		material:=&Material{
			CodeLink:    r.FormValue("codelink"),
			Description: r.FormValue("description"),
			Video:lessonName,
			LessonId:    0,
		}
		defer file.Close()
		st,err:=ef.materialRep.Add(material)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}
		saveFile(w, lessonName,file, handle)
		respondJSON(w,http.StatusOK,st)
		//data,err:=ioutil.ReadAll(r.Body)
		//if err!=nil{
		//	respondJSON(w,http.StatusInternalServerError,err.Error())
		//	return
		//}
		//material:=&Material{}
		//if err:= json.Unmarshal(data,&material);err!=nil{
		//	respondJSON(w,http.StatusBadRequest,err.Error())
		//	return
		//}
		//st,err:=ef.materialRep.Add(material)
		//if err!=nil{
		//	respondJSON(w,http.StatusBadRequest,err.Error())
		//	return
		//}
		//respondJSON(w,http.StatusOK,st)


	}
}
func saveFile(w http.ResponseWriter, lessonName string,file multipart.File, handle *multipart.FileHeader) {

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err)
		return
	}

	videoname := lessonName
	err = os.Mkdir(fileuploaddir+videoname, 0700)
	if err != nil {
		log.Fatal(err)
		return
	}
	folder_hls := fileuploaddir + videoname + "/hls/"
	err = os.Mkdir(fileuploaddir+videoname+"/hls/", 0700)
	if err != nil {
		log.Fatal(err)
		return
	}
	file_url := fileuploaddir + videoname + "/" + handle.Filename

	err = ioutil.WriteFile(file_url, data, 0700)
	if err != nil {
		log.Fatal(err)
		return
	}
	//ffmpeg -i input.mp4 -profile:v baseline -level 3.0 -s 640x360 -start_number 0 -hls_time 10 -hls_list_size 0 -f hls index.m3u8
	fmt.Println(file_url)
	go func(file_url, folder_hls, videoname string) {
		cmd := exec.Command("ffmpeg", "-i", file_url, "-profile:v", "baseline", "-level", "3.0", "-s", "640x360", "-start_number", "0", "-hls_time", "10", "-hls_list_size", "0", "-f", "hls", folder_hls+"playlist"+videoname+".m3u8")
		cmd.Run()
	}(file_url, folder_hls, videoname)
	respondJSON(w, http.StatusCreated, "File uploaded successfully!.")
}
func(ef *endpointsFactory)  StreamHandler(idParam ,segParam string) func(w http.ResponseWriter,r *http.Request){
	return func(w http.ResponseWriter,r *http.Request) {
		vars := mux.Vars(r)
		mId := vars[idParam]
		id, _ := strconv.ParseInt(mId, 10, 10)
		fmt.Println("asdsad")
		segName, ok := vars[segParam]
		material, err := ef.materialRep.GetById(id)
		fmt.Println(material)
		if err != nil {
			respondJSON(w, http.StatusBadRequest, err.Error())
		}
		if !ok {
			fmt.Println("in seq id")

			mediaBase := getMediaBase(material.Video)
			fmt.Println(mediaBase)
			m3u8Name := fmt.Sprintf("playlist%s.m3u8", material.Video)
			serverHlsM3u8(w, r, mediaBase, m3u8Name)
		} else {
			fmt.Println("in seg name")
			mediaBase := getMediaBase(material.Video)
			fmt.Println(mediaBase)
			serverHlsTs(w, r, mediaBase, segName)
		}
	}
}

func getMediaBase(mId string) string {
	mediaRoot := "assets/media"
	return fmt.Sprintf("%s/%s", mediaRoot, mId)
}
func serverHlsM3u8(w http.ResponseWriter, r *http.Request, mediaBase, m3u8Name string) {
	enableCors(&w)
	mediaFile := fmt.Sprintf("%s/hls/%s", mediaBase, m3u8Name)
	fmt.Println(mediaFile)
	http.ServeFile(w, r, mediaFile)
	w.Header().Set("Content-Type", "application/x-mpegURL")

}
func serverHlsTs(w http.ResponseWriter, r *http.Request, mediaBase, segName string) {
	enableCors(&w)
	mediaFile := fmt.Sprintf("%s/hls/%s", mediaBase, segName)
	fmt.Println(mediaFile)
	http.ServeFile(w, r, mediaFile)
	w.Header().Set("Content-Type", "video/MP2T")

}
func(ef *endpointsFactory) GetMaterials() func(w http.ResponseWriter,r *http.Request){
	return func(w http.ResponseWriter,r *http.Request) {
		mymaterials, err := ef.materialRep.Get()
		if err != nil {
			respondJSON(w, http.StatusInternalServerError, "Ошибка"+err.Error())
			return
		}
		respondJSON(w, http.StatusOK, mymaterials)
	}
}
func (ef *endpointsFactory) GetMaterial(idParam string) func(w http.ResponseWriter,r *http.Request) {
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
		material,err:=ef.materialRep.GetById(id)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}
		respondJSON(w,http.StatusOK,material)
	}
}
func (ef *endpointsFactory) DeleteMaterial(idParam string) func(w http.ResponseWriter,r *http.Request){
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
		material,err:=ef.materialRep.GetById(id)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}
		err=ef.materialRep.Remove(material)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}
		respondJSON(w,http.StatusOK,"Material was deleted")
	}
}
func(ef *endpointsFactory) UpdateMaterial(idParam string) func(w http.ResponseWriter,r *http.Request){
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
		material,err:=ef.materialRep.GetById(id)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}
		data,err:=ioutil.ReadAll(r.Body)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}
		if err:=json.Unmarshal(data,&material);err!=nil{
			respondJSON(w,http.StatusInternalServerError,err.Error())
			return
		}
		updated_material,err:=ef.materialRep.Update(material)
		if err!=nil{
			respondJSON(w,http.StatusInternalServerError,err)
			return
		}
		respondJSON(w,http.StatusOK,updated_material)
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

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}