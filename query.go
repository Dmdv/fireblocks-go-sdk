package fireblocksdk

import (
	"fmt"
	"net/url"
	"reflect"
)

type QueryItems []QueryItem

//go:generate stringer -type=EnvEntry
type QueryItem struct {
	Key, Value string
}

func (e QueryItem) String() string {
	return fmt.Sprintf("%s = %s", e.Key, e.Value)
}

// BuildQuery uses `env` and `envDefault` as tag to bind config to viper bindings
// Example: `env:"USERNAME" envDefault:"admin"`
func BuildQuery(in any) QueryItems {
	var vars QueryItems
	var t = reflect.TypeOf(in)
	var kind = t.Kind()
	if kind == reflect.Ptr {
		t = t.Elem()
		kind = t.Kind()
	}

	if kind == reflect.Struct {
		iterateStructFields(t, in, &vars)
	}

	return vars
}

func (items QueryItems) UrlValues() url.Values {
	values := make(url.Values, len(items))

	for _, val := range items {
		values.Add(val.Key, val.Value)
	}

	return values
}

func iterateStructFields(t reflect.Type, v any, vars *QueryItems) {
	val := reflect.ValueOf(v)
	for i := 0; i < t.NumField(); i++ {
		var field = t.Field(i)
		tag := field.Tag.Get("json")
		value := val.Field(i).Interface()

		entry := QueryItem{
			Key:   tag,
			Value: fmt.Sprintf("%v", value),
		}

		*vars = append(*vars, entry)
	}
}
