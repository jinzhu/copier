package domain

import (
	"errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type MaterialGroupType string
type MaterialGroupScope string

var (
	Null                 MaterialGroupType  = ""
	Welcome              MaterialGroupType  = "welcome"
	OrganizationTemplate MaterialGroupType  = "organization_template"
	SelfTemplate         MaterialGroupType  = "self_template"
	BroadcastJob         MaterialGroupType  = "broadcast_job"
	Organization         MaterialGroupScope = "org"
	Group                MaterialGroupScope = "group"
	User                 MaterialGroupScope = "user"
)

var types map[MaterialGroupType]byte
var scopes map[MaterialGroupScope]byte

func init() {
	types := make(map[MaterialGroupType]byte)
	scopes := make(map[MaterialGroupScope]byte)
	types[Null] = 0
	types[Welcome] = 0
	types[OrganizationTemplate] = 0
	types[SelfTemplate] = 0
	types[BroadcastJob] = 0
	scopes[Organization] = 0
	scopes[Group] = 0
	scopes[User] = 0
}

type MaterialGroup struct {
	Id         primitive.ObjectID `bson:"_id,omitempty"`
	Name       string             `bson:"name"`
	Type       MaterialGroupType  `bson:"type,omitempty"`
	Scope      MaterialGroupScope `bson:"scope,omitempty"`
	Order      int64              `bson:"order"`
	CreateTime time.Time          `bson:"createTime,omitempty"`
	UpdateTime time.Time          `bson:"updateTime"`
}

func ParseMaterialGroupType(s string) (*MaterialGroupType, error) {
	t := MaterialGroupType(s)
	if _, exist := types[t]; !exist {
		return nil, errors.New("invalid material group type")
	}
	return &t, nil
}
