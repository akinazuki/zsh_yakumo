package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"zsh_yakumo/defs"

	"github.com/joho/godotenv"
	"github.com/valyala/fastjson"
)

type OpenAIConfig struct {
	Model string
	Token string
}

var openAIConfig = readConfig()
var system_prompt = "You are a zsh shell expert, please help me complete the following command, you should only output the completed command, no need to include any other explanation. Do not put completed command in a code block."
var zsh_prefix = "#!/bin/zsh\n\n"

func readConfig() OpenAIConfig {
	HOME := os.Getenv("HOME")
	configFile := HOME + "/.config/zsh_yakumo.env"

	godotenv.Load(configFile)
	model := os.Getenv("OPENAI_MODEL")

	if model == "" {
		model = "gpt-4o-mini"
	}

	token := os.Getenv("OPENAI_TOKEN")
	if token == "" {
		panic(fmt.Sprintf("Please set OPENAI_TOKEN environment variable in %s", configFile))
	}

	return OpenAIConfig{
		Model: model,
		Token: token,
	}
}

func main() {
	args := os.Args[1:]
	cursor_pointer := args[0]

	stdinBuf, _ := io.ReadAll(os.Stdin)

	bufferPrefix := string(stdinBuf[:len(stdinBuf)-len(cursor_pointer)])
	bufferSuffix := string(stdinBuf[len(stdinBuf)-len(cursor_pointer):])
	fullCommand := zsh_prefix + bufferPrefix + bufferSuffix

	completion := requestOpenAICompletion(fullCommand)
	processedCmd := postProcessCommand(completion, bufferPrefix, bufferSuffix)
	fmt.Println(processedCmd)
}

func postProcessCommand(
	completion string,
	bufferPrefix string,
	bufferSuffix string,
) string {
	completion = strings.TrimLeft(completion, zsh_prefix)

	linePrefix := bufferPrefix
	if idx := strings.LastIndex(bufferPrefix, "\n"); idx != -1 {
		linePrefix = bufferPrefix[idx+1:]
	}

	for _, prefix := range []string{bufferPrefix, linePrefix} {
		if strings.HasPrefix(completion, prefix) {
			completion = completion[len(prefix):]
			break
		}
	}

	if bufferSuffix != "" && strings.HasSuffix(completion, bufferSuffix) {
		completion = completion[:len(completion)-len(bufferSuffix)]
	}

	completion = strings.Trim(completion, "\n")

	if strings.HasPrefix(strings.TrimSpace(linePrefix), "#") {
		completion = "\n" + completion
	}
	return completion
}

func requestOpenAICompletion(fullCommand string) string {
	url := "https://api.openai.com/v1/responses"
	method := "POST"

	payloadJSON := defs.CompletionResponse{
		Model: openAIConfig.Model,
		Input: []defs.Input{
			{
				Role:    "system",
				Content: system_prompt,
			},
			{
				Role:    "user",
				Content: fullCommand,
			},
		},
		Text: defs.Text{
			Format: defs.Format{
				Type: "json_schema",
				Name: "shell_completion",
				Schema: defs.Schema{
					Type: "object",
					Properties: defs.Properties{
						Command: defs.Command{
							Type: "string",
						},
					},
					Required:             []string{"command"},
					AdditionalProperties: false,
				},
			},
		},
		Stream: false,
	}

	payloadBytes, _ := json.Marshal(payloadJSON)
	payload := bytes.NewBuffer(payloadBytes)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return ""
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", openAIConfig.Token))

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	var parser = fastjson.Parser{}
	jsonData, err := parser.ParseBytes(body)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return ""
	}
	command := string(jsonData.GetArray("output")[0].GetArray("content")[0].GetStringBytes("text"))

	var commandStruct struct {
		Command string `json:"command"`
	}
	err = json.Unmarshal([]byte(command), &commandStruct)
	if err != nil {
		fmt.Printf("Error unmarshalling command: %v, raw body: %s\n", err, body)
		return ""
	}
	return commandStruct.Command
}
