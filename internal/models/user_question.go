package models

type UserQuestion struct {
    UUID       string `gorm:"type:uuid;primaryKey"`
    UserID     string `gorm:"type:uuid;index"`
    QuestionID string `gorm:"type:uuid;index"`
    Answered   bool   `gorm:"type:boolean"`
}