package cmd

import (
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"gitlab.connectwisedev.com/RMM/rmm-scripts/script-generator/models"

	"encoding/json"
	"github.com/spf13/cobra"
)

var Create = &cobra.Command{
	Use:   "create [script-name|string] [path to script-name.ps1|string]",
	Short: "Creates new JSON template for script for given name",
	Args:  cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			scriptName       = args[0]
			originScriptPath = args[1]
			outputFileName   = scriptName + "/" + scriptName + ".json"
			outputFile       *os.File
			err              error
		)

		err = CreateDirIfNotExist(scriptName)
		if err != nil {
			fmt.Println("cannot open output directory: ", err)
			return
		}
		outputFile, err = os.Create(outputFileName)
		if err != nil {
			fmt.Println("cannot create output file: ", err)
			return
		}
		defer outputFile.Close()

		// origin script file
		scriptFile, err := os.Open(originScriptPath)
		if err != nil {
			fmt.Println("cannot open origin script: ", err)
			return
		}
		defer scriptFile.Close()

		// copy origin script to target destination
		err = CopyFile(scriptName+"/"+originScriptPath, scriptFile)
		if err != nil {
			fmt.Println("cannot copy origin script to destination folder: ", err)
			return
		}

		scriptContent, err := encodeScriptBody(scriptFile)
		if err != nil {
			fmt.Println("cannot encode script body: ", err)
			return
		}

		scriptDefinition := models.NewScriptDefinition()
		scriptDefinition.Category = []string{}
		scriptDefinition.JSONSchema = ""
		scriptDefinition.UISchema = ""
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
		fmt.Println("Script definition is created for: ", scriptName)
	},
}

func encodeScriptBody(file io.Reader) (string, error) {
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}
	encoded := base64.StdEncoding.EncodeToString(bytes)

	return encoded, nil
}

func CreateDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.Mkdir(dir, 0755)
		if err != nil {
			return err
		}
	} else {
		return fmt.Errorf("dir [%s] exists", dir)
	}

	return nil
}

func CopyFile(dstPath string, src *os.File) error {
	defer src.Seek(0, 0)

	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	if err != nil {
		return err
	}
	return dst.Close()
}
