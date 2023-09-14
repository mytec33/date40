package models

type OutputResults struct {
	AcscEuropean          string `json:"AcscEuropean"`      // 15.07.23
	AcscHundredYear       string `json:"AcscHundredYear"`   // 4/15/73 -> 26768
	AcscInternational     string `json:"AcscInternational"` // 23-07-15
	AcscJulian            string `json:"AcscJulian"`        // 8/31/2023 -> 23-243
	AcscUsaStandard       string `json:"AcscUsaStandard"`   // 7/15/23
	DayOfWeek             string `json:"DayOfWeek"`         // THU FRI
	ErrorFlag             string `json:"ErrorFlag"`         // ???
	ErrorText             string `json:"ErrorText"`
	EuropeanStandard      string `json:"EuropeanStandard"`      // 15.07.2023
	InternationalStandard string `json:"InternationalStandard"` // 2023-07-15
	UsaStandard           string `json:"UsaStandard"`           // 7/15/2023
}
