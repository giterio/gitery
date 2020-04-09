package models

// Text ...
type Text interface {
	Fetch(id int) (err error)
	Create() (err error)
	Update() (err error)
	Delete() (err error)
}
