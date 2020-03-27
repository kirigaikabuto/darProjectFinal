package lessons

type LessonRepository interface {
	Add(l *Lesson) (*Lesson,error)
	Get() ([]*Lesson,error)
	GetById(id int64) (*Lesson,error)
	Remove(l *Lesson) error
	Update(l *Lesson)  (*Lesson,error)
}
type Lesson struct{
	Id int64 `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	VideoLink string `json:"videolink"`
	Resource string `json:"resource"`
	CourseId int64 `json:"courseid"`
}
