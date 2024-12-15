package utils_test

import (
	"testing"

	"github.com/go-park-mail-ru/2024_2_VKatuny/internal/utils"
	"github.com/stretchr/testify/require"
)

func TestEscapeHTMLStruct(t *testing.T) {
	type work struct {
		Description string
		ISBN  string
	}
	type testingStruct struct {
		Name  string
		Email string
		Age   int
		Notes []string
		Works map[string]*work
	}

	input := &testingStruct{
		Name: "<script>alert('1')</script>Alice",
		Email: "alice@mail.com",
		Age:   23,
		Notes: []string{"<b>Note 1</b>", "Note 2"},
		Works: map[string]*work{
			"<h1>Go: <p>Reflect<p></h1>": {
				Description: "Article about usage reflect package",
				ISBN:  `<img src="hacker-cracker.org">2-266-11156-6`,
			},
			"Alice in Wonderland": {
				Description: "story about Alice",
				ISBN:  "4-833-12256-3",
			},
		},
	}

	expected := &testingStruct{
		Name: "Alice",
		Email: "alice@mail.com",
		Age:   23,
		Notes: []string{"Note 1", "Note 2"},
		Works: map[string]*work{
			"Go: Reflect": {
				Description: "Article about usage reflect package",
				ISBN:  "2-266-11156-6",
			},
			"Alice in Wonderland": {
				Description: "story about Alice",
				ISBN:  "4-833-12256-3",
			},
		},
	}

	utils.EscapeHTMLStruct(input)
	require.EqualValues(t, expected, input)
}
