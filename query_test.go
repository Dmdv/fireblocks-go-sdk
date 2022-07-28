package fireblocksdk_test

import (
	sdk "fireblocksdk"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type QueryStruct struct {
	Float  float64 `json:"float,omitempty"`
	Int    int64   `json:"integer,omitempty"`
	String string  `json:"string,omitempty"`
}

type QuerySuite struct {
	suite.Suite
	str QueryStruct
	pt  *QueryStruct
}

func TestQuerySuite(t *testing.T) {
	suite.Run(t, new(QuerySuite))
}

func (suite *QuerySuite) SetupTest() {
	value := QueryStruct{
		Float:  1.01,
		Int:    11,
		String: "name",
	}

	suite.str = value
	suite.pt = &value
}

func (suite *QuerySuite) TestMustConvertToValue() {
	values := sdk.BuildQuery(suite.str)
	require.NotEmpty(suite.T(), values)
	require.Equal(suite.T(), 3, len(values))
}
