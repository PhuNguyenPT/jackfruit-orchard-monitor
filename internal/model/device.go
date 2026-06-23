package model

import "time"

type DeviceStatus struct {
	Connected bool
	UpdatedAt time.Time
}
