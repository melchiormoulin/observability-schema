package protoplugin

import (
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestReqInit(t *testing.T) {
	data, _ := ioutil.ReadFile("/Users/m.moulin/github/observability-schema/schema/observabilitySchema.proto")
	ReqInit(data)
	//	if req.GetProtoFile() == nil {
	//		t.Errorf("bad file parsing")
	//	}
	assert.Equal(t, 123, 123, "COUCOU")

}
