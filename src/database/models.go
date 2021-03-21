package database

import (
	"gorm.io/gorm"
	"time"
)

var ModelsList = []Model{
	&User{}, &Password{}, &Collection{}, &Pass2Coll{}, &Token{},
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
	Username			string
	Password			string
}


type Collection struct {
	gorm.Model
	Id				int				`gorm:"primaryKey"`
	Title			string
	UserRefer		string
	User			User			`gorm:"foreignKey:UserRefer;constraint:OnDelete:CASCADE;"`
}


type Pass2Coll struct {
	gorm.Model
	Id					int			`gorm:"primaryKey"`
	PasswordRefer		string
	Password			Password	`gorm:"foreignKey:PasswordRefer;constraint:OnDelete:CASCADE;"`
	CollectionRefer		string
	Collection			Collection	`gorm:"foreignKey:CollectionRefer;constraint:OnDelete:CASCADE;"`
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


