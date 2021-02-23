package jsonhal

import (
	"bytes"
	"log"
	"reflect"
	"testing"

	"encoding/json"

	"github.com/stretchr/testify/assert"
)

// HelloWorld is a simple test struct
type HelloWorld struct {
	Hal
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// Foobar is a simple test struct
type Foobar struct {
	Hal
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// Qux is a simple test struct
type Qux struct {
	Hal
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

var expectedJSON = []byte(`{
	"id": 1,
	"name": "Hello World"
}`)

var expectedJSON2 = []byte(`{
	"_links": {
		"self": {
			"href": "/v1/hello/world/1"
		}
	},
	"id": 1,
	"name": "Hello World"
}`)

var expectedJSON3 = []byte(`{
	"_links": {
		"next": {
			"href": "/v1/hello/world?offset=4\u0026limit=2"
		},
		"previous": {
			"href": "/v1/hello/world?offset=0\u0026limit=2"
		},
		"self": {
			"href": "/v1/hello/world?offset=2\u0026limit=2"
		}
	},
	"_embedded": {
		"foobar": {
			"_links": {
				"self": {
					"href": "/v1/foo/bar/1"
				}
			},
			"id": 1,
			"name": "Foo bar 1"
		}
	},
	"id": 1,
	"name": "Hello World"
}`)

var expectedJSON4 = []byte(`{
	"_links": {
		"self": {
			"href": "/v1/hello/world/1"
		}
	},
	"_embedded": {
		"foobars": [
			{
				"_links": {
					"self": {
						"href": "/v1/foo/bar/1"
					}
				},
				"id": 1,
				"name": "Foo bar 1"
			},
			{
				"_links": {
					"self": {
						"href": "/v1/foo/bar/2"
					}
				},
				"id": 2,
				"name": "Foo bar 2"
			}
		]
	},
	"id": 1,
	"name": "Hello World"
}`)

var expectedJSON5 = []byte(`{
	"_links": {
		"self": {
			"href": "/v1/hello/world/1"
		}
	},
	"_embedded": {
		"foobars": [
			{
				"_links": {
					"self": {
						"href": "/v1/foo/bar/1"
					}
				},
				"id": 1,
				"name": "Foo bar 1"
			},
			{
				"_links": {
					"self": {
						"href": "/v1/foo/bar/2"
					}
				},
				"id": 2,
				"name": "Foo bar 2"
			}
		],
		"quxes": [
			{
				"_links": {
					"self": {
						"href": "/v1/qux/1"
					}
				},
				"id": 1,
				"name": "Qux 1"
			},
			{
				"_links": {
					"self": {
						"href": "/v1/qux/2"
					}
				},
				"id": 2,
				"name": "Qux 2"
			}
		]
	},
	"id": 1,
	"name": "Hello World"
}`)

func TestHal(t *testing.T) {
	var (
		helloWorld *HelloWorld
		expected   *bytes.Buffer
		actual     []byte
		err        error
		foobar     *Foobar
		foobars    []*Foobar
		quxes      []*Qux
	)

	// Let's test the simplest scenario without links
	helloWorld = &HelloWorld{ID: 1, Name: "Hello World"}

	expected = bytes.NewBuffer([]byte{})
	err = json.Compact(expected, expectedJSON)
	if err != nil {
		log.Fatal(err)
	}
	actual, err = json.Marshal(helloWorld)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, expected.String(), string(actual))

	// Let's add a self link
	helloWorld.SetLink(
		"self",              // name
		"/v1/hello/world/1", // href
		"",                  // title
	)

	expected = bytes.NewBuffer([]byte{})
	err = json.Compact(expected, expectedJSON2)
	if err != nil {
		log.Fatal(err)
	}
	actual, err = json.Marshal(helloWorld)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, expected.String(), string(actual))

	// Let's add more links and a single embedded resource
	helloWorld = &HelloWorld{ID: 1, Name: "Hello World"}
	helloWorld.SetLink(
		"self", // name
		"/v1/hello/world?offset=2&limit=2", // href
		"", // title
	)
	helloWorld.SetLink(
		"next", // name
		"/v1/hello/world?offset=4&limit=2", // href
		"", // title
	)
	helloWorld.SetLink(
		"previous",                         // name
		"/v1/hello/world?offset=0&limit=2", // href
		"", // title
	)
	foobar = &Foobar{ID: 1, Name: "Foo bar 1"}
	foobar.SetLink("self", "/v1/foo/bar/1", "")
	helloWorld.SetEmbedded("foobar", Embedded(foobar))

	// Assert JSON after marshalling is as expected
	expected = bytes.NewBuffer([]byte{})
	err = json.Compact(expected, expectedJSON3)
	if err != nil {
		log.Fatal(err)
	}
	actual, err = json.Marshal(helloWorld)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, expected.String(), string(actual))

	// Let's test embedded resources
	helloWorld = &HelloWorld{ID: 1, Name: "Hello World"}
	helloWorld.SetLink(
		"self",              // name
		"/v1/hello/world/1", // href
		"",                  // title
	)

	// Add embedded foobars
	foobars = []*Foobar{
		&Foobar{
			Hal: Hal{
				Links: map[string]*Link{
					"self": &Link{Href: "/v1/foo/bar/1"},
				},
			},
			ID:   1,
			Name: "Foo bar 1",
		},
		&Foobar{
			Hal: Hal{
				Links: map[string]*Link{
					"self": &Link{Href: "/v1/foo/bar/2"},
				},
			},
			ID:   2,
			Name: "Foo bar 2",
		},
	}
	helloWorld.SetEmbedded("foobars", Embedded(foobars))

	// Assert JSON after marshalling is as expected
	expected = bytes.NewBuffer([]byte{})
	err = json.Compact(expected, expectedJSON4)
	if err != nil {
		log.Fatal(err)
	}
	actual, err = json.Marshal(helloWorld)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, expected.String(), string(actual))

	// Let's test multiple embedded resources
	helloWorld = &HelloWorld{ID: 1, Name: "Hello World"}
	helloWorld.SetLink(
		"self",              // name
		"/v1/hello/world/1", // href
		"",                  // title
	)

	// Add embedded foobars
	foobars = []*Foobar{
		&Foobar{
			Hal: Hal{
				Links: map[string]*Link{
					"self": &Link{Href: "/v1/foo/bar/1"},
				},
			},
			ID:   1,
			Name: "Foo bar 1",
		},
		&Foobar{
			Hal: Hal{
				Links: map[string]*Link{
					"self": &Link{Href: "/v1/foo/bar/2"},
				},
			},
			ID:   2,
			Name: "Foo bar 2",
		},
	}
	helloWorld.SetEmbedded("foobars", Embedded(foobars))

	// Add embedded quxes
	quxes = []*Qux{
		&Qux{
			Hal: Hal{
				Links: map[string]*Link{
					"self": &Link{Href: "/v1/qux/1"},
				},
			},
			ID:   1,
			Name: "Qux 1",
		},
		&Qux{
			Hal: Hal{
				Links: map[string]*Link{
					"self": &Link{Href: "/v1/qux/2"},
				},
			},
			ID:   2,
			Name: "Qux 2",
		},
	}
	helloWorld.SetEmbedded("quxes", Embedded(quxes))

	// Assert JSON after marshalling is as expected
	expected = bytes.NewBuffer([]byte{})
	err = json.Compact(expected, expectedJSON5)
	if err != nil {
		log.Fatal(err)
	}
	actual, err = json.Marshal(helloWorld)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, expected.String(), string(actual))
}

