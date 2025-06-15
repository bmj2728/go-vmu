package metadata

type Checker interface {
	Compare() bool
}

type MetaChecker struct {
	ExistingMetadata map[string]interface{}
	CompareMetadata  map[string]interface{}
}

func NewMetaChecker(existing map[string]interface{}, compare map[string]interface{}) *MetaChecker {
	return &MetaChecker{
		ExistingMetadata: existing,
		CompareMetadata:  compare,
	}
}

func (m *MetaChecker) Compare() bool {
	//compare is new data since we want this to be the data we check each value
	for k, v := range m.CompareMetadata {
		//if the key-value is not found or doesn't match we return
		if m.ExistingMetadata[k] != v {
			return false
		}
	}
	//if map processes without a false return true
	return true
}
