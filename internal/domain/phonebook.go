package domain


type Phonebook struct {
	Contacts map[string]Contact
}

func NewPhonebook() *Phonebook {
	return &Phonebook{
		Contacts: make(map[string]Contact),
	}
}


// func (p *Phonebook) ValidateContact(contact Contact) error {
// 	if contact.Name == "" || contact.Phone == "" {
// 		return ErrInvalidContact
// 	}
// 	return nil
// }
