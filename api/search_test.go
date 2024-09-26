package api_test

import (
	"testing"

	"github.com/pandatix/go-acm/api"
	"github.com/stretchr/testify/assert"
)

func Test_I_Search(t *testing.T) {
	t.Parallel()

	var tests = map[string]struct {
		Params    *api.SearchParams
		ExpectErr bool
	}{
		"capture-the-flag": {
			Params: &api.SearchParams{
				Request: `"capture the flag"`,
			},
			ExpectErr: false,
		},
	}

	for testname, tt := range tests {
		t.Run(testname, func(t *testing.T) {
			assert := assert.New(t)

			cli := api.NewACMClient()
			_, err := cli.Search(tt.Params)

			if tt.ExpectErr {
				assert.NotNil(err)
			} else {
				assert.Nil(err)
			}
		})
	}
}
