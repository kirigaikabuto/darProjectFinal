package main

import (
	"DarProject-master/attendances"
	"DarProject-master/config"
	"DarProject-master/courses"
	"DarProject-master/lessons"
	"DarProject-master/materials"
	"DarProject-master/redis_connect"
	"DarProject-master/schedule"
	"DarProject-master/students"
	"DarProject-master/teachers"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/urfave/cli"
	"io/ioutil"

	"net/http"
	"os"
)
var(
	PATH string  = ""
	conf config.MongoConfig
)
var flags []cli.Flag = []cli.Flag{

	&cli.StringFlag{
		Name:        "config, c",
		Usage:       "Load configuration from `FILE`",
		Destination: &PATH,
	},
	&cli.StringFlag{
		Name: "host",
		Usage: "MongoDB Hostname to connect",
		Destination: &conf.Host,
	},

	&cli.StringFlag{
		Name: "database",
		Usage: "MongoDB Name to connect",
		Destination: &conf.Database,
	},
	&cli.StringFlag{
		Name: "port",
		Usage: "MongoDB Port to connect",
		Destination: &conf.Port,
	},
}
func main(){
	app := cli.NewApp()
	app.Name = "Rest Api Cli"
	app.Flags = flags
	app.Action = runRestApi
	fmt.Println(app.Run(os.Args))
}
func ExtractConfig(path string, config *config.MongoConfig) config.MongoConfig{
	file, _ := ioutil.ReadFile(path)
	json.Unmarshal(file, &config)
	return *config
}
func runRestApi(*cli.Context) error{
	router:=mux.NewRouter()
	redisconnect:=redis_connect.ConnectRedis()
	if PATH != "" {
		ExtractConfig(PATH, &conf)
	}
	if conf.Host == "" {
		return errors.New("Nothing found in MongoDB Hostname variable")
	}
	if conf.Database == "" {
		return errors.New("Nothing found in MongoDB Database variable")
	}
	if conf.Port == "" {
		return errors.New("Nothing found in MongoDB Port variable")
	}
	//materials
	materialrepo,err:=materials.NewMaterialRepository(conf)
	if err!=nil{
		return err

	}
	materialsendpoints:=materials.NewEndpointsFactory(materialrepo)
	router.Methods("GET").Path("/materials/").HandlerFunc(materialsendpoints.GetMaterials())
	router.Methods("POST").Path("/materials/").HandlerFunc(materialsendpoints.AddMaterial())
	router.Methods("GET").Path("/materials/{id}").HandlerFunc(materialsendpoints.GetMaterial("id"))
	router.Methods("DELETE").Path("/materials/{id}").HandlerFunc(materialsendpoints.DeleteMaterial("id"))
	router.Methods("PUT").Path("/materials/{id}").HandlerFunc(materialsendpoints.UpdateMaterial("id"))
	router.Methods("GET").Path("/materials/{id}/stream/").HandlerFunc(materialsendpoints.StreamHandler("id","segName"))
	router.Methods("GET").Path("/materials/{id}/stream/{segName}").HandlerFunc(materialsendpoints.StreamHandler("id","segName"))


	//attendances
	attendancerepo,err:=attendances.NewAttendanceRepository(conf)
	if err!=nil{
		return err

	}
	attendancesendpoints:=attendances.NewEndpointsFactory(attendancerepo)
	router.Methods("GET").Path("/attendances/").HandlerFunc(attendancesendpoints.GetAttendances())
	router.Methods("POST").Path("/attendances/").HandlerFunc(attendancesendpoints.AddAttendance())
	router.Methods("GET").Path("/attendances/{id}").HandlerFunc(attendancesendpoints.GetAttendance("id"))
	router.Methods("DELETE").Path("/attendances/{id}").HandlerFunc(attendancesendpoints.DeleteAttendance("id"))
	router.Methods("PUT").Path("/attendances/{id}").HandlerFunc(attendancesendpoints.UpdateAttendance("id"))
	//lessons
	lessonsrepo,err:=lessons.NewLessonRepository(conf)
	if err!=nil{
		return err
	}
	lessonsendpoints:=lessons.NewEndpointsFactory(lessonsrepo)
	router.Methods("GET").Path("/lessons/").HandlerFunc(lessonsendpoints.GetLessons())
	router.Methods("POST").Path("/lessons/").HandlerFunc(lessonsendpoints.AddLesson())
	router.Methods("GET").Path("/lessons/{id}").HandlerFunc(lessonsendpoints.GetLesson("id"))
	router.Methods("DELETE").Path("/lessons/{id}").HandlerFunc(lessonsendpoints.DeleteLesson("id"))
	router.Methods("PUT").Path("/lessons/{id}").HandlerFunc(lessonsendpoints.UpdateLesson("id"))
	//schedule
	schedulerepo,err:=schedule.NewScheduleRepository(conf)
	if err!=nil{
		return err
	}
	scheduleendpoints:=schedule.NewEndpointsFactory(schedulerepo)
	router.Methods("GET").Path("/schedule/").HandlerFunc(scheduleendpoints.GetSchedules())
	router.Methods("POST").Path("/schedule/").HandlerFunc(scheduleendpoints.AddSchedule())
	router.Methods("GET").Path("/schedule/{id}").HandlerFunc(scheduleendpoints.GetSchedule("id"))
	router.Methods("DELETE").Path("/schedule/{id}").HandlerFunc(scheduleendpoints.DeleteSchedule("id"))
	router.Methods("PUT").Path("/schedule/{id}").HandlerFunc(scheduleendpoints.UpdateSchedule("id"))
	//courses

	coursesrepo,err:=courses.NewCourseRepository(conf)
	if err!=nil{
		return err
	}
	coursesendpoints:=courses.NewEndpointsFactory(coursesrepo,lessonsrepo)
	router.Methods("GET").Path("/courses/").HandlerFunc(coursesendpoints.GetCourses())
	router.Methods("POST").Path("/courses/").HandlerFunc(coursesendpoints.AddCourse())
	router.Methods("GET").Path("/courses/{id}").HandlerFunc(coursesendpoints.GetCourse("id"))
	router.Methods("DELETE").Path("/courses/{id}").HandlerFunc(coursesendpoints.DeleteCourse("id"))
	router.Methods("PUT").Path("/courses/{id}").HandlerFunc(coursesendpoints.UpdateCourse("id"))
	router.Methods("GET").Path("/courses/{id}/lessons/").HandlerFunc(coursesendpoints.GetLessons("id"))
	//

	//students
	studentrepo,err:=students.NewStudentRepository(conf)
	if err!=nil{
		return err
	}


	studentendpoints:=students.NewEndpointsFactory(studentrepo,redisconnect,coursesrepo)
	router.Methods("GET").Path("/students/").HandlerFunc(studentendpoints.GetStudents())
	router.Methods("GET").Path("/students/{id}").HandlerFunc(studentendpoints.GetStudent("id"))
	router.Methods("DELETE").Path("/students/{id}").HandlerFunc(studentendpoints.DeleteStudent("id"))
	router.Methods("PUT").Path("/students/{id}").HandlerFunc(studentendpoints.UpdateStudent("id"))
	router.Methods("POST").Path("/students/").HandlerFunc(studentendpoints.Register())
	router.Methods("POST").Path("/students/login/").HandlerFunc(studentendpoints.Login())
	router.Methods("POST").Path("/students/profile/").HandlerFunc(studentendpoints.Profile())


	//teachers
	teacherrepo,err:=teachers.NewTeacherRepository(conf)
	if err!=nil{
		return err
	}
	teachersendpoints:=teachers.NewEndpointsFactory(teacherrepo,redisconnect,coursesrepo)
	router.Methods("GET").Path("/teachers/").HandlerFunc(teachersendpoints.GetTeachers())
	router.Methods("GET").Path("/teachers/{id}").HandlerFunc(teachersendpoints.GetTeacher("id"))
	router.Methods("DELETE").Path("/teachers/{id}").HandlerFunc(teachersendpoints.DeleteTeacher("id"))
	router.Methods("POST").Path("/teachers/").HandlerFunc(teachersendpoints.AddTeacher())
	router.Methods("PUT").Path("/teachers/{id}").HandlerFunc(teachersendpoints.UpdateTeacher("id"))
	router.Methods("POST").Path("/teachers/login/").HandlerFunc(teachersendpoints.Login())
	router.Methods("POST").Path("/teachers/profile/").HandlerFunc(teachersendpoints.Profile())
	fmt.Println("Server is running")
	http.ListenAndServe(":8000",router)
	return nil
}