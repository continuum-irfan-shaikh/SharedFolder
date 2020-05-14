package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

const delimiter = '\b'

func importFromCSV() (err error) {
	cmd := exec.Command("cqlsh", "-e", "use platform_tasking_db; COPY task_definitions (deleted, id, partner_id, origin_id, name, description, type, categories, created_at, created_by, updated_at, updated_by, user_parameters) FROM '"+csvFileName+"' WITH DELIMITER = '\b' AND ESCAPE='\r';")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		log.Printf("ERROR: can not execute command copy from for tasks_by_runtime_mv table, err: %v", err)
		return err
	}
	return
}

func exportToCSV(data []TaskDefinition) (err error) {
	file, err := os.Create(csvFileName)
	if err != nil {
		log.Fatal("Cannot create file", err)
		return
	}
	defer file.Close()

	writer := NewWriter(file)
	defer writer.Flush()

	writer.Comma = delimiter

	for _, td := range data {
		var (
			createdAt string
			updatedAt string
			cat       string
		)

		if !td.CreatedAt.IsZero() {
			createdAt = td.CreatedAt.Format(CassandraTimeFormat)
		}

		if !td.UpdatedAt.IsZero() {
			updatedAt = td.UpdatedAt.Format(CassandraTimeFormat)
		}

		if len(td.Categories) > 0 {
			cat = fmt.Sprintf(`'%s'`, strings.Join(td.Categories, `,`))
			cat = fmt.Sprintf(`{%s}`, cat)
		}

		td.UserParameters = strings.ReplaceAll(td.UserParameters, `\r`, ``)

		var b []byte

		b, err = json.Marshal(td.Description)
		if err != nil {
			log.Printf("Cannot marshal description %v", td.Description)
			return err
		}

		td.Description = strings.ReplaceAll(string(b), `\r`, ``)

		fields := []string{
			"false",
			td.ID.String(),
			td.PartnerID,
			td.OriginID.String(),
			td.Name,
			td.Description,
			td.Type,
			cat,
			createdAt,
			td.CreatedBy,
			updatedAt,
			td.UpdatedBy,
			td.UserParameters,
		}

		err = writer.Write(fields)
		if err != nil {
			log.Printf("Cannot write to file %v", err)
			return
		}
	}
	return
}
