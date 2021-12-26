// Code generated by constructor, DO NOT EDIT.

package examples

import uuid "github.com/google/uuid"

func NewUser(id uuid.UUID, name string, age int, socialId string, hobbies []string, languages map[string]bool, bar *Bar, maritalStatus MaritalStatus, permission *Permission, haha map[*int]Bar) *User {
	return &User{
		Age:           age,
		Bar:           bar,
		Haha:          haha,
		Hobbies:       hobbies,
		ID:            id,
		Languages:     languages,
		MaritalStatus: maritalStatus,
		Name:          name,
		Permission:    permission,
		SocialID:      socialId,
	}
}