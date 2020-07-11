package protoplugin

import (
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/assert"
	"testing"
)

func TestTemplatePath(t *testing.T) {
	expectedTemplatePath := "/my/input-template"
	templatePath:=TemplatePath("template_in="+expectedTemplatePath)
	assert.Equal(t,expectedTemplatePath,templatePath,"template path should be the same")
}
func TestTemplatePathDefaultValue(t *testing.T) {
	templatePath1:=TemplatePath("mybadparam")
	assert.Equal(t,templatePath,templatePath1,"template path should be the default one")
}