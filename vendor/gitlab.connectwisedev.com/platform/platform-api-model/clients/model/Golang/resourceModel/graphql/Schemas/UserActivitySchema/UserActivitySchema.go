package UserActivitySchema

import "github.com/graphql-go/graphql"

//UserEventResponse : user event response Structure
type UserEventResponse struct {
	Response string `json:"response"`
}

//UserEventResponseType : User event response GraphQL Schema
var UserEventResponseType = graphql.NewObject(graphql.ObjectConfig{
	Name: "UserEventResponse",
	Fields: graphql.Fields{
		"response": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(UserEventResponse); ok {
					return CurData.Response, nil
				}
				return nil, nil
			},
		},
	},
})
