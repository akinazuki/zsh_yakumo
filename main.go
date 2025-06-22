package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/akinazuki/zsh_yakumo/internal/defs"
	"github.com/akinazuki/zsh_yakumo/internal/logger"

	"github.com/joho/godotenv"
	"github.com/valyala/fastjson"
)

type OpenAIConfig struct {
	Model string
	Token string
}

var appLogger = logger.NewLogger()
var openAIConfig = readConfig()
var system_prompt = "You are a zsh shell expert, please help me complete the following command, you should only output the completed command, no need to include any other explanation. Do not put completed command in a code block."
var zshPrefix = "#!/bin/zsh\n\n"

func readConfig() OpenAIConfig {
	HOME := os.Getenv("HOME")
	configFile := HOME + "/.config/zsh_yakumo.env"

	if appLogger != nil {
		appLogger.Info("Starting configuration read from %s", configFile)
		appLogger.Info("model: %s", os.Getenv("OPENAI_MODEL"))
		appLogger.Info("token: %s", os.Getenv("OPENAI_TOKEN"))
	}

	err := godotenv.Load(configFile)
	if err != nil && appLogger != nil {
		appLogger.Warn("Failed to load config file %s: %v", configFile, err)
	}

	model := os.Getenv("OPENAI_MODEL")

	if model == "" {
		model = "gpt-4o-mini"
		if appLogger != nil {
			appLogger.Info("Using default OpenAI model: %s", model)
		}
	} else {
		if appLogger != nil {
			appLogger.Info("Using configured OpenAI model: %s", model)
		}
	}

	token := os.Getenv("OPENAI_TOKEN")
	if token == "" {
		if appLogger != nil {
			appLogger.Error("OPENAI_TOKEN not found in environment or config file: %s", configFile)
		}
		panic(fmt.Sprintf("Please set OPENAI_TOKEN environment variable in %s", configFile))
	}

	if appLogger != nil {
		appLogger.Info("Configuration loaded successfully - Model: %s, Token: %s", model, "***REDACTED***")
	}

	return OpenAIConfig{
		Model: model,
		Token: token,
	}
}

func main() {
	appLogger.Info("=== ZSH Yakumo started ===")

	args := os.Args[1:]
	appLogger.Debug("Command line arguments: %v", args)

	if len(args) == 0 {
		appLogger.Error("No cursor pointer argument provided")
		appLogger.Error("Usage: %s <cursor_pointer>", os.Args[0])
		os.Exit(1)
	}

	cursorPointer := args[0]
	appLogger.Debug("Cursor pointer: %s", cursorPointer)

	stdinBuf, err := io.ReadAll(os.Stdin)
	if err != nil {
		appLogger.Error("Failed to read from stdin: %v", err)
		os.Exit(1)
	}
	appLogger.Debug("Read %d bytes from stdin", len(stdinBuf))
	appLogger.Debug("Raw stdin content: %q", string(stdinBuf))

	if len(stdinBuf) < len(cursorPointer) {
		appLogger.Error("Invalid input: stdin buffer shorter than cursor pointer")
		os.Exit(1)
	}

	cursorPosition, _ := strconv.Atoi(cursorPointer)
	bufferPrefix := string(stdinBuf[:cursorPosition])
	bufferSuffix := string(stdinBuf[cursorPosition:])
	fullCommand := zshPrefix + bufferPrefix + bufferSuffix

	appLogger.Info("Buffer analysis - Prefix: %q, Suffix: %q", bufferPrefix, bufferSuffix)
	appLogger.Debug("Full command to be sent: %q", fullCommand)

	appLogger.Info("Requesting OpenAI completion...")
	completion := requestOpenAICompletion(fullCommand)
	appLogger.Info("Received completion: %q", completion)

	appLogger.Info("Post-processing completion...")
	processedCmd := postProcessCommand(completion, bufferPrefix, bufferSuffix)
	appLogger.Info("Final processed command: %q", processedCmd)

	fmt.Println(processedCmd)
	appLogger.Info("=== ZSH Yakumo completed successfully ===")
}

