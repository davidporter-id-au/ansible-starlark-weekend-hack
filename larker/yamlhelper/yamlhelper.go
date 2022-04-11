package yamlhelper

import (
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

const DefaultYamlMarshalIndent = 2
const generationComment = "# Generated code. Do not edit"

func PrettyPrint(node *yaml.Node) (string, error) {
	builder := &strings.Builder{}
	encoder := yaml.NewEncoder(builder)
	encoder.SetIndent(DefaultYamlMarshalIndent)
	err := encoder.Encode(node)
	if err != nil {
		return "", err
	}
	err = encoder.Close()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s\n%s", generationComment, builder.String()), nil
}

func NewSeqNode(content []*yaml.Node) *yaml.Node {
	var result yaml.Node
	result.Kind = yaml.SequenceNode
	result.Tag = "!!seq"
	result.Content = content
	return &result
}

func NewMapNode(content []*yaml.Node) *yaml.Node {
	var result yaml.Node
	result.Kind = yaml.MappingNode
	result.Tag = "!!map"
	result.Content = content
	return &result
}

func NewStringNode(text string) *yaml.Node {
	var result yaml.Node
	result.SetString(text)
	return &result
}
