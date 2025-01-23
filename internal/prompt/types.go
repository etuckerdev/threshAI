package prompt

import "encoding/xml"

// TaskPrompt represents a prompt template structure
type TaskPrompt struct {
	XMLName xml.Name `xml:"TaskPrompt"`
	System  string   `xml:"System"`
	Inputs  struct {
		RepoContent string `xml:"RepoContent"`
		UserRequest string `xml:"UserRequest"`
		FocusAreas  string `xml:"FocusAreas"`
	} `xml:"Inputs"`
	ProcessFlow    string `xml:"ProcessFlow"`
	OutputTemplate struct {
		Task struct {
			ProblemStatement string `xml:"ProblemStatement"`
			CodeTargets      struct {
				Files     []string `xml:"File"`
				Functions []string `xml:"Function"`
			} `xml:"CodeTargets"`
			SuccessMetrics struct {
				Performance []string `xml:"Performance"`
				Readability []string `xml:"Readability"`
			} `xml:"SuccessMetrics"`
		} `xml:"Task"`
	} `xml:"OutputTemplate"`
	Example struct {
		RepoContent string `xml:"RepoContent"`
		UserRequest string `xml:"UserRequest"`
		Output      string `xml:"Output"`
	} `xml:"Example"`
}
