package model

import "time"

type User struct {
	ID        uint      `gorm:"primarykey" json:"-"`
	CreatedAt time.Time `json:"-"`
	Login     string    `json:"login"`
	Password  string    `gorm:"-" json:"password"`
	Hash      string    `json:"-"`
}

func (u User) IsAlreadyExist() bool {
	return u.ID != 0
}
