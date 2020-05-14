package AssetSchema

import (
	"github.com/graphql-go/graphql"
)

//AssetDrivePartition : AssetDrivePartition Structure
type AssetDrivePartition struct {
      Name		string
      Label		string
      FileSystem	string
      Description	string
      SizeBytes		int64
}

//AssetDrivePartitionType : AssetDrivePartition GraphQL Schema
var AssetDrivePartitionType = graphql.NewObject(graphql.ObjectConfig{
	Name: "assetDrivePartition",
	Fields: graphql.Fields{
		"name": &graphql.Field{
			Type:        graphql.String,
			Description: "Name of the partition e.g. C",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetDrivePartition); ok {
					return CurData.Name, nil
				}
				return nil, nil
			},
		},

		"label": &graphql.Field{
			Type:        graphql.String,
			Description: "Label or VolumeName of the partition",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetDrivePartition); ok {
					return CurData.Label, nil
				}
				return nil, nil
			},
		},

		"fileSystem": &graphql.Field{
			Type:        graphql.String,
			Description: "FileSystem of the partition e.g. NTFS",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetDrivePartition); ok {
					return CurData.FileSystem, nil
				}
				return nil, nil
			},
		},

		"description": &graphql.Field{
			Type:        graphql.String,
			Description: "A textual description of the partition",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetDrivePartition); ok {
					return CurData.Description, nil
				}
				return nil, nil
			},
		},

		"sizeBytes": &graphql.Field{
			Type:        graphql.String,
			Description: "Size of the partition in bytes",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetDrivePartition); ok {
					return CurData.SizeBytes, nil
				}
				return nil, nil
			},
		},

		
	},
})


//AssetDriveData : AssetDriveData Structure
type AssetDriveData struct {
	Product      		string
	Manufacturer 		string
	MediaType    		string
	InterfaceType    	string
	LogicalName  		string
	SerialNumber 		string
	Partitions   		[]string
	SizeBytes    		int64
	NumberOfPartitions 	int64
	PartitionData		[]AssetDrivePartition
}

//AssetDriveType : AssetDrive GraphQL Schema
var AssetDriveType = graphql.NewObject(graphql.ObjectConfig{
	Name: "assetDrive",
	Fields: graphql.Fields{
		"product": &graphql.Field{
			Type:        graphql.String,
			Description: "Drive product name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetDriveData); ok {
					return CurData.Product, nil
				}
				return nil, nil
			},
		},

		"manufacturer": &graphql.Field{
			Type:        graphql.String,
			Description: "Drive manufacturer name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetDriveData); ok {
					return CurData.Manufacturer, nil
				}
				return nil, nil
			},
		},

		"mediaType": &graphql.Field{
			Type:        graphql.String,
			Description: "Type of media",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetDriveData); ok {
					return CurData.MediaType, nil
				}
				return nil, nil
			},
		},

		"interfaceType": &graphql.Field{
			Type:        graphql.String,
			Description: "Interface Type of disk for ex. SCSI, SATA",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetDriveData); ok {
					return CurData.InterfaceType, nil
				}
				return nil, nil
			},
		},

		"logicalName": &graphql.Field{
			Type:        graphql.String,
			Description: "Logical name of drive",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetDriveData); ok {
					return CurData.LogicalName, nil
				}
				return nil, nil
			},
		},

		"serialNumber": &graphql.Field{
			Type:        graphql.String,
			Description: "Serial number of drive",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetDriveData); ok {
					return CurData.SerialNumber, nil
				}
				return nil, nil
			},
		},

		"partitions": &graphql.Field{
			Type:        graphql.NewList(graphql.String),
			Description: "Disk partitions information",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetDriveData); ok {
					return CurData.Partitions, nil
				}
				return nil, nil
			},
		},

		"sizeBytes": &graphql.Field{
			Type:        graphql.String,
			Description: "Drive size in bytes",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetDriveData); ok {
					return CurData.SizeBytes, nil
				}
				return nil, nil
			},
		},

		"numberOfPartitions": &graphql.Field{
			Type:        graphql.String,
			Description: "Number of partitions",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetDriveData); ok {
					return CurData.NumberOfPartitions, nil
				}
				return nil, nil
			},
		},

		"partitionData": &graphql.Field{
			Type:        graphql.NewList(AssetDrivePartitionType),
			Description: "Disk partitions Data",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetDriveData); ok {
					return CurData.PartitionData, nil
				}
				return nil, nil
			},
		},
	},
})
