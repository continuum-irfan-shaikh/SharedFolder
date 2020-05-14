package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"gitlab.connectwisedev.com/RMM/rmm-scripts/script-generator/models"

	"github.com/spf13/cobra"
)

var Update = &cobra.Command{
	Use:   "update [path to the updated script|string] [path to existing script definition file (json)|string]",
	Short: "Updates existing JSON definition of script",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			originScriptFile     = args[0]
			scriptDefinitionFile = args[1]
			outputFile           *os.File
			err                  error
		)
		// open for update
		outputFile, err = os.Open(scriptDefinitionFile)
		if err != nil {
			fmt.Println("cannot open script definition file for update: ", err)
			return
		}
		defer outputFile.Close()

		// origin script file
		scriptFile, err := os.Open(originScriptFile)
		if err != nil {
			fmt.Println("cannot open origin script: ", err)
			return
		}
		defer scriptFile.Close()

		scriptContent, err := encodeScriptBody(scriptFile)
		if err != nil {
			fmt.Println("cannot encode script body: ", err)
			return
		}

		data, err := ioutil.ReadFile(scriptDefinitionFile)
		if err != nil {
			fmt.Println("cannot read output file: ", err)
			return
		}

		// recreating output file
		err = os.Remove(scriptDefinitionFile)
		if err != nil {
			fmt.Println("cannot remove old file: ", err)
			return
		}
		outputFile, err = os.Create(scriptDefinitionFile)
		if err != nil {
			fmt.Println("cannot recreate output file: ", err)
			return
		}

		var scriptDefinition models.ScriptDefinition
		err = json.Unmarshal(data, &scriptDefinition)
		if err != nil {
			fmt.Println("cannot unmarshal script definition: ", err)
			return
		}
		scriptDefinition.Content = scriptContent

		rawJSON, err := json.Marshal(scriptDefinition)
		if err != nil {
			fmt.Println("cannot marshal script definition: ", err)
			return
		}

		_, err = outputFile.Write(rawJSON)
		if err != nil {
			fmt.Println("cannot write to output file: ", err)
			return
		}
		fmt.Println(scriptDefinitionFile, " is updated!")
	},
}
