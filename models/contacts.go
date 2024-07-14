package models

import gonanoid "github.com/matoous/go-nanoid/v2"

type Contact struct {
	Id    string
	Name  string
	Email string
}

func (c Contact) New(name string, email string) Contact {
	return Contact{
		Id:    gonanoid.MustGenerate("0123456789abcdefghijklmnopqrstuvwxyz", 16),
		Name:  name,
		Email: email,
	}
}

type Contacts map[string]Contact

func (c Contacts) New() Contacts {
	return make(Contacts)
}

func (c Contacts) Add(contact Contact) {
	c[contact.Id] = contact
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
