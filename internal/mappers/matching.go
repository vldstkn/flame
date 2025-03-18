package mappers

import (
	"flame/internal/models"
	"strconv"
)

func FromModelToMapMatchingUser(user models.GetMatchingUser) map[string]interface{} {
	var birthdate string
	var photoId int64
	var photoUrl string
	var gender string
	var city string

	if user.BirthDate != nil {
		birthdate = *user.BirthDate
	}
	if user.PhotoId != nil {
		photoId = *user.PhotoId
	}
	if user.PhotoUrl != nil {
		photoUrl = *user.PhotoUrl
	}
	if user.Gender != nil {
		gender = *user.Gender
	}
	if user.City != nil {
		city = *user.City
	}
	return map[string]interface{}{
		"id":         user.Id,
		"name":       user.Name,
		"birth_date": birthdate,
		"photo_id":   photoId,
		"photo_url":  photoUrl,
		"gender":     gender,
		"city":       city,
		"lon":        user.Lon,
		"lat":        user.Lat,
	}
}

func FromMapToModelMatchingUser(user map[string]string) models.GetMatchingUser {
	var id int64
	var photoId *int64

	birthdate := user["birth_date"]

	photoUrl := user["photo_url"]
	var gender *string
	if user["gender"] == "" {
		gender = nil
	} else {
		g := user["gender"]
		gender = &g
	}
	var city *string
	if user["city"] == "" {
		city = nil
	} else {
		c := user["city"]
		city = &c
	}

	photoIdStr := user["photo_id"]
	i, _ := strconv.Atoi(photoIdStr)
	ii := int64(i)
	photoId = &ii

	idStr := user["id"]
	i, _ = strconv.Atoi(idStr)
	id = int64(i)

	lonStr := user["lon"]
	lon, _ := strconv.ParseFloat(lonStr, 64)

	latStr := user["lat"]
	lat, _ := strconv.ParseFloat(latStr, 64)

	return models.GetMatchingUser{
		User: models.User{
			Id:        id,
			BirthDate: &birthdate,
			Name:      user["name"],
			Gender:    gender,
			City:      city,
		},
		PhotoId:  photoId,
		PhotoUrl: &photoUrl,
		Lat:      lat,
		Lon:      lon,
	}
}
