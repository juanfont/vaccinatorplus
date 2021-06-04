package vaccinatorplus

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func (v *Vaccinator) initDB() error {
	db, err := gorm.Open(sqlite.Open(v.dbPath), &gorm.Config{})
	if err != nil {
		return err
	}
	db.AutoMigrate(&Conversation{})
	return err
}

func (v *Vaccinator) db() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(v.dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
