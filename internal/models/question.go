package models

type QuestionLevel string

const (
    Easy   QuestionLevel = "Easy"
    Medium QuestionLevel = "Medium"
    Hard   QuestionLevel = "Hard"
)

type Question struct {
    UUID   string         `gorm:"type:uuid;primaryKey"`
    Text   string         `gorm:"type:text"`
    Answer string         `gorm:"type:text"`
    Level  QuestionLevel  `gorm:"type:varchar(20)"`
}
