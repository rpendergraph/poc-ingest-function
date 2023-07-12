package moveit

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getTestDataPath() string {
	_, me, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}
	return fmt.Sprintf("%s/testdata.txt", filepath.Dir(me))
}

func TestParseBody(t *testing.T) {
	contents, err := os.ReadFile(getTestDataPath())
	assert.Nil(t, err)
	events, err := parseMoveItBody(contents)
	assert.Nil(t, err)
	assert.Len(t, events, 9)

}

func TestConvertToIndexWorks(t *testing.T) {
	contents, err := os.ReadFile(getTestDataPath())
	assert.Nil(t, err)
	events, err := parseMoveItBody(contents)
	assert.Nil(t, err)
	indexed := handler{}.createIndexedDocuments(events)
	assert.Len(t, indexed, len(events))
}
