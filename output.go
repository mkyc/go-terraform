package terra

import "encoding/json"

// OutputAll calls terraform output returns all values as a map.
// If there is error fetching the output, fails the test
func OutputAll(options *Options) (map[string]interface{}, error) {
	return OutputForKeysE(options, nil)
}

// OutputForKeysE calls terraform output for the given key list and returns values as a map.
// The returned values are of type interface{} and need to be type casted as necessary. Refer to output_test.go
func OutputForKeysE(options *Options, keys []string) (map[string]interface{}, error) {
	out, err := RunTerraformCommandAndGetStdoutE(options, FormatArgs(options, "output", "-no-color", "-json")...)
	if err != nil {
		return nil, err
	}

	outputMap := map[string]map[string]interface{}{}
	if err := json.Unmarshal([]byte(out), &outputMap); err != nil {
		return nil, err
	}

	if keys == nil {
		outputKeys := make([]string, 0, len(outputMap))
		for k := range outputMap {
			outputKeys = append(outputKeys, k)
		}
		keys = outputKeys
	}

	resultMap := make(map[string]interface{})
	for _, key := range keys {
		value, containsValue := outputMap[key]["value"]
		if !containsValue {
			return nil, OutputKeyNotFound(key)
		}
		resultMap[key] = value
	}
	return resultMap, nil
}
