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
func NewProp(db *gorm.DB, proAlgo Prop_mat) Prop {
	return Prop{
		DB:             db,
		MatchAlgorithm: proAlgo,
	}
}
func (plP Prop) GetMatchingReqs(pro Prop_list) ([]Matched_req, error) {
	var err error
	var matchingReqs []Matched_req
	var candidateReqs []Req
	var rMargins ReqMargins
	isValid := plP.validate(pro)
	if !isValid {
		return matchingReqs, errors.Wrap(err, "--Prop couldn't validate")
	}
	err = plP.addToDB(pro)
	if err != nil {
		return matchingReqs, errors.Wrap(err, "--Prop couldn't add to DB")
	}
	candidateReqs, rMargins, err = plP.getCandidateReqs(pro)
	if err != nil {
		return matchingReqs, errors.Wrap(err "--Processor couldn't get Candidate Reqs")
	}
	matchingReqs = plP.MatchAlgorithm.Match(pro, candidateReqs, rMargins)
	return matchingReqs, nil
}

func (plP Prop) validate(pro Prop_list) bool {
	if !validCoordinate(pro.Latitude, pro.Longitude) {
		log.Printf("lat: %f or lon: %f", pro.Latitude, pro.Longitude)
		return false
	}
	if !validPrice(pro.Price) {
		log.Printf("price val: %f", pro.Price)
		return false
	}
	if !validBedrooms(pro.Bedrooms) {
		log.Printf("bed val: %d", pro.Bedrooms)
		return false
	}
	if !validBathrooms(pro.Bathrooms) {
		log.Printf("bath val: %d", pro.Bathrooms)
		return false
	}
	return true
}

func (plP Prop) addToDB(pro Prop_list) error {
	newProp := NewProp(pro.Latitude, pro.Longitude, pro.Price, pro.Bedrooms, pro.Bathrooms)

	err := plP.DB.Debug().Create(newProp).Error
	if err != nil {
		log.Printf("Processor error: (prop: %v, err: %v)", newProp, err)
		return errors.Wrap(err "--Processor error")
	}
	return nil
}
func (plP Prop) getCandidateReqs(pro Prop_list) ([]Req, ReqMargins, error) {
	req := []Req
	dist := float32(10) 
	// distance in miles
	rMargins := plP.getReqMargins(pro, distanceRange)

	queryString := plP.getQueryString()

	err := plP.DB.Debug().
		Raw(queryString, pro.Latitude, pro.Latitude, pro.Longitude, EarthRadius,
			rMargins.MinLat, rMargins.MaxLat, rMargins.MinLon, rMargins.MaxLon,
			distanceRange, pro.Price, rMargins.MinPrice, rMargins.MaxPrice, rMargins.MinPrice, rMargins.MaxPrice,
			pro.Bedrooms, rMargins.MinBeds, rMargins.MaxBeds, rMargins.MinBeds, rMargins.MaxBeds,
			pro.Bathrooms, rMargins.MinBaths, rMargins.MaxBaths, rMargins.MinBaths, rMargins.MaxBaths).
		Scan(&req).Error
	if err != nil {
		log.Printf("Processor error(property: %v, err: %v)", pro, err)
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

func (plP Prop) getReqMargins(pro Prop_list, dist float32) ReqMargins {
	minLat, maxLat := GetMinMaxLat(pro.Latitude, dist)
	minLon, maxLon := GetMinMaxLon(pro.Latitude, pro.Longitude, dist)
	minPrice, maxPrice := plP.getMinMaxPrice(pro.Price)
	minBeds, maxBeds := plP.getMinMaxBedrooms(pro.Bedrooms)
	minBaths, maxBaths := plP.getMinMaxBathrooms(pro.Bathrooms)

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
