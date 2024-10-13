package models

import "time"

type Community struct {
	ID   int64  `json:"ID" db:"community_id"`
	Name string `json:"Name" db:"community_name"`
}

type CommunityDetail struct {
	ID           int64     `json:"ID" db:"community_id"`
	Name         string    `json:"Name" db:"community_name"`
	Introduction string    `json:"introduction,omitempty" db:"introduction"`
	CreateTime   time.Time `json:"CreateTime" db:"create_time"`
}
