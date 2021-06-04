package vaccinatorplus

import "gorm.io/gorm"

type Conversation struct {
	gorm.Model
	ChatID    int
	FirstName string
	LastName  string
	Username  string

	Year int
}
