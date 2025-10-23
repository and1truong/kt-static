package internal

import "strings"

type Translation struct {
	Name  string
	Books []Book
}

type Book struct {
	Number    uint
	Code      string
	Name      string
	Chapters  []string
	Group     string
	Testament string
}

type Chapter struct {
	Number     int
	Blocks     []Block
	AudioLinks []string
}

type Block struct {
	Kind       string
	Number     string
	Content    []InnerBlock
	Classes    []string
	References []string // …
	NewLine    bool
}

type InnerBlock struct {
	Content    string
	References []string
}

func (b Block) String() string {
	out := ""

	if len(b.Content) == 0 {
		return out
	}

	if b.Kind == "title" {
		out += "\n## " + b.Content[0].Content + "\n"
	} else {
		num := b.Number
		num = strings.ReplaceAll(num, "0", "⁰")
		num = strings.ReplaceAll(num, "1", "¹")
		num = strings.ReplaceAll(num, "2", "²")
		num = strings.ReplaceAll(num, "3", "³")
		num = strings.ReplaceAll(num, "4", "⁴")
		num = strings.ReplaceAll(num, "5", "⁵")
		num = strings.ReplaceAll(num, "6", "⁶")
		num = strings.ReplaceAll(num, "7", "⁷")
		num = strings.ReplaceAll(num, "8", "⁸")
		num = strings.ReplaceAll(num, "9", "⁹")
		out += " " + num + " "

		parts := make([]string, len(b.Content))
		for i, part := range b.Content {
			parts[i] += part.Content + " "

			if false && len(part.References) > 0 {
				parts[i] += " ⚓ "
				// parts[i] += "\n\n    ⚓ " + strings.Join(part.References, "; ") + "\n\n"
			}
		}

		out += strings.Join(parts, "\n")

		if b.NewLine {
			out += "\n\n"
		} else {
			out += " "
		}
	}

	return out
}
