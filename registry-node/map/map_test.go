package regmap

import (
	"testing"

	"github.com/brendenehlers/go-distributed-cache/registry-node"
	"github.com/stretchr/testify/assert"
)

func setup() (regmap *RegistryMap) {
	regmap = &RegistryMap{
		values: make(map[string]*registry.RegistryEntry),
	}

	return
}

func TestRegisterNewValue(t *testing.T) {
	regmap := setup()

	expectedUrl := "google.com"
	err := regmap.Register(expectedUrl)
	assert.Nil(t, err)

	value := regmap.values[expectedUrl]
	assert.NotNil(t, value)
	assert.Equal(t, expectedUrl, value.Url)
}

func TestUnregister(t *testing.T) {
	regmap := setup()

	url := "google.com"
	regmap.values[url] = &registry.RegistryEntry{
		Url: url,
	}

	err := regmap.Unregister(url)
	assert.Nil(t, err)

	value, ok := regmap.values[url]
	assert.False(t, ok)
	assert.Nil(t, value)
}

func TestGetNodeFailsWithNoValues(t *testing.T) {
	regmap := setup()

	_, err := regmap.GetNode()

	assert.Error(t, err)
} 

func TestGetNodeFailsWhenNodeIsNil(t *testing.T) {
	regmap := setup()

	regmap.values["google.com"] = nil

	_, err := regmap.GetNode()
	
	assert.Error(t, err)
}

func TestGetNodeReturnsNode(t *testing.T) {
	regmap := setup()

	url := "google.com"
	regmap.values[url] = &registry.RegistryEntry{
		Url: url,
	}

	node, err := regmap.GetNode()

	assert.Nil(t, err)
	assert.NotNil(t, node)
	assert.Equal(t, url, node.Url)
}

func TestGetRandomNodeReturnsNilWhenNoValues(t *testing.T) {
	regmap := setup()

	node := regmap.getRandomNode()

	assert.Nil(t, node)
}

func TestGetRandomNodeReturnsNode(t *testing.T) {
	regmap := setup()

	url := "google.com"
	regmap.values[url] = &registry.RegistryEntry{
		Url: url,
	}

	node := regmap.getRandomNode()

	assert.NotNil(t, node)
	assert.Equal(t, url, node.Url)
}
