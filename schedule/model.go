package schedule

type ScheduleRepo interface {
	Add(l *Schedule) (*Schedule,error)
	Get() ([]*Schedule,error)
	GetById(id int64) (*Schedule,error)
	Remove(l *Schedule) error
	Update(l *Schedule)  (*Schedule,error)
}
type Schedule struct {
	Id int64 `json:"id"`
	CourseId int64 `json:"courseid"`
}

