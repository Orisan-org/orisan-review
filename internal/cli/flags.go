package cli

type commonFlags struct {
	output            string
	outPath           string
	severityThreshold string
	configPath        string
	maxPatchBytes     int64
	maxFileBytes      int64
	noColor           bool
	failOnReview      bool
	includeGenerated  bool
}

func defaultCommonFlags() commonFlags {
	return commonFlags{
		output:            "table",
		severityThreshold: "high",
		maxPatchBytes:     5_242_880,
		maxFileBytes:      1_048_576,
	}
}
