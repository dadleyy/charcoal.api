package defs

const (
	// BlueprintDefaultLimit is the default size of a page in a json api response.
	BlueprintDefaultLimit = 100

	// BlueprintMaxLimit is the maximum size of a page in a json api response.
	BlueprintMaxLimit = 500

	// BlueprintMinLimit is the minimum size of a page in a json api response.
	BlueprintMinLimit = 1

	// BlueprintFilterStart defines the prefix used in blueprint filter query params
	BlueprintFilterStart = "filter["

	// BlueprintFilterEnd defines the closing part of a blueprint filter query param
	BlueprintFilterEnd = "]"
)
