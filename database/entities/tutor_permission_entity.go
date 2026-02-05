package entities

import (
	"time"

	"gorm.io/gorm"
)

// TutorPermission grants a user (tutor) permission to create questions for a specific Try Out package.
type TutorPermission struct {
	gorm.Model

	TryOutPackageID uint      `json:"tryOutPackageId" gorm:"uniqueIndex:idx_tutor_package;not null"`
	UserID          uint      `json:"userId" gorm:"uniqueIndex:idx_tutor_package;not null"`
	GrantedByUserID uint      `json:"grantedByUserId"`
	GrantedAt       time.Time `json:"grantedAt" gorm:"autoCreateTime"`

	// Relations
	TryOutPackage TryOut `json:"tryOutPackage,omitempty" gorm:"foreignKey:TryOutPackageID"`
	User          User          `json:"user,omitempty" gorm:"foreignKey:UserID"`
	GrantedBy     User          `json:"grantedBy,omitempty" gorm:"foreignKey:GrantedByUserID"`
}
