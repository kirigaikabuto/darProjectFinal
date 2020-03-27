package courses



type CourseRepository interface {
	AddCourse(course *Course) (*Course,error)
	GetCourses() ([]*Course,error)
	GetCourse(id int64) (*Course,error)
	DeleteCourse(course *Course) error
	UpdateCourse(course *Course) (*Course,error)
	GetCourseByTeacherId(id int64)(*Course,error)
}
type Course struct {
	Id int64 `json:"id,pk"`
	Name string `json:"name,omitempty"`
	DateStart string `json:"datestart,omitempty"`
	DateEnd string `json:"dateend,omitempty"`
	Description string `json:"description,omitempty"`
	TeacherId int64 `json:"teacherid"`
}


