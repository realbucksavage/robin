package types

import (
	"time"
)

type Model struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index"`
}

type Host struct {
	Model

	FQDN           string `json:"fqdn" gorm:"unique;not null"`
	Origin         string `json:"origin" gorm:"not null"`
	NickName       string `json:"nick_name"`
	RSAKey         []byte `json:"rsa_key" gorm:"type:blob"`
	SSLCertificate []byte `json:"certificate" gorm:"type:blob"`
}
