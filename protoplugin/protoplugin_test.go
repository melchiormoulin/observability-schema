package protoplugin

import (
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/assert"
	"testing"
)

func TestTemplatePath(t *testing.T) {
	expectedTemplatePath := "/my/input-template"
	templatePath,_ := GetParameters("template_in=" + expectedTemplatePath)
	assert.Equal(t, expectedTemplatePath, templatePath, "template path should be the same")
}
func TestTemplatePathBadFormat(t *testing.T) {
	_,err := GetParameters("mybadparam")
	assert.Errorf(t,err,"the template input param is not in the right format")
}
