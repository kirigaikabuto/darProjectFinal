package students

import "DarProject-master/users"

type Repository interface {
	GetStudents() ([]*Student,error)
	AddStudent(st *Student) (*Student,error)
	GetStudent(id int64) (*Student,error)
	DeleteStudent(st *Student) error
	UpdateStudent(st *Student) (*Student,error)
	GetStudentByUser(user *users.User)(*Student,error)

}


type Student struct{
	users.User
	FirstName string `json:"firstname,omitempty"`
	LastName string `json:"lastname,omitempty"`
	CourseId int64 `json:"courseid"`
}