func TestGetLink(t *testing.T) {
	helloWorld := new(HelloWorld)

	var (
		link *Link
		err  error
	)

	// Test when object has no links
	link, err = helloWorld.GetLink("self")
	assert.Nil(t, link)
	if assert.NotNil(t, err) {
		assert.Equal(t, "Link \"self\" not found", err.Error())
	}

	helloWorld.SetLink(
		"self",              // name
		"/v1/hello/world/1", // href
		"",                  // title
	)

	// Test getting a bogus link
	link, err = helloWorld.GetLink("bogus")
	assert.Nil(t, link)
	if assert.NotNil(t, err) {
		assert.Equal(t, "Link \"bogus\" not found", err.Error())
	}

	// Test getting a valid link
	link, err = helloWorld.GetLink("self")
	assert.Nil(t, err)
	if assert.NotNil(t, link) {
		assert.Equal(t, "/v1/hello/world/1", link.Href)
		assert.Equal(t, "", link.Title)
	}
}

func TestDeleteLink(t *testing.T) {
	helloWorld := new(HelloWorld)
	helloWorld.SetLink(
		"self",              // name
		"/v1/hello/world/1", // href
		"",                  // title
	)
	link, err := helloWorld.GetLink("self")
	assert.NotNil(t, link)
	assert.NoError(t, err)

	helloWorld.DeleteLink("self")
	link, err = helloWorld.GetLink("self")
	assert.Nil(t, link)
	assert.EqualError(t, err, "Link \"self\" not found")
}

func TestGetEmbedded(t *testing.T) {
	helloWorld := new(HelloWorld)

	var (
		embedded Embedded
		err      error
		foobars  []*Foobar
	)

	// Test when object has no embedded resources
	embedded, err = helloWorld.GetEmbedded("foobars")
	assert.Nil(t, embedded)
	if assert.NotNil(t, err) {
		assert.Equal(t, "Embedded \"foobars\" not found", err.Error())
	}

	// Add embedded foobars
	foobars = []*Foobar{
		&Foobar{ID: 1, Name: "Foo bar 1"},
		&Foobar{ID: 2, Name: "Foo bar 2"},
	}
	helloWorld.SetEmbedded("foobars", Embedded(foobars))

	// Test geting bogus embedded resources
	embedded, err = helloWorld.GetEmbedded("bogus")
	assert.Nil(t, embedded)
	if assert.NotNil(t, err) {
		assert.Equal(t, "Embedded \"bogus\" not found", err.Error())
	}

	// Test geting valid embedded resources
	embedded, err = helloWorld.GetEmbedded("foobars")
	assert.Nil(t, err)
	if assert.NotNil(t, embedded) {
		reflectedValue := reflect.ValueOf(embedded)
		expectedType := reflect.SliceOf(reflect.TypeOf(new(Foobar)))
		if assert.Equal(t, expectedType, reflectedValue.Type()) {
			assert.Equal(t, 2, reflectedValue.Len())
		}
	}
}

func TestDeleteEmbedded(t *testing.T) {
	helloWorld := new(HelloWorld)
	var (
		embedded Embedded
		err      error
		foobars  []*Foobar
	)

	// Add embedded foobars
	foobars = []*Foobar{
		&Foobar{ID: 1, Name: "Foo bar 1"},
		&Foobar{ID: 2, Name: "Foo bar 2"},
	}
	helloWorld.SetEmbedded("foobars", Embedded(foobars))

	// Test geting valid embedded resources
	embedded, err = helloWorld.GetEmbedded("foobars")
	assert.NoError(t, err)
	assert.NotNil(t, embedded)

	helloWorld.DeleteEmbedded("foobars")
	embedded, err = helloWorld.GetEmbedded("bogus")
	assert.Nil(t, embedded)
	assert.EqualError(t, err, "Embedded \"bogus\" not found")

}
