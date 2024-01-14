package backend

import(
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGetTomls(t *testing.T) {
	correctTestData := "testdata"

	correctTomls, err := getTomls(&correctTestData)
	assert.NoError(t, err)
	assert.ElementsMatch(t, correctTomls,  []string{"testdata/.at.toml", "testdata/foldera/.at.toml", "testdata/folderb/folderc/.at.toml"})

	incorrectTestData := "incorrecttestdata"
	_, err = getTomls(&incorrectTestData)
	assert.Error(t, err)
}

func TestTomlsToInputData(t *testing.T) {
	correctTomlPaths := []string{"testdata/.at.toml", "testdata/foldera/.at.toml", "testdata/folderb/folderc/.at.toml"}
	_, err := tomlsToInputData(correctTomlPaths)
	assert.Error(t, err)
}