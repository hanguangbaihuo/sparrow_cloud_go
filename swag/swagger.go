package swag

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/swaggo/swag"
)

func Build(config Config) error {
	if _, err := os.Stat(config.SearchDir); os.IsNotExist(err) {
		return fmt.Errorf("dir: %s is not exist", config.SearchDir)
	}

	log.Println("Register swagger docs....")
	p := swag.New(swag.SetMarkdownFileDirectory(config.MarkdownFilesDir),
		swag.SetExcludedDirsAndFiles(config.Excludes),
		swag.SetCodeExamplesDirectory(config.CodeExampleFilesDir))
	p.PropNamingStrategy = config.PropNamingStrategy
	p.ParseVendor = config.ParseVendor
	p.ParseDependency = config.ParseDependency
	p.ParseInternal = config.ParseInternal

	if err := p.ParseAPI(config.SearchDir, config.MainAPIFile, config.ParseDepth); err != nil {
		return err
	}
	swagger := p.GetSwagger()

	b, err := json.MarshalIndent(swagger, "", "    ")
	if err != nil {
		return err
	}
	// TODO: send data to swagger api

	// according to user setting, log out to json file
	if !config.OutputFlag {
		return nil
	}
	log.Println("Generate swagger docs....")

	if err := os.MkdirAll(config.OutputDir, os.ModePerm); err != nil {
		log.Printf("successfully register schema to remote server, but failed to create local file %v\n", err)
		return err
	}

	jsonFileName := filepath.Join(config.OutputDir, "swagger.json")
	err = writeFile(b, jsonFileName)
	if err != nil {
		log.Printf("register schema successful, write local file failed %v\n", err)
		return err
	}
	log.Printf("create swagger.json at %+v", jsonFileName)
	return nil
}

func writeFile(b []byte, file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(b)
	return err
}
