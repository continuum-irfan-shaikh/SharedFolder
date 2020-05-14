package mec

const (
	HeaderAssetChangeAttributes = "CHANGE-ATTRIBUTES"
	//NewMecTopicName isNew Managed endpoint Change topic for new Kafka message implementation
	HeaderNewAsset string = "NEWASSET"
	//HeaderAssetChange is the kafka header to publish asset change event
	HeaderAssetChange string = "ASSET-CHANGE"
	//HeaderAssetCollectionChange is the kafka header to publish asset collection change event
	HeaderAssetCollectionChange string = "ASSET-COLLECTION-CHANGE"
	//HeaderNetworkIPv4ListChange is the constant for message type in header of new kafka message implementation
	HeaderNetworkIPv4ListChange string = "NETWORKIPV4LISTCHANGE"
	//HeaderChangeFriendlyName is the message type for change of friendly name
	HeaderChangeFriendlyName string = "CHANGEFRIENDLYNAME"
	//HeaderInstalledSoftwareChange is the constant for message type in header of new kafka message implementation
	HeaderInstalledSoftwareChange string = "INSTALLEDSOFTWARECHANGE"
	//HeaderChangeEndpointType is the message type for change of endpoint type
	HeaderChangeEndpointType string = "CHANGEENDPOINTTYPE"
)
