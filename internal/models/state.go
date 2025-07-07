package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

//nolint:recvcheck
type State struct {
	Name StateName `json:"name"`
	Data StateData `json:"data"`
}

func (s *State) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to assert []byte type")
	}

	return json.Unmarshal(bytes, s)
}

func (s State) Value() (driver.Value, error) {
	return json.Marshal(s)
}

type StateName string

const (
	StateName_Default                  StateName = "Default"
	StateName_WaitingForUserField      StateName = "WaitingForUserField"
	StateName_WaitingForQuestionAnswer StateName = "WaitingForQuestionAnswer"
)

type StateData struct {
	Default                  DefaultStateData                  `json:"default"`
	WaitingForUserField      WaitingForUserFieldStateData      `json:"waiting_for_user_field"`
	WaitingForQuestionAnswer WaitingForQuestionAnswerStateData `json:"waiting_for_question_answer"`
}

type DefaultStateData struct{}

func NewDefaultState() State {
	return State{
		Name: StateName_Default,
		Data: StateData{
			Default: DefaultStateData{},
		},
	}
}

type WaitingForUserFieldStateData struct {
	FieldName UserFieldName `json:"field_name"`
}

type UserFieldName string

const (
	UserFieldNameDisplayName = "DisplayName"
)

func NewWaitingForUserFieldState(fieldName UserFieldName) State {
	return State{
		Name: StateName_WaitingForUserField,
		Data: StateData{
			WaitingForUserField: WaitingForUserFieldStateData{
				FieldName: fieldName,
			},
		},
	}
}

type WaitingForQuestionAnswerStateData struct {
	UserQuestionID string `json:"user_question_id"`
}

func NewWaitingForQuestionAnswerState(userQuestionID string) State {
	return State{
		Name: StateName_WaitingForQuestionAnswer,
		Data: StateData{
			WaitingForQuestionAnswer: WaitingForQuestionAnswerStateData{
				UserQuestionID: userQuestionID,
			},
		},
	}
}
