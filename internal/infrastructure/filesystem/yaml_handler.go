package filesystem

import (
	"gopkg.in/yaml.v3"
)

func SaveToYAML[T any](data T, filePath string) error {
	/* if err := CreateDirectoryFilePath(filePath); err != nil {
		return err
	} */

	file, err := CreateFile(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	encoder := yaml.NewEncoder(file)
	defer encoder.Close()
	
	return encoder.Encode(data)
}

func LoadFromYAML[T any](filePath string) (T, error) {
	var result T

	file, err := OpenFile(filePath)
	if err != nil {
		return result, err
	}
	defer file.Close()
	
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&result)
	
	return result, err
}

