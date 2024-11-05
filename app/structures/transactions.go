package structures

import "strconv"

func Incr(key string) (int, error) {

	item, ok := getMapValue(key)
	if !ok {
		setMapValue(key, MapValue{
			Typ:    "string",
			String: "1",
		})

		return 1, nil
	}

	intVal, err := strconv.Atoi(item.String)
	if err != nil {
		return 0, err
	}

	intVal++

	item.String = strconv.Itoa(intVal)

	setMapValue(key, item)

	return intVal, nil
}
