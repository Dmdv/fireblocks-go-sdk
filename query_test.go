package fireblocksdk_test

import (
	"encoding/json"
	sdk "fireblocksdk"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type QueryStruct struct {
	Float  float64 `json:"float,omitempty"`
	Int    int64   `json:"integer,omitempty"`
	String string  `json:"string"`
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

func (suite *QuerySuite) TestMustConvertToValueFromStruct() {
	values := sdk.BuildQuery(suite.str)
	require.NotEmpty(suite.T(), values)
	require.Equal(suite.T(), 3, len(values))
}

func (suite *QuerySuite) TestMustConvertUrlValuesFromStruct() {
	values := sdk.BuildQuery(suite.str)
	query := values.UrlValues()

	require.NotEmpty(suite.T(), query)
	require.Equal(suite.T(), 3, len(query))
	require.NotEmpty(suite.T(), query.Get("float"))
	require.NotEmpty(suite.T(), query.Get("integer"))
	require.NotEmpty(suite.T(), query.Get("string"))
}

func (suite *QuerySuite) TestMustConvertToValueFromPointer() {
	values := sdk.BuildQuery(suite.pt)
	require.NotEmpty(suite.T(), values)
	require.Equal(suite.T(), 3, len(values))
}

func (suite *QuerySuite) TestMustConvertUrlValuesFromPointer() {
	values := sdk.BuildQuery(suite.pt)
	query := values.UrlValues()

	require.NotEmpty(suite.T(), query)
	require.Equal(suite.T(), 3, len(query))
	require.NotEmpty(suite.T(), query.Get("float"))
	require.NotEmpty(suite.T(), query.Get("integer"))
	require.NotEmpty(suite.T(), query.Get("string"))
}

func (suite *QuerySuite) TestMarshal() {
	body, err := json.Marshal(suite.pt)
	require.NoError(suite.T(), err)

	target := &QueryStruct{}
	err = json.Unmarshal(body, target)
	require.NoError(suite.T(), err)
}
