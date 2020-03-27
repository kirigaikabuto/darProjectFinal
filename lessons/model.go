package lessons

type LessonRepository interface {
	Add(l *Lesson) (*Lesson,error)
	Get() ([]*Lesson,error)
	GetById(id int64) (*Lesson,error)
	Remove(l *Lesson) error
	Update(l *Lesson)  (*Lesson,error)
	GetLessonsByCourseId(id int64) ([]*Lesson,error)
}
type Lesson struct{
	Id int64 `json:"id"`
	Name string `json:"name"`
	Date string   `json:"date"`
	TimeStart string `json:"timestart"`
	TimeEnd string `json:"timeend"`
	ScheduleId int64 `json:"scheduleid"`
}
//Description string `json:"description"`
//VideoLink string `json:"videolink"`
//Resource string `json:"resource"`