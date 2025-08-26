package model

import (
    "time"
    "gorm.io/gorm"
)

// Task represents a user task
type Task struct {
    gorm.Model
    UserID      uint      `gorm:"index;not null" json:"user_id"`
    Name        string    `gorm:"not null" json:"name"`
    Subject     string    `json:"subject"`
    Description string    `json:"description"`
    HowItsDone  string    `json:"how_its_done"`
    Completed   bool      `gorm:"not null;default:false" json:"completed"`
    CompletedAt *time.Time `json:"completed_at"`
}

