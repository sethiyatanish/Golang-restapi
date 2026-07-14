package storage

type Storage interface {
	CreateStudent(name string, email string, afe int) (int64, error)

}