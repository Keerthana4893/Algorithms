package main

import (
	"time"
)
type Prop struct {
	Prop_ID int64  `gorm:"primary_key"`
	Latitude   float32 `gorm:"index:idx_properties_latitude_longitude"`
	Longitude  float32 `gorm:"index:idx_properties_latitude_longitude"`
	Price      float32
	Bedrooms   int16
	Bathrooms  int16
	AddedDate  time.Time
}

func NewProp(lat, lon, price float32, bedrooms, bathrooms int16)*Prop {
	p := Prop{
		Latitude:  lat,
		Longitude: lon,
		Price:     price,
		Bedrooms:  bedrooms,
		Bathrooms: bathrooms,
		AddedDate: time.Now(),
		//we can also use time.NOw().UTC() to get the UTC time. Although, in UNIX() it always returns the number of secs elapsed and hence both time.NOw() and time.Now().UTC() is the same.
	}
	return &p
}

type Req struct {
	ReqID int64  `gorm:"primary_key"`
	Latitude      float32 `gorm:"index:idx_requirements_latitude_longitude"`
	Longitude     float32 `gorm:"index:idx_requirements_latitude_longitude"`
	MinBudget     float32
	MaxBudget     float32
	MinBedrooms   int16
	MaxBedrooms   int16
	MinBathrooms  int16
	MaxBathrooms  int16
	AddedDate     time.Time
}

func NewReq(lat, lon, minBudget, maxBudget float32, minBedrooms, maxBedrooms, minBathrooms, maxBathrooms int16)*Req {
	r := Req{
		Latitude:     lat,
		Longitude:    lon,
		MinBudget:    minBudget,
		MaxBudget:    maxBudget,
		MinBedrooms:  minBedrooms,
		MaxBedrooms:  maxBedrooms,
		MinBathrooms: minBathrooms,
		MaxBathrooms: maxBathrooms,
		AddedDate:    time.Now(),
	}
	return &r
}

/* To find the nearest, the sql query should be (
        IN lat float,  -- Latitude [-90..90]
        IN lon float,  -- Longitude [-180..180]
        IN _START_dist DOUBLE,  -- Starting estimate of how far to search: miles or km
        IN _max_dist DOUBLE,  -- Limit how far to search: miles or km
        IN _limit INT,     -- How many items to try to get
        IN _condition VARCHAR(1111)   -- will be ANDed in a WHERE clause
    )*/
 

func validCoordinate(lat, long float32) bool {
	if lat < -90 || lat > 90 || long < -180 || long > 180 {
		return false
	}
	return true
}

func validBudget(minBudget, maxBudget float32) bool {
	if minBudget < 0 || maxBudget < 0 || minBudget > maxBudget {
		return false
	}
	return true
}

func validBedroomsRange(minRooms, maxRooms int16) bool {
	return validIntRange(minRooms, maxRooms)
}

func validBathroomsRange(minRooms, maxRooms int16) bool {
	return validIntRange(minRooms, maxRooms)
}

func validIntRange(x, y int16) bool {
	if x < 0 || y < 0 || x > y {
		return false
	}
	return true
}

func validPrice(price float32) bool {
	return price > 0
}

func validBedrooms(bed int16) bool {
	return bed > 0
}
func validBathroom(bath int16) bool {
	return bath > 0
}
