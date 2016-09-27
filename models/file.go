package models

type File struct {
	Common
	Key string
	Status string
	Mime string
}

func (file *File) Marshal() interface{} {
	return file
}
