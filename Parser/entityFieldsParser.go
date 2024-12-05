package Parser

import (
	"fmt"
	"regexp"
	"strings"
)

type Field struct {
	Type        string
	MinChar     int
	MaxChar     int
	EnumOptions []string
	Options     []string
}

type Entity struct {
	Name   string
	Fields map[string]Field
}

func parseField(fieldStr string) Field {
	field := Field{}
	re := regexp.MustCompile(`(\w+)\s*:\s*([\w.]+)\s*\{([^}]*)\}`)
	matches := re.FindStringSubmatch(fieldStr)
	if len(matches) == 4 {
		field.Type = matches[2]
		options := strings.Split(matches[3], ",")
		for _, option := range options {
			option = strings.TrimSpace(option)
			if strings.HasPrefix(option, "minChar=") {
				fmt.Sscanf(option, "minChar=%d", &field.MinChar)
			} else if strings.HasPrefix(option, "maxChar=") {
				fmt.Sscanf(option, "maxChar=%d", &field.MaxChar)
			} else if field.Type == "enum" {
				field.EnumOptions = append(field.EnumOptions, option)
			} else {
				field.Options = append(field.Options, option)
			}
		}
	}
	return field
}

func parseEntity(content string) Entity {
	lines := strings.Split(content, "\n")
	entity := Entity{
		Name:   strings.TrimSpace(strings.Split(lines[0], ":")[0]),
		Fields: make(map[string]Field),
	}
	for _, line := range lines[1:] {
		line = strings.TrimSpace(line)
		if line != "" {
			fieldName := strings.TrimSpace(strings.Split(line, ":")[0])
			entity.Fields[fieldName] = parseField(line)
		}
	}
	return entity
}

func ParseFileContents(fileContents []FileContent) (map[string]Entity, []string) {
	entities := make(map[string]Entity)
	order := make([]string, 0)
	for _, fileContent := range fileContents {
		entity := parseEntity(fileContent.Content)
		entities[entity.Name] = entity
		order = append(order, entity.Name)
	}
	return entities, order
}