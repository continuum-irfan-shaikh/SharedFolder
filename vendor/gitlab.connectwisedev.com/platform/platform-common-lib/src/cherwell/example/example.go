package main

import (
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/cherwell"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/webClient"
)

func main() {
	server, mux := newTestServer()
	defer func() {
		server.CloseClientConnections()
		server.Close()
	}()
	mockExampleHandlers(mux)

	client, err := cherwell.NewClient(cherwell.Config{Host: server.URL}, webClient.ClientFactoryImpl{}.GetClientServiceByType(webClient.BasicClient, webClient.ClientConfig{}))
	if err != nil {
		// handle error
		handleError(err)
	}

	bo := saveExampleCase.bo
	boInfo, err := saveExample(client, bo)
	if err != nil {
		// handle error
		handleError(err)
	}
	handleBOInfo(boInfo)

	respBO, err := getExample(client, boInfo.ID, boInfo.RecordID)
	if err != nil {
		// handle error
		handleError(err)
	}
	handleBO(respBO)

	boInfo, err = deleteExample(client, boInfo.ID, boInfo.RecordID)
	if err != nil {
		// handle error
		handleError(err)
	}
	handleBOInfo(boInfo)
}

func saveExample(client *cherwell.Client, bo cherwell.BusinessObject) (*cherwell.BusinessObjectInfo, error) {
	resp, err := client.Save(bo)
	return resp, err
}

func getExample(client *cherwell.Client, ID, RecordID string) (*cherwell.BusinessObject, error) {
	resp, err := client.GetByRecordID(ID, RecordID)
	return resp, err
}

func deleteExample(client *cherwell.Client, ID, RecordID string) (*cherwell.BusinessObjectInfo, error) {
	resp, err := client.DeleteByRecordID(ID, RecordID)
	return resp, err
}
