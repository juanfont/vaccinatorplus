package vaccinatorplus

import (
	"fmt"

	"gorm.io/gorm"
)

type Conversation struct {
	gorm.Model
	ChatID    int64 `gorm:"index"`
	FirstName string
	LastName  string
	Username  string

	NotifiedYear   int
	RequestedYear  int
	NotifyAllYears *bool
}

func (c Conversation) ToHumanName() string {
	return fmt.Sprintf("%s %s (%s)", c.FirstName, c.LastName, c.Username)
}
