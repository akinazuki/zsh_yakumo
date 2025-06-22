package defs

import "encoding/json"

func UnmarshalCompletionResponse(data []byte) (CompletionResponse, error) {
	var r CompletionResponse
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *CompletionResponse) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type CompletionResponse struct {
	Model  string  `json:"model"`
	Input  []Input `json:"input"`
	Text   Text    `json:"text"`
	Stream bool    `json:"stream"`
}

type Input struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Text struct {
	Format Format `json:"format"`
}

type Format struct {
	Type   string `json:"type"`
	Name   string `json:"name"`
	Schema Schema `json:"schema"`
}

type Schema struct {
	Type                 string     `json:"type"`
	Properties           Properties `json:"properties"`
	Required             []string   `json:"required"`
	AdditionalProperties bool       `json:"additionalProperties"`
}

type Properties struct {
	Command Command `json:"command"`
}

type Command struct {
	Type string `json:"type"`
}
