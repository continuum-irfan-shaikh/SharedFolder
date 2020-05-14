package main

import (
	"fmt"

	"gitlab.connectwisedev.com/platform/platform-common-lib/src/checksum"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/downloader"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/downloader/http"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/webClient"
)

func main() {
	service := http.GetDownloader(webClient.TLSClient, webClient.ClientConfig{
		MaxIdleConns:                100,
		MaxIdleConnsPerHost:         10,
		IdleConnTimeoutMinute:       1,
		TimeoutMinute:               1,
		DialTimeoutSecond:           100,
		DialKeepAliveSecond:         100,
		TLSHandshakeTimeoutSecond:   100,
		ExpectContinueTimeoutSecond: 100,
	})

	res := service.Download(&downloader.Config{
		URL:              "http://cdn.itsupport247.net/InstallJunoAgent/Plugin/Windows/platform-installation-manager/1.0.216/platform_installation_manager_windows32_1.0.216.zip",
		DownloadLocation: "/home/juno/Desktop/test",
		FileName:         "platform_installation_manager_windows32_1.0.216.zip",
		TransactionID:    "1",
		CheckSumType:     checksum.NONE,
	})

	if res.Error != nil {
		fmt.Println("Download failure with error : ", res.Error)
		return
	}
	fmt.Println("File successfully download at location")
}
