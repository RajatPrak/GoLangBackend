package storage

import "github.com/RajatPrak/students/internal/types"

type Storage interface {
	CreateStudent(name string, email string, age int) (int64, error)

	GetStudentById(id int64) (types.Student, error)

	GetStudents() ([]types.Student, error)

	//Done by me

	UpdateStudent(name string, email string, age int, id int64) (int64, error)

	PartiallyUpdateStudent(name string, email string, age int, id int64) (int64, error)

	// DeleteStudent(id int64) (types.Student, error)
}
