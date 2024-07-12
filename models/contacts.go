package models

import gonanoid "github.com/matoous/go-nanoid/v2"

type Contact struct {
	Id    string
	Name  string
	Email string
}

func (c Contact) New(name string, email string) Contact {
	return Contact{
		Id:    gonanoid.Must(),
		Name:  name,
		Email: email,
	}
}

type Contacts []Contact

func (c Contacts) Add(contact Contact) Contacts {
	return append(c, contact)
}

type QuickContactAddFormData struct {
	Name  string
	Email string
}

func (c QuickContactAddFormData) Validate() map[string]string {
	errors := map[string]string{}

	if c.Name == "" {
		errors["name"] = "Name is required"
	}

	if c.Email == "" {
		errors["email"] = "Email is required"
	}

	return errors
}
