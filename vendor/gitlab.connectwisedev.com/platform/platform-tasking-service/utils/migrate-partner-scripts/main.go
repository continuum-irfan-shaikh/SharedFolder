package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	"github.com/json-iterator/go"
	"github.com/urfave/cli"
)

const (
	modeFlag        = "mode"
	hostFlag        = "host"
	exportMode      = "export"
	exportCassandra = "exportCassandra"
	importMode      = "import"
	partnersFlag    = "partners"
	templateIDs     = "templateIDs"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func main() {
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  modeFlag,
			Value: "export",
			Usage: "mode for migration process",
		},
		cli.StringFlag{
			Name:  hostFlag,
			Value: "",
		},
		cli.StringSliceFlag{
			Name:  partnersFlag,
			Usage: "partnerIDs to import scripts",
		},
		cli.StringSliceFlag{
			Name:  templateIDs,
			Value: &cli.StringSlice{powershellTemplateID, bashTemplateID, cmdTemplateID},
			Usage: "partnerIDs to import scripts",
		},
		cli.StringFlag{
			Name:  "file, f",
			Usage: "export to/import from file",
		},
	}

	setUpLogging()
	log.Println("INFO: Logging was set up successfully")

	log.Println("INFO: config was loaded successfully")
	defer func() {
		if MsSQLConn != nil {
			MsSQLConn.Close()
		}
	}()

	app.Action = Migrate

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalf("ERROR: failed to run app, err: %v", err)
	}
	log.Println("INFO: app finished the work successfully")
}

func Migrate(c *cli.Context) (err error) {
	if c.String(modeFlag) == exportMode || c.String(modeFlag) == exportCassandra {
		// loading MSSQL
		if c.String(hostFlag) == "" {
			LoadConfig()
			err = LoadMsSQLDBFromConfig()
		} else {
			err = LoadMsSQLDBFromString(c.String(hostFlag))
		}
		if err != nil {
			log.Fatalf(`Can't connect to MsSQL DB' err: %v`, err)
		}

		var useCassandra bool
		if c.String(modeFlag) == exportCassandra {
			useCassandra = true
			err = LoadCassandra()
			if err != nil {
				log.Fatalf(`Can't connect to Cassandra DB', err: %v`, err)
			}
			log.Println("Cassandra loaded")
		}

		log.Printf("MsSQL loaded, stats: %v", MsSQLConn.Stats())

		//export from MySQL to Cassandra and csv file
		partners := c.StringSlice(partnersFlag)
		if len(partners) == 0 {
			return fmt.Errorf("ERROR: There is no mode:" + c.String(modeFlag))
		}
		templates := c.StringSlice(templateIDs)

		log.Printf("You chose: Partners: %v and Templates: %v\n", partners, templates)
		err = exportScripts(partners, templates, useCassandra)
		if err != nil {
			log.Printf("ERROR: export: %v\n", err)
		}
	} else if c.String(modeFlag) == importMode {
		//import from csv to Cassandra
		err := importFromCSV()
		if err != nil {
			log.Printf("ERROR: import: %v\n", err)
		}
	} else {
		log.Println("ERROR: There is no mode:" + c.String(modeFlag))
	}

	return nil
}

func setUpLogging() {
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	logFile, err := os.Create(currentDir + "/" + time.Now().Truncate(time.Minute).Format(time.RFC3339) + "_migrate_partner_scripts.log")
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
}

func exportScripts(partnerIDs, templateIDs []string, useCassandra bool) (err error) {
	partnerIDsString := fmt.Sprintf(`%s`, strings.Join(partnerIDs, `,`))
	templateIDsString := fmt.Sprintf(`%s`, strings.Join(templateIDs, `,`))

	allScripts, err := getOldScripts(partnerIDsString, templateIDsString)
	if err != nil {
		return
	}
	var taskDefs []TaskDefinition

	for tempID, oldScripts := range allScripts {
		converter, ok := converters[tempID]
		if !ok {
			//TODO: can be some common handler in future
			continue
		}
		defs, err := converter(oldScripts)
		if err != nil {
			return fmt.Errorf("can't convert oldScripts to task definitions : %s", err.Error())
		}
		taskDefs = append(taskDefs, defs...)
	}

	if useCassandra {
		for _, t := range taskDefs {
			err = Upsert(t)
			if err != nil {
				log.Printf("Upsert err :%v", err)
				return
			}
		}
	}

	return exportToCSV(taskDefs)
}

func getOldScripts(partnerIDs, templateIDs string) (scripts map[int][]ScriptMsSQL, err error) {
	ctx := context.Background()

	query := fmt.Sprintf("SELECT ScriptDesc, ScriptName, ScriptData, CreatedBy, TemplateID, MemberID from MSTScript where MemberID IN(%s) AND TemplateID IN (%s)", partnerIDs, templateIDs) // todo remove TOP 10

	rows, err := MsSQLConn.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	scripts = make(map[int][]ScriptMsSQL)
	for rows.Next() {
		var (
			script ScriptMsSQL

			name sql.NullString

			scriptData sql.NullString
			scriptDesc sql.NullString
			createdBy  sql.NullString
			templateID sql.NullInt64
			memberID   sql.NullInt64
		)

		err := rows.Scan(&scriptDesc, &name, &scriptData, &createdBy, &templateID, &memberID)
		if err != nil {
			return nil, err
		}

		script.ScriptXML = scriptData.String
		script.ScriptDesc = scriptDesc.String
		script.CreatedBy = createdBy.String
		script.CreatedAt = time.Now()
		script.MemberID = int(memberID.Int64)
		script.TemplateID = int(templateID.Int64)
		script.ScriptName = name.String
		scripts[script.TemplateID] = append(scripts[script.TemplateID], script)
	}
	log.Printf("Got %v CMD, %v Powershell, %v Bash  scripts to migrate", len(scripts[cmdTemplateIDint]), len(scripts[powershellTemplateIDint]), len(scripts[bashTemplateIDint]))
	return
}

// Upsert creates new Task Definition or updates existed in repository
func Upsert(taskDefinition TaskDefinition) error {
	taskDefinitionFields := []interface{}{
		taskDefinition.ID,
		taskDefinition.PartnerID,
		taskDefinition.OriginID,
		taskDefinition.Name,
		taskDefinition.Description,
		taskDefinition.Type,
		taskDefinition.Categories,
		taskDefinition.CreatedAt,
		taskDefinition.CreatedBy,
		taskDefinition.UpdatedAt,
		taskDefinition.UpdatedBy,
		taskDefinition.UserParameters,
		taskDefinition.Deleted,
	}

	query := "INSERT INTO task_definitions (id, partner_id, origin_id, name, description, type, categories, created_at, created_by, updated_at, updated_by, user_parameters, deleted) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	cassandraQuery := QueryCassandra(query, taskDefinitionFields...)
	return cassandraQuery.Exec()
}
