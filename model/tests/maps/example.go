package unsupported

type Type struct {
	MapOfStringToString map[string]string
	MapOfMapsToMaps     map[string]map[string]string
	MapOfMapValue       map[string]MapValue
}

type MapValue struct {
	FieldA string
}
