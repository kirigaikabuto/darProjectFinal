package teachers

import "DarProject-master/users"

type TeacherRepository interface {
	AddTeacher(t *Teacher) (*Teacher,error)
	GetTeachers() ([]*Teacher,error)
	GetTeacher(id int64) (*Teacher,error)
	DeleteTeacher(t *Teacher) error
	UpdateTeacher(t *Teacher)  (*Teacher,error)
	GetTeacherByUser(user *users.User)(*Teacher,error)

}
type Teacher struct {
	users.User
	FirstName string `json:"firstname,omitempty"`
	LastName string `json:"lastname,omitempty"`
	Rating float32 `json:"rating"`
}