func postProcessCommand(
	completion string,
	bufferPrefix string,
	bufferSuffix string,
) string {
	appLogger.Debug("Post-processing started - Input completion: %q", completion)
	appLogger.Debug("Post-processing - Buffer prefix: %q, Buffer suffix: %q", bufferPrefix, bufferSuffix)

	originalCompletion := completion
	completion = strings.TrimPrefix(completion, zshPrefix)
	if originalCompletion != completion {
		appLogger.Debug("Trimmed zsh prefix from completion")
	}

	linePrefix := bufferPrefix
	if idx := strings.LastIndex(bufferPrefix, "\n"); idx != -1 {
		linePrefix = bufferPrefix[idx+1:]
		appLogger.Debug("Extracted line prefix: %q from buffer prefix", linePrefix)
	} else {
		appLogger.Debug("Using full buffer prefix as line prefix")
	}

	appLogger.Debug("Checking prefixes for removal - Full buffer: %q, Line: %q", bufferPrefix, linePrefix)
	for i, prefix := range []string{bufferPrefix, linePrefix} {
		if strings.HasPrefix(completion, prefix) {
			appLogger.Debug("Removing prefix %d (%q) from completion", i, prefix)
			completion = completion[len(prefix):]
			break
		}
	}

	if bufferSuffix != "" && strings.HasSuffix(completion, bufferSuffix) {
		appLogger.Debug("Removing suffix %q from completion", bufferSuffix)
		completion = completion[:len(completion)-len(bufferSuffix)]
	}

	beforeTrim := completion
	completion = strings.Trim(completion, "\n")
	if beforeTrim != completion {
		appLogger.Debug("Trimmed newlines from completion")
	}

	if strings.HasPrefix(strings.TrimSpace(linePrefix), "#") {
		appLogger.Debug("Line prefix starts with comment, adding newline to completion")
		completion = "\n" + completion
	}

	appLogger.Debug("Post-processing completed - Final result: %q", completion)
	return completion
}

func requestOpenAICompletion(fullCommand string) string {
	url := "https://api.openai.com/v1/responses"
	method := "POST"

	appLogger.Info("Starting OpenAI API request")
	appLogger.Debug("API URL: %s", url)
	appLogger.Debug("HTTP Method: %s", method)
	appLogger.Debug("Using model: %s", openAIConfig.Model)
	appLogger.Debug("Command length: %d characters", len(fullCommand))

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

	appLogger.Debug("Constructed request payload with %d input messages", len(payloadJSON.Input))

	payloadBytes, err := json.Marshal(payloadJSON)
	if err != nil {
		appLogger.Error("Failed to marshal request payload: %v", err)
		return ""
	}
	appLogger.Debug("Request payload size: %d bytes", len(payloadBytes))

	payload := bytes.NewBuffer(payloadBytes)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	appLogger.Debug("Created HTTP client with 30s timeout")

	req, err := http.NewRequest(method, url, payload)
	if err != nil {
		appLogger.Error("Failed to create HTTP request: %v", err)
		return ""
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", openAIConfig.Token))
	appLogger.Debug("Set request headers: Content-Type, Accept, Authorization")

	appLogger.Info("Sending request to OpenAI API...")
	startTime := time.Now()
	res, err := client.Do(req)
	duration := time.Since(startTime)

	if err != nil {
		appLogger.Error("HTTP request failed after %v: %v", duration, err)
		return ""
	}
	defer res.Body.Close()

	appLogger.Info("Received response in %v - Status: %s (%d)", duration, res.Status, res.StatusCode)

	if res.StatusCode != http.StatusOK {
		appLogger.Warn("Non-200 status code received: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		appLogger.Error("Failed to read response body: %v", err)
		return ""
	}
	appLogger.Debug("Response body size: %d bytes", len(body))
	appLogger.Debug("Raw response body: %s", string(body))

	var parser = fastjson.Parser{}
	jsonData, err := parser.ParseBytes(body)
	if err != nil {
		appLogger.Error("Failed to parse JSON response: %v", err)
		appLogger.Error("Raw response: %s", string(body))
		return ""
	}
	appLogger.Debug("Successfully parsed JSON response")

	appLogger.Debug("Extracting command from JSON path: output[0].content[0].text")
	outputArray := jsonData.GetArray("output")
	if len(outputArray) == 0 {
		appLogger.Error("No 'output' array found in response or array is empty")
		return ""
	}

	contentArray := outputArray[0].GetArray("content")
	if len(contentArray) == 0 {
		appLogger.Error("No 'content' array found in output[0] or array is empty")
		return ""
	}

	textBytes := contentArray[0].GetStringBytes("text")
	if textBytes == nil {
		appLogger.Error("No 'text' field found in content[0]")
		return ""
	}

	command := string(textBytes)
	appLogger.Debug("Extracted raw command text: %q", command)

	var commandStruct struct {
		Command string `json:"command"`
	}
	err = json.Unmarshal([]byte(command), &commandStruct)
	if err != nil {
		appLogger.Error("Failed to unmarshal command JSON: %v", err)
		appLogger.Error("Command text was: %q", command)
		appLogger.Error("Full response body: %s", string(body))
		return ""
	}

	appLogger.Info("Successfully extracted command: %q", commandStruct.Command)
	return commandStruct.Command
}
