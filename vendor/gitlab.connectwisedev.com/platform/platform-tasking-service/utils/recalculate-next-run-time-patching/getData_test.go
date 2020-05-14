package main

import (
	"fmt"
	"testing"
)

func TestGet(t *testing.T) {
	got := getTasksData(testTask)
	fmt.Println(got)
}
