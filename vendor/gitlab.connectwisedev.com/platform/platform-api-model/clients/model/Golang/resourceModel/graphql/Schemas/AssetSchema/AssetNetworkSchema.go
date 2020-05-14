package AssetSchema

import (
	"time"
	"github.com/graphql-go/graphql"
	"github.com/ContinuumLLC/rmm-device-graphql-service/Schemas/CustomDataTypes"
)

//AssetNetworkData : AssetNetworkData Structure
type AssetNetworkData struct {
	Vendor           	string
	Product          	string
	DhcpEnabled      	bool
	DhcpServer       	string
	DhcpLeaseObtained	time.Time
	DhcpLeaseExpires	time.Time
	DNSServers       	[]string
	IPEnabled		bool
	IPV4             	string
	IPv4List		[]string
	IPV6             	string
	IPv6List		[]string
	SubnetMask       	string
	SubnetMasks		[]string
	DefaultIPGateway 	string
	DefaultIPGateways	[]string
	MacAddress       	string
	WinsPrimaryServer	string
	WinsSecondaryServer	string
	LogicalName      	string
}

//AssetNetworkType : AssetNetwork GraphQL Schema
var AssetNetworkType = graphql.NewObject(graphql.ObjectConfig{
	Name: "assetNetwork",
	Fields: graphql.Fields{
		"vendor": &graphql.Field{
			Type:        graphql.String,
			Description: "Ethernet vendor name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetNetworkData); ok {
					return CurData.Vendor, nil
				}
				return nil, nil
			},
		},

		"product": &graphql.Field{
			Type:        graphql.String,
			Description: "Ethernet product name",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetNetworkData); ok {
					return CurData.Product, nil
				}
				return nil, nil
			},
		},

		"dhcpEnabled": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "DHCP is enabled or not",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetNetworkData); ok {
					return CurData.DhcpEnabled, nil
				}
				return nil, nil
			},
		},

		"dhcpServer": &graphql.Field{
			Type:        graphql.String,
			Description: "DHCP server IP address",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetNetworkData); ok {
					return CurData.DhcpServer, nil
				}
				return nil, nil
			},
		},

		"dhcpLeaseObtained": &graphql.Field{
			Type:        CustomDataTypes.DateTimeType,
			Description: "Date and time the lease was obtained for the IP address assigned to the computer by the DHCP server.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetNetworkData); ok {
					return CurData.DhcpLeaseObtained, nil
				}
				return nil, nil
			},
		},

		"dhcpLeaseExpires": &graphql.Field{
			Type:        CustomDataTypes.DateTimeType,
			Description: "Expiration date and time for a leased IP address that was assigned to the computer by the DHCP server.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetNetworkData); ok {
					return CurData.DhcpLeaseExpires, nil
				}
				return nil, nil
			},
		},

		"dnsServers": &graphql.Field{
			Type:        graphql.NewList(graphql.String),
			Description: "DNS servers IP addresses",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetNetworkData); ok {
					return CurData.DNSServers, nil
				}
				return nil, nil
			},
		},

		"ipEnabled": &graphql.Field{
			Type:        graphql.Boolean,
			Description: "If TRUE, TCP/IP is bound and enabled on this network adapter.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetNetworkData); ok {
					return CurData.IPEnabled, nil
				}
				return nil, nil
			},
		},

		"ipv4": &graphql.Field{
			Type:        graphql.String,
			Description: "IPv4 address",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetNetworkData); ok {
					return CurData.IPV4, nil
				}
				return nil, nil
			},
		},

		"ipv4List": &graphql.Field{
			Type:        graphql.NewList(graphql.String),
			Description: "Array of all of the IPv4 addresses associated with the current network adapter.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetNetworkData); ok {
					return CurData.IPv4List, nil
				}
				return nil, nil
			},
		},

		"ipv6": &graphql.Field{
			Type:        graphql.String,
			Description: "IPv6 address",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetNetworkData); ok {
					return CurData.IPV6, nil
				}
				return nil, nil
			},
		},
		
		"ipv6List": &graphql.Field{
			Type:        graphql.NewList(graphql.String),
			Description: "Array of all of the IPv6 addresses associated with the current network adapter",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetNetworkData); ok {
					return CurData.IPv6List, nil
				}
				return nil, nil
			},
		},

		"subnetMask": &graphql.Field{
			Type:        graphql.String,
			Description: "Subnet mask",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetNetworkData); ok {
					return CurData.SubnetMask, nil
				}
				return nil, nil
			},
		},

		"subnetMasks": &graphql.Field{
			Type:        graphql.NewList(graphql.String),
			Description: "Subnet mask",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetNetworkData); ok {
					return CurData.SubnetMasks, nil
				}
				return nil, nil
			},
		},

		"defaultIPGateway": &graphql.Field{
			Type:        graphql.String,
			Description: "Default gateway IP address",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetNetworkData); ok {
					return CurData.DefaultIPGateway, nil
				}
				return nil, nil
			},
		},

		"defaultIPGateways": &graphql.Field{
			Type:        graphql.NewList(graphql.String),
			Description: "Array of IP addresses of default gateways that the computer system uses",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetNetworkData); ok {
					return CurData.DefaultIPGateways, nil
				}
				return nil, nil
			},
		},

		"macAddress": &graphql.Field{
			Type:        graphql.String,
			Description: "MAC address",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetNetworkData); ok {
					return CurData.MacAddress, nil
				}
				return nil, nil
			},
		},

		"winsPrimaryServer": &graphql.Field{
			Type:        graphql.String,
			Description: "IP address for the primary WINS server.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetNetworkData); ok {
					return CurData.WinsPrimaryServer, nil
				}
				return nil, nil
			},
		},

		"winsSecondaryServer": &graphql.Field{
			Type:        graphql.String,
			Description: "IP address for the secondary WINS server.",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetNetworkData); ok {
					return CurData.WinsSecondaryServer, nil
				}
				return nil, nil
			},
		},

		"logicalName": &graphql.Field{
			Type:        graphql.String,
			Description: "Logical name like eth0 eth1",
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				if CurData, ok := p.Source.(AssetNetworkData); ok {
					return CurData.LogicalName, nil
				}
				return nil, nil
			},
		},
	},
})
