package main

import (
	"log"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	// package errors provide simple error handling...
)
type Prop_list struct {
	Latitude  float32
	Longitude float32
	Price     float32
	Bedrooms  int16
	Bathrooms int16
}

type Req_Dist struct {
	Req
        Dist float32
}
type Prop struct {
	DB             *gorm.DB
	MatchAlgorithm Prop_mat
}
func NewProp(db *gorm.DB, pAlgo Prop_mat) Prop {
	return Prop{
		DB:             db,
		MatchAlgorithm: pAlgo,
	}
}
func (plP Prop) GetMatchingReqs(p Prop_list) ([]Matched_req, error) {
	var err error
	var matchingReqs []Matched_req
	var candidateReqs []Req
	var rMargins ReqMargins
	isValid := plP.validate(p)
	if !isValid {
		return matchingReqs, errors.Wrap(err, "--Prop couldn't validate")
	}
	err = plP.addToDB(p)
	if err != nil {
		return matchingReqs, errors.Wrap(err, "--Prop couldn't add to DB")
	}
	candidateReqs, rMargins, err = plP.getCandidateReqs(p)
	if err != nil {
		return matchingReqs, errors.Wrap(err "--Processor couldn't get Candidate Reqs")
	}
	matchingReqs = plP.MatchAlgorithm.Match(p, candidateReqs, rMargins)
	return matchingReqs, nil
}

func (plP Prop) validate(p Prop_list) bool {
	if !validCoordinate(p.Latitude, p.Longitude) {
		log.Printf("lat: %f or lon: %f", p.Latitude, p.Longitude)
		return false
	}
	if !validPrice(p.Price) {
		log.Printf("price val: %f", p.Price)
		return false
	}
	if !validBedrooms(p.Bedrooms) {
		log.Printf("bed val: %d", p.Bedrooms)
		return false
	}
	if !validBathrooms(p.Bathrooms) {
		log.Printf("bath val: %d", p.Bathrooms)
		return false
	}
	return true
}

func (plP Prop) addToDB(p Prop_list) error {
	newProp := NewProp(p.Latitude, p.Longitude, p.Price, p.Bedrooms, p.Bathrooms)

	err := plP.DB.Debug().Create(newProp).Error
	if err != nil {
		log.Printf("Processor error: (prop: %v, err: %v)", newProp, err)
		return errors.Wrap(err "--Processor error")
	}
	return nil
}
func (plP Prop) getCandidateReqs(p Prop_list) ([]Req, ReqMargins, error) {
	req := []Req
	dist := float32(10) 
	// distance in miles
	rMargins := plP.getReqMargins(p, distanceRange)

	queryString := plP.getQueryString()

	err := plP.DB.Debug().
		Raw(queryString, p.Latitude, p.Latitude, p.Longitude, EarthRadius,
			rMargins.MinLat, rMargins.MaxLat, rMargins.MinLon, rMargins.MaxLon,
			distanceRange, p.Price, rMargins.MinPrice, rMargins.MaxPrice, rMargins.MinPrice, rMargins.MaxPrice,
			p.Bedrooms, rMargins.MinBeds, rMargins.MaxBeds, rMargins.MinBeds, rMargins.MaxBeds,
			p.Bathrooms, rMargins.MinBaths, rMargins.MaxBaths, rMargins.MinBaths, rMargins.MaxBaths).
		Scan(&req).Error
	if err != nil {
		log.Printf("Processor error(property: %v, err: %v)", p, err)
		return req, rMargins, errors.Wrap(err "--Processor error")
	}

	return req, rMargins, nil
}

func (plP Prop) getQueryString() string {
	select_Base := "SELECT req_id, latitude, longitude, "
	select_Dist := "acos(sin(radians(latitude))*sin(radians(?)) + cos(radiand(latitude))*cos(radians(?))*cos(radians(?) - radians(longitude)) ) * ? as distance, "
	select_others := "min_budget, max_budget, min_bedrooms, max_bedrooms, min_bathrooms, max_bathrooms "

	from := "FROM requirements " where

	lat := "latitude BETWEEN ? AND ? "
	lon := "AND longitude BETWEEN ? AND ? "
	dist := "AND distance <= ? "
	price := "AND ((? BETWEEN min_budget AND max_budget) OR (min_budget BETWEEN ? AND ?) OR (max_budget BETWEEN ? AND ?) "
	bed := "AND ((? BETWEEN min_bedrooms AND max_bedrooms) OR (min_bedrooms BETWEEN ? AND ?) OR (max_bedrooms BETWEEN ? AND ?) "
	bath := "AND ((? BETWEEN min_bathrooms AND max_bathrooms) OR (min_bathrooms BETWEEN ? AND ?) OR (max_bathrooms BETWEEN ? AND ?)"

	return select_Base + select_Dist + select_others +
		fromClause + "Where " + lat + lon + dist +
		price + bed + bath
}

func (plP Prop) getReqMargins(p Prop_list, dist float32) ReqMargins {
	minLat, maxLat := GetMinMaxLat(p.Latitude, dist)
	minLon, maxLon := GetMinMaxLon(p.Latitude, p.Longitude, dist)
	minPrice, maxPrice := plP.getMinMaxPrice(p.Price)
	minBeds, maxBeds := plP.getMinMaxBedrooms(p.Bedrooms)
	minBaths, maxBaths := plP.getMinMaxBathrooms(p.Bathrooms)

	return NewReqMargins(minLat, maxLat, minLon, maxLon, minPrice, maxPrice, minBeds, maxBeds, minBaths, maxBaths)
}

func (plP Prop) getMinMaxPrice(price float32) (float32, float32) {
	return MaxF((price - (0.25 * price)), 1.0), MaxF((price + (0.25 * price)), 1.25)
}

func (plP Prop) getMinMaxBedrooms(bedrooms int16) (int16, int16) {
	return Max(bedrooms-2, 1), Max(bedrooms+2, 3)
}

func (plP Prop) getMinMaxBathrooms(bathrooms int16) (int16, int16) {
	return Max(bathrooms-2, 1), Max(bathrooms+2, 3)
}
