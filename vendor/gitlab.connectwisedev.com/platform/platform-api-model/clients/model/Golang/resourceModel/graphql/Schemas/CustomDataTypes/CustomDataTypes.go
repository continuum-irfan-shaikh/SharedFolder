package CustomDataTypes

import (
	"time"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

//DateTimeType : Custom date Time Type
var DateTimeType = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "DateTime",
	Description: "DateTime is a DateTime in ISO 8601 format",
	Serialize: func(value interface{}) interface{} {
		switch value := value.(type) {
		case time.Time:
			return value.Format(time.RFC3339)
		}
		//WriteLog(ERROR, "DateTimeType Conversion Got invalid type ["+value.(string)+"]")
		return "INVALID"
	},
	ParseValue: func(value interface{}) interface{} {
		switch tvalue := value.(type) {
		case string:
			var tval time.Time
			var err error
			if tval, err = time.Parse(time.RFC3339, tvalue); err != nil {
				return nil
			}
			return tval

		}
		return nil
	},
	ParseLiteral: func(valueAST ast.Value) interface{} {
		switch valueAST := valueAST.(type) {
		case *ast.StringValue:
			return valueAST.Value
		}
		return nil
	},
},
)
