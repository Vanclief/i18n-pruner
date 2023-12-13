package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/urfave/cli/v2"
	"github.com/vanclief/ez"

	"github.com/sashabaranov/go-openai"
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
		&cli.StringFlag{
			Name:    "translate",
			Aliases: []string{"t"},
			Usage:   "Language to translate to",
		},
		&cli.BoolFlag{
			Name:    "read-only",
			Aliases: []string{"r"},
			Usage:   "Read only",
		},
	}

	app.Action = func(c *cli.Context) error {
		pruner, err := NewPruner(c.String("translate"), c.Bool("read-only"))
		if err != nil {
			return cli.Exit(err, 1)
		}

		err = pruner.Format(c.String("source"), c.String("destination"))
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

type Pruner struct {
	Prompt   string
	ReadOnly bool
	Client   *openai.Client
}

func NewPruner(translateTo string, readOnly bool) (*Pruner, error) {
	const op = "NewPruner"

	p := &Pruner{
		ReadOnly: readOnly,
	}

	if translateTo != "" {
		OPENAI_API_KEY := os.Getenv("OPENAI_API_KEY")
		if OPENAI_API_KEY == "" {
			return nil, ez.New(op, ez.EINVALID, "OPENAI_API_KEY is not set", nil)
		}

		p.Client = openai.NewClient(OPENAI_API_KEY)
		p.Prompt = fmt.Sprintf(`Translate this sentence to %s`, translateTo)
	}

	return p, nil
}

func (p *Pruner) Translate(text any) (string, error) {
	ctx := context.Background()

	req := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo,
		MaxTokens: 40,
		Messages: []openai.ChatCompletionMessage{{
			Role:    "user",
			Content: fmt.Sprintf("%s: %s", p.Prompt, text),
		}},
	}

	resp, err := p.Client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func (p *Pruner) Format(sourceFilePath, destinationFilePath string) error {
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

	p.ValidateAndFillMissing(sourceMap, destinationMap, "")

	sourceJSON, _ := json.MarshalIndent(sourceMap, "", "    ")
	destinationJson, _ := json.MarshalIndent(destinationMap, "", "    ")

	if p.ReadOnly {
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

func (p *Pruner) ValidateAndFillMissing(source map[string]interface{}, destination map[string]interface{}, keyPath string) {
	for key, val := range source {
		_, exists := destination[key]
		if !exists {
			if _, isObject := val.(map[string]interface{}); isObject {
				destination[key] = make(map[string]interface{})
			} else if p.Prompt != "" {
				translation, err := p.Translate(source[key])
				if err != nil {
					fmt.Println("Err:", err)
					destination[key] = ""
				}
				destination[key] = translation
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
			p.ValidateAndFillMissing(nestedSource, nestedDestination, keyPath)
		}
	}
}
