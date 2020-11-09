package swag

import (
	"github.com/swaggo/swag"
)

const (
	searchDir        = "./"           //Directory you want to parse
	exclude          = ""             //Exclude directories and files when searching, comma separated
	generalInfo      = "main.go"      //Go file path in which 'swagger general API Info' is written
	propertyStrategy = swag.CamelCase //Property Naming Strategy like snakecase,camelcase,pascalcase
	outputFlag       = false
	output           = "./docs" //Output directory for all the generated files(swagger.json, swagger.yaml and doc.go)
	parseVendor      = false    //Parse go files in 'vendor' folder, disabled by default
	parseDependency  = false    //Parse go files in outside dependency folder, disabled by default
	markdownFiles    = ""       //Parse folder containing markdown files to use as description, disabled by default
	codeExampleFiles = ""       //Parse folder containing code example files to use for the x-codeSamples extension, disabled by default
	parseInternal    = false    //Parse go files in internal packages, disabled by default
	generatedTime    = false    //Generate timestamp at the top of docs.go, disabled by default
	parseDepth       = 100      //Dependency parse depth
)

type Config struct {
	// SearchDir the swag would be parse
	SearchDir string

	// excludes dirs and files in SearchDir,comma separated
	Excludes string

	//OutputFlag represents whether output swagger.json file to local directory
	OutputFlag bool

	// OutputDir represents the output directory for all the generated files
	OutputDir string

	// MainAPIFile the Go file path in which 'swagger general API Info' is written
	MainAPIFile string

	// PropNamingStrategy represents property naming strategy like snakecase,camelcase,pascalcase
	PropNamingStrategy string

	// ParseVendor whether swag should be parse vendor folder
	ParseVendor bool

	// ParseDependencies whether swag should be parse outside dependency folder
	ParseDependency bool

	// ParseInternal whether swag should parse internal packages
	ParseInternal bool

	// MarkdownFilesDir used to find markdownfiles, which can be used for tag descriptions
	MarkdownFilesDir string

	// GeneratedTime whether swag should generate the timestamp at the top of docs.go
	GeneratedTime bool

	// CodeExampleFilesDir used to find code example files, which can be used for x-codeSamples
	CodeExampleFilesDir string

	// ParseDepth dependency parse depth
	ParseDepth int
}

func DefaultConfig() Config {
	return Config{
		SearchDir:           searchDir,
		Excludes:            exclude,
		MainAPIFile:         generalInfo,
		PropNamingStrategy:  propertyStrategy,
		OutputFlag:          outputFlag,
		OutputDir:           output,
		ParseVendor:         parseVendor,
		ParseDependency:     parseDependency,
		MarkdownFilesDir:    markdownFiles,
		ParseInternal:       parseInternal,
		GeneratedTime:       generatedTime,
		CodeExampleFilesDir: codeExampleFiles,
		ParseDepth:          parseDepth,
	}
}
