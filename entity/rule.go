package entity

import (
	"time"

	"gorm.io/gorm"
)

type Rule struct {
	ID        uint `gorm:"primarykey;autoIncrement"`
	CreatedAt time.Time
	StartTime time.Time `gorm:"not null" json:"start_time"`
	Suspect   string    `gorm:"not null" json:"suspect"`
	EndTime   int64     `gorm:"not null"`
	Domain    string    `gorm:"not null" json:"domain"`
	IP        string    `gorm:"not null"`
	TTL       int       `gorm:"not null" json:"ttl"`
}

func (r *Rule) BeforeSave(tx *gorm.DB) (err error) {
	r.EndTime = getEndTime(r.StartTime, r.TTL)
	return
}

func (r *Rule) BeforeCreate(tx *gorm.DB) (err error) {
	r.ID = 0
	return
}

func getEndTime(startTime time.Time, ttl int) int64{
	return startTime.Add(time.Duration(ttl) * time.Second).Unix()
}
