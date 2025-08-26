package model

import (
    "time"
    "gorm.io/gorm"
)

// Attendance represents a clock in/out record
type Attendance struct {
    gorm.Model
    UserID   uint       `gorm:"index;not null" json:"user_id"`
    TaskID   *uint      `gorm:"index" json:"task_id"`
    ClockIn  time.Time  `gorm:"not null" json:"clock_in"`
    ClockOut *time.Time `json:"clock_out"`
}

