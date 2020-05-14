package site

import (
	"fmt"
	"net/http"
	"testing"

	"gopkg.in/jarcoal/httpmock.v1"
)

func TestNew(t *testing.T) {
	cli := NewSite(http.DefaultClient, "")
	cli.Sites("1", "")

	httpmock.Activate()
	url := fmt.Sprintf("%s/partner/%s/sites", "", "1")

	resp, err := httpmock.NewJsonResponder(http.StatusOK, `{"siteDetailList":{["siteId":1]}}`)
	if err != nil {
		t.Fatal(err)
	}
	httpmock.RegisterResponder(http.MethodGet, url, resp)
	cli.Sites("1", "")

	httpmock.DeactivateAndReset()

	httpmock.Activate()
	resp = httpmock.NewBytesResponder(http.StatusOK, []byte(`{"siteDetailList":[{"siteId":1}]}`))

	httpmock.RegisterResponder(http.MethodGet, url, resp)
	cli.Sites("1", "")
}
