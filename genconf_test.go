package genconf

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type TestTarget struct {
	A string
	B string
	C string
	D string
}

type Test struct {
	init   TestTarget
	expect TestTarget
	order  []string
	env    map[string]string
}

func TestLoad(t *testing.T) {
	testLoad(t, Test{
		init: TestTarget{
			C: "baz",
		},
		expect: TestTarget{
			A: "foo",
			B: "bar",
			C: "baz",
			D: "lol",
		},
		order: []string{"env", "type", "color"},
		env:   map[string]string{"env": "prod", "color": "green", "type": "foo"},
	})

	testLoad(t, Test{
		init: TestTarget{
			C: "baz",
		},
		expect: TestTarget{
			A: "foo",
			B: "bar",
			C: "baz",
			D: "bop",
		},
		order: []string{"env", "color", "type"},
		env:   map[string]string{"env": "prod", "color": "green", "type": "foo"},
	})

	testLoad(t, Test{
		init: TestTarget{
			C: "baz",
		},
		expect: TestTarget{
			A: "foo",
			B: "cat",
			C: "baz",
		},
		order: []string{"env", "color", "type"},
		env:   map[string]string{},
	})

	testLoad(t, Test{
		init: TestTarget{
			C: "baz",
		},
		expect: TestTarget{
			A: "wat",
			B: "cat",
			C: "baz",
		},
		order: []string{"env", "color", "type"},
		env:   map[string]string{"env": "qa"},
	})
}

func testLoad(t *testing.T, test Test) {
	err := Load(
		"./fixtures",
		test.order,
		test.env,
		&test.init,
	)
	require.NoError(t, err)
	require.Equal(t, test.init, test.expect)
}
