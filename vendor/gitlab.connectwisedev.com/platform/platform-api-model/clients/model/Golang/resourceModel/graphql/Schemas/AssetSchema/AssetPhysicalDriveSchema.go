package AssetSchema

import (
	"github.com/graphql-go/graphql"
)

//AssetPhysicalDriveData : AssetPhysicalDriveData Structure
type AssetPhysicalDriveData struct {
	Type		string
	PartitionData	[]AssetDrivePartition
}

//AssetPhysicalDriveType : AssetPhysicalDrive GraphQL Schema
var AssetPhysicalDriveType = graphql.NewObject(graphql.ObjectConfig{
	Name: "assetPhysicalDrive",
	Fields: graphql.Fields{
		"type": &graphql.Field{
			Type:        graphql.String,
			Description: "Type of physical drive, it could be RemovableDrive,NetworkDrive,CDDrive, etc.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetPhysicalDriveData); ok {
					return CurData.Type, nil
				}
				return nil, nil
			},
		},

		"partitionData": &graphql.Field{
			Type:        graphql.NewList(AssetDrivePartitionType),
			Description: "Number of partitions of a physical drive",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetPhysicalDriveData); ok {
					return CurData.PartitionData, nil
				}
				return nil, nil
			},
		},
	},
})
