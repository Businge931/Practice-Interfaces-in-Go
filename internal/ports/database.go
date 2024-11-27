package ports

type Database interface {
	Create(location string, data map[string]interface{}) (bool, string)
	Read(location string) (bool, string, map[string]interface{})
	Update(location string, data map[string]interface{}) (bool, string)
	Delete(location string) (bool, string)
}
