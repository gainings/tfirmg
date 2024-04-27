package resourceid

type Generator interface {
	Generate(attributes map[string]interface{}) *ResourceID
}

type GeneratorFactory interface {
	Generator(typeName string) Generator
}
type ResourceID string

func (rID ResourceID) String() string {
	return string(rID)
}
