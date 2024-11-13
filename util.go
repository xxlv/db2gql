package main

import "strings"

func AsName(name string) string {
	result := strings.Split(asCamStyle(name), "")
	return strings.ToLower(result[0]) + asLowCaseCamStyle(strings.Join(result[1:], ""))
}

// asTypeNameFromKeys [Aa,Bb,Cc] => AaBbCc.
func asTypeNameFromKeys(types map[string]any) string {
	keys := make([]string, 0, len(types))
	for k := range types {
		keys = append(keys, asCamStyle(k))
	}
	return strings.Join(keys, "")
}

func asCamStyle(name string) string {
	if len(name) <= 0 {
		return name
	}

	var result string
	parts := strings.Split(name, "_")
	for i, part := range parts {
		if len(part) > 0 {
			if i > 0 && len(part) == 1 {
				// If it's a single letter, keep it lowercase
				result += part
			} else {
				// Capitalize the first letter of each part
				result += strings.ToUpper(string(part[0])) + strings.ToLower(part[1:])
			}
		}
	}

	return result
}

func asLowCaseCamStyle(name string) string {
	if len(name) <= 0 {
		return name
	}
	nameArr := strings.Split(name, "_")
	for i, part := range nameArr {
		if len(part) > 0 {
			nameArr[i] = strings.ToLower(string(part[0])) + (part[1:])
		}
	}

	return strings.Join(nameArr, "")
}

func asCamStyleWithoutUnderline(name string) string {
	if len(name) <= 0 {
		return name
	}
	nameArr := strings.Split(name, "_")
	for i, part := range nameArr {
		if len(part) > 0 {
			nameArr[i] = (string(part[0])) + (part[1:])
		}
	}
	return strings.Join(nameArr, "")
}
