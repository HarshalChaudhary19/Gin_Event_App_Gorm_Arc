package models

import "time"

type FabricUser struct {
	ID           uint   `gorm:"primaryKey"`
	EnrollmentID string `gorm:"uniqueIndex;not null"` // e.g. "user1"
	Org          string `gorm:"not null"`             // e.g. "org1"
	MSPID        string `gorm:"not null"`             // e.g. "Org1MSP"
	Affiliation  string
	CertPEM      []byte `gorm:"type:bytea"`    // certificate PEM
	KeyPEMEnc    []byte `gorm:"type:bytea"`    // encrypted private key PEM should be stored in Encrypted
	IsAdmin      bool   `gorm:"default:false"` // differentiate admin users
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type UserSentFabric struct {
	EnrollId     string `json:"enrollId" binding:"required"`
	EnrollSecret string `json:"enrollsecret" binding:"required"`
	Org          string `json:"org" binding:"required"`
}
