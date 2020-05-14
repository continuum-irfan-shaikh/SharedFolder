package main

import (
	"fmt"

	"gitlab.connectwisedev.com/platform/platform-common-lib/src/runtime/util"
	"gitlab.connectwisedev.com/platform/platform-common-lib/src/wal"
)

func main() {
	fmt.Println("Done", process())
}

func process() error {
	c := wal.NewConfig()
	c.Name = fmt.Sprintf("./wal/%s.wal", util.ProcessName())
	w := wal.Create(c)

	err := w.Write("Test-1")
	if err != nil {
		return err
	}

	err = w.WriteObject("Test-2")
	if err != nil {
		return err
	}

	w.Flush()
	records, err := w.Read(5)
	if err != nil {
		return err
	}

	for _, r := range records {
		fmt.Println("Record ", r.Messages)
		r.Commit()
	}

	return nil
}
