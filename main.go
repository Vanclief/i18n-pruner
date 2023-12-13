package main

import (
	"encoding/json"
	"fmt"
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/urfave/cli/v2"
)

var duplicatesMap = make(map[interface{}]bool)

func main() {
	app := cli.NewApp()
	app.Name = "pruner"
	app.Usage = "Formating JSON files for i18n"

	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:     "source",
			Aliases:  []string{"s"},
			Usage:    "Source json file path",
			Required: true,
		},
		&cli.StringFlag{
			Name:     "destination",
			Aliases:  []string{"d"},
			Usage:    "Destination json file path",
			Required: true,
		},
		&cli.BoolFlag{
			Name:    "read-only",
			Aliases: []string{"r"},
			Usage:   "Read only",
		},
	}

	app.Action = func(c *cli.Context) error {
		err := format(c.String("source"), c.String("destination"), c.Bool("read-only"))
		if err != nil {
			return cli.Exit(err, 1)
		}
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		os.Exit(0)
	}
}

func format(sourceFilePath, destinationFilePath string, readOnly bool) error {
	sourceFile, err := os.ReadFile(sourceFilePath)
	if err != nil {
		return err
	}

	destinationFile, err := os.ReadFile(destinationFilePath)
	if err != nil {
		return err
	}

	var sourceMap map[string]interface{}
	var destinationMap map[string]interface{}

	api := jsoniter.Config{
		SortMapKeys: true,
	}.Froze()

	err = api.Unmarshal(sourceFile, &sourceMap)
	if err != nil {
		return err
	}

	err = api.Unmarshal(destinationFile, &destinationMap)
	if err != nil {
		return err
	}

	validateAndFillMissing(sourceMap, destinationMap, "")

	sourceJSON, _ := json.MarshalIndent(sourceMap, "", "    ")
	destinationJson, _ := json.MarshalIndent(destinationMap, "", "    ")

	if readOnly {
		return nil
	}

	err = os.WriteFile(sourceFilePath, sourceJSON, os.ModePerm)
	if err != nil {
		return err
	}

	err = os.WriteFile(destinationFilePath, destinationJson, os.ModePerm)
	if err != nil {
		return err
	}

	msg := fmt.Sprintf("Formatted %s and %s successfuly", sourceFilePath, destinationFilePath)
	fmt.Println(msg)

	return nil
}

func validateAndFillMissing(source map[string]interface{}, destination map[string]interface{}, keyPath string) {
	for key, val := range source {
		_, exists := destination[key]
		if !exists {
			if _, isObject := val.(map[string]interface{}); isObject {
				destination[key] = make(map[string]interface{})
			} else {
				destination[key] = ""
			}
		}

		if _, isObject := val.(map[string]interface{}); !isObject {
			if _, exist := duplicatesMap[val]; exist {
				msg := fmt.Sprintf(`Found duplicate at "%s": "%s"`, keyPath, val)
				fmt.Println(msg)
			} else {
				duplicatesMap[val] = true
			}
		}

		nestedSource, sourceOk := val.(map[string]interface{})
		nestedDestination, destinationOk := destination[key].(map[string]interface{})
		if sourceOk && destinationOk {
			if keyPath == "" {
				keyPath = key
			} else {
				keyPath = fmt.Sprintf("%s.%s", keyPath, key)
			}
			validateAndFillMissing(nestedSource, nestedDestination, keyPath)
		}
	}
}
