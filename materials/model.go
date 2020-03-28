package materials

type MaterialRepository interface {
	Add(l *Material) (*Material,error)
	Get() ([]*Material,error)
	GetById(id int64) (*Material,error)
	Remove(l *Material) error
	Update(l *Material)  (*Material,error)
}
type Material struct {
	Id int64 `json:"id"`
	CodeLink string `json:"codelink"`
	Video string `json:"video"`
	Description string `json:"description"`
	LessonId int64 `json:"lessonid"`
}