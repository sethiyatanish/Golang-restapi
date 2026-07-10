package types

type Student struct{
	Id int 
	Name string `validate:"required"`
	Email string `validate:"required"`
	age int `validate: "required"`
}