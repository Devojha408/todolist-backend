package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Todo struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	User       primitive.ObjectID `bson:"user,omitempty"`
	Taskname   string             `json:"taskname,omitempty"`
	Task       string             `json:"task,omitempty"`
	Status     bool               `json:"status,omitempty"`
	Targetdate primitive.DateTime `json:"targetdate",omitempty"`
}
