package types

import (
	"time"
)

type Model struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type Vhost struct {
	Model

	FQDN   string      `json:"fqdn" gorm:"unique;not null"`
	Origin string      `json:"origin" gorm:"not null"`
	Cert   Certificate `json:"certificate" gorm:"foreignkey:VhostID"`
}

// TODO: Add support for password protected keys
type Certificate struct {
	Model

	VhostID uint   `json:"-"`
	RSAKey  []byte `json:"rsa_key" gorm:"type:blob"`
	X509    []byte `json:"certificate" gorm:"type:blob"`
	CAChain []byte `json:"ca_chain" gorm:"type:blob"`
}
