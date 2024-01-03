package main

import (
	"html/template"
	"os"
)

type User struct {
	FirstName  string
	SecondName string
	Bio        string
	Age        int
	Working    bool
	Pi         float64
	Pet        []Pet
}

type Pet struct {
	Name   string
	Sex    string
	Intact bool
	Age    string
	Breed  string
}

func main() {
	t, err := template.ParseFiles("hello.gohtml")
	if err != nil {
		panic(err)
	}

	petData := []Pet{
		{
			Name:   "Fluffy",
			Sex:    "Male",
			Intact: true,
			Age:    "2 Months old",
			Breed:  "Labrador",
		},
		{
			Name:   "Fido",
			Sex:    "Female",
			Intact: false,
			Age:    "4 Years old",
			Breed:  "Persian",
		},
		{
			Name:   "Kiera",
			Sex:    "Female",
			Intact: true,
			Age:    "6 Years Old",
			Breed:  "French Bulldog",
		},
	}

	user := User{
		FirstName:  "Jamie",
		SecondName: "Smith",
		Bio:        "I live in the United Kingdom",
		Age:        123,
		Working:    true,
		Pi:         3.14,
		Pet:        petData,
	}

	err = t.Execute(os.Stdout, user)
	if err != nil {
		panic(err)
	}
}
