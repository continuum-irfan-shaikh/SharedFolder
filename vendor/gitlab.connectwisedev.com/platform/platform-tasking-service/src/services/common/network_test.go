package common

import (
	"fmt"
	"strings"
	"testing"
)

func TestGetNetworkInterfaces(t *testing.T) {

	tests := []struct {
		testURL string
		want    []string
		err     error
		isIP    bool
	}{
		{
			testURL: "test:1234",
			want:    []string{"test:1234"},
			err:     nil,
		},
		{
			testURL: "localhost8080",
			want:    []string{"localhost:8080"},
			err:     fmt.Errorf("error is needed"),
		},
		{
			testURL: "127.0.0.1",
			want:    []string{"127.0.0.1"},
			err:     fmt.Errorf("error is needed"),
		},
		{ //need to mock getIPs() ([]string, error)
			testURL: ":8888",
			err:     nil,
			isIP:    true,
		},
	}

	for _, test := range tests {

		gotInterfaces, err := GetNetworkInterfaces(test.testURL)
		if err != nil {
			if test.err == nil {
				t.Fatal(err)
			}
			continue
		}

		if test.isIP {
			gotInterfacesLen := len(gotInterfaces)
			if gotInterfacesLen < 1 {
				t.Fatal("Small len() for GetNetworkInterfaces() result, got: ", gotInterfacesLen)
			}

			got := gotInterfaces[0]
			if strings.HasSuffix(got, test.testURL) {
				continue
			} else {
				t.Fatalf("Want %q suffix, but got %q string", test.testURL, got)
			}
		}

		for i, want := range test.want {
			got := gotInterfaces[i]
			if got != want {
				t.Fatalf("Want %q but got %q", want, got)
			}
		}
	}
}

func TestGetIps(t *testing.T) {

	got, err := getIPs()
	if err != nil {
		t.Fatal(err)
	}

	if len(got) < 1 {
		t.Fatal("getIPs: returned empty array")
	}
}
