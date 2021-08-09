package database

import (
	"gorm.io/gorm"
	"time"
)

var ModelsList = []Model{
	&User{}, &Password{}, &Collection{}, &Token{},
}

type Model interface {}

type User struct {
	gorm.Model
	Id					int			`gorm:"primaryKey"`
	Username			string
	Password			string
	RegistrationDate	time.Time
}

type Password struct {
	gorm.Model
	Id					int			`gorm:"primaryKey"`
	Title				string
	Email				string
	Login				string
	Password			string

	CollectionRefer		string
	Collection			Collection	`gorm:"foreignKey:CollectionRefer;constraint:OnDelete:CASCADE;"`
}


type Collection struct {
	gorm.Model
	Id				int				`gorm:"primaryKey"`
	Title			string
	UserRefer		string
	User			User			`gorm:"foreignKey:UserRefer;constraint:OnDelete:CASCADE;"`
}


type Token struct {
	gorm.Model
	Id					int			`gorm:"primaryKey"`
	Token				string
	UserRefer			string
	User				User		`gorm:"foreignKey:UserRefer;constraint:OnDelete:CASCADE;"`
	CreationDate		time.Time
	ExpireDate			time.Time
}


