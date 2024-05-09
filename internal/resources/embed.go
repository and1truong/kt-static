package resources

import (
	_ "embed"
)

//go:embed templates/chapter.tpl
var ChapterTemplateFile string
