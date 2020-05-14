package main

import (
	"fmt"

	"gitlab.connectwisedev.com/platform/platform-common-lib/src/checksum"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/communication/http/client"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/downloader"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/downloader/grab"
)

func main() {
	service := grab.GetDownloader(&client.Config{
		MaxIdleConns:                100,
		MaxIdleConnsPerHost:         10,
		IdleConnTimeoutMinute:       1,
		TimeoutMinute:               1,
		DialTimeoutSecond:           100,
		DialKeepAliveSecond:         100,
		TLSHandshakeTimeoutSecond:   100,
		ExpectContinueTimeoutSecond: 100,
	})

	resp := service.Download(&downloader.Config{
		URL:              "http://cdn.itsupport247.net/InstallJunoAgent/Plugin/Windows/platform-installation-manager/1.0.216/platform_installation_manager_windows32_1.0.216.zip",
		DownloadLocation: "/home/juno/Desktop/test",
		FileName:         "platform_installation_manager_windows32_1.0.216.zip",
		TransactionID:    "1",
		CheckSumType:     checksum.MD5,
	})

	if resp.Error != nil {
		fmt.Printf("Download failure with error : %+v", resp)
		return
	}
	fmt.Println("File successfully download at location")
}
