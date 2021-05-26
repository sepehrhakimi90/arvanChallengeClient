package entity

import (
	"time"
)

type Rule struct {
	ID        uint `gorm:"primarykey;autoIncrement"`
	CreatedAt time.Time
	StartTime time.Time `gorm:"not null" json:"start_time"`
	EndTime   int64 `gorm:"not null"`
	Domain    string `gorm:"not null" json:"domain"`
	TTL       int `gorm:"not null" json:"ttl"`
}
