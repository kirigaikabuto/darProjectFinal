package attendances
type AttendanceRepository interface {
	Add(l *Attendance) (*Attendance,error)
	Get() ([]*Attendance,error)
	GetById(id int64) (*Attendance,error)
	Remove(l *Attendance) error
	Update(l *Attendance)  (*Attendance,error)

}
type Attendance struct {
	Id int64 `json:"id,pk"`
	LessonId int64 `json:"lessonid"`
	StudentId int64 `json:"studentid"`
}
