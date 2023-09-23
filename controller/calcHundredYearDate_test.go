package controller

import (
	"date_calculation/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCalcHundreYearDate_ErrorHandling(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	router.POST("/api/CalcHundredYearDate", CalcHundreYearDate)

	testCases := []struct {
		name           string
		payload        string
		expectedStatus int
		expectedValues models.OutputResults
	}{
		{
			name:           "Invalid value - mixed characters",
			payload:        `{"date": "abcde3d4"}`,
			expectedStatus: http.StatusBadRequest,
			expectedValues: models.OutputResults{
				AcscEuropean:          "",
				AcscHundredYear:       "",
				AcscInternational:     "",
				AcscJulian:            "",
				AcscUsaStandard:       "",
				DayOfWeek:             "",
				ErrorFlag:             "HTTP 400",
				ErrorText:             "invalid 100 year date: must be a positive number",
				EuropeanStandard:      "",
				InternationalStandard: "",
				UsaStandard:           "",
			},
		},
		{
			name:           "Invalid value - letters",
			payload:        `{"date": "abcde"}`,
			expectedStatus: http.StatusBadRequest,
			expectedValues: models.OutputResults{
				AcscEuropean:          "",
				AcscHundredYear:       "",
				AcscInternational:     "",
				AcscJulian:            "",
				AcscUsaStandard:       "",
				DayOfWeek:             "",
				ErrorFlag:             "HTTP 400",
				ErrorText:             "invalid 100 year date: must be a positive number",
				EuropeanStandard:      "",
				InternationalStandard: "",
				UsaStandard:           "",
			},
		},
		{
			name:           "Invalid value - empty",
			payload:        `{"date": ""}`,
			expectedStatus: http.StatusBadRequest,
			expectedValues: models.OutputResults{
				AcscEuropean:          "",
				AcscHundredYear:       "",
				AcscInternational:     "",
				AcscJulian:            "",
				AcscUsaStandard:       "",
				DayOfWeek:             "",
				ErrorFlag:             "HTTP 400",
				ErrorText:             "invalid 100 year date: empty",
				EuropeanStandard:      "",
				InternationalStandard: "",
				UsaStandard:           "",
			},
		},
		{
			name:           "Invalid value - decimal",
			payload:        `{"date": "3.14"}`,
			expectedStatus: http.StatusBadRequest,
			expectedValues: models.OutputResults{
				AcscEuropean:          "",
				AcscHundredYear:       "",
				AcscInternational:     "",
				AcscJulian:            "",
				AcscUsaStandard:       "",
				DayOfWeek:             "",
				ErrorFlag:             "HTTP 400",
				ErrorText:             "invalid 100 year date: must be a positive number",
				EuropeanStandard:      "",
				InternationalStandard: "",
				UsaStandard:           "",
			},
		},
		{
			name:           "Hundred Year Out of Range - Negative",
			payload:        `{"date": "-1"}`,
			expectedStatus: http.StatusBadRequest,
			expectedValues: models.OutputResults{
				AcscEuropean:          "",
				AcscHundredYear:       "",
				AcscInternational:     "",
				AcscJulian:            "",
				AcscUsaStandard:       "",
				DayOfWeek:             "",
				ErrorFlag:             "HTTP 400",
				ErrorText:             "100 year date out of range: must be between 0 and 99999",
				EuropeanStandard:      "",
				InternationalStandard: "",
				UsaStandard:           "",
			},
		},
		{
			name:           "Hundred Year Out of Range - High",
			payload:        `{"date": "100000"}`,
			expectedStatus: http.StatusBadRequest,
			expectedValues: models.OutputResults{
				AcscEuropean:          "",
				AcscHundredYear:       "",
				AcscInternational:     "",
				AcscJulian:            "",
				AcscUsaStandard:       "",
				DayOfWeek:             "",
				ErrorFlag:             "HTTP 400",
				ErrorText:             "100 year date out of range: must be between 0 and 99999",
				EuropeanStandard:      "",
				InternationalStandard: "",
				UsaStandard:           "",
			},
		},
		{
			name:           "Malformed JSON",
			payload:        `{"date": ""`,
			expectedStatus: http.StatusBadRequest,
			expectedValues: models.OutputResults{
				AcscEuropean:          "",
				AcscHundredYear:       "",
				AcscInternational:     "",
				AcscJulian:            "",
				AcscUsaStandard:       "",
				DayOfWeek:             "",
				ErrorFlag:             "HTTP 400",
				ErrorText:             "unexpected EOF",
				EuropeanStandard:      "",
				InternationalStandard: "",
				UsaStandard:           "",
			},
		},
		{
			name:           "Valid HYD Date - 12345",
			payload:        `{"date": "12345"}`,
			expectedStatus: http.StatusOK,
			expectedValues: models.OutputResults{
				AcscEuropean:          "19.10.33",
				AcscHundredYear:       "12345",
				AcscInternational:     "33-10-19",
				AcscJulian:            "33-292",
				AcscUsaStandard:       "10/19/33",
				DayOfWeek:             "THU.",
				ErrorFlag:             "0",
				ErrorText:             "",
				EuropeanStandard:      "19.10.1933",
				InternationalStandard: "1933-10-19",
				UsaStandard:           "10/19/1933",
			},
		},
		{
			name:           "Valid HYD Date Used by Date Test- 45189",
			payload:        `{"date": "45189"}`,
			expectedStatus: http.StatusOK,
			expectedValues: models.OutputResults{
				AcscEuropean:          "21.09.23",
				AcscHundredYear:       "45189",
				AcscInternational:     "23-09-21",
				AcscJulian:            "23-264",
				AcscUsaStandard:       " 9/21/23",
				DayOfWeek:             "THU.",
				ErrorFlag:             "0",
				ErrorText:             "",
				EuropeanStandard:      "21.09.2023",
				InternationalStandard: "2023-09-21",
				UsaStandard:           " 9/21/2023",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/CalcHundredYearDate", strings.NewReader(tc.payload))
			req.Header.Set("Content-Type", "application/json")

			router.ServeHTTP(w, req)
			assert.Equal(t, tc.expectedStatus, w.Code)

			var responseWrapper ResponseWrapper
			err := json.NewDecoder(w.Body).Decode(&responseWrapper)

			assert.NoError(t, err)
			assert.NotEmpty(t, responseWrapper.Results.ErrorFlag)

			assert.Equal(t, tc.expectedValues.AcscEuropean, responseWrapper.Results.AcscEuropean)
			assert.Equal(t, tc.expectedValues.AcscHundredYear, responseWrapper.Results.AcscHundredYear)
			assert.Equal(t, tc.expectedValues.AcscInternational, responseWrapper.Results.AcscInternational)
			assert.Equal(t, tc.expectedValues.AcscJulian, responseWrapper.Results.AcscJulian)
			assert.Equal(t, tc.expectedValues.AcscUsaStandard, responseWrapper.Results.AcscUsaStandard)
			assert.Equal(t, tc.expectedValues.DayOfWeek, responseWrapper.Results.DayOfWeek)
			assert.Equal(t, tc.expectedValues.ErrorFlag, responseWrapper.Results.ErrorFlag)
			assert.Equal(t, tc.expectedValues.ErrorText, responseWrapper.Results.ErrorText)
			assert.Equal(t, tc.expectedValues.EuropeanStandard, responseWrapper.Results.EuropeanStandard)
			assert.Equal(t, tc.expectedValues.InternationalStandard, responseWrapper.Results.InternationalStandard)
			assert.Equal(t, tc.expectedValues.UsaStandard, responseWrapper.Results.UsaStandard)
		})
	}
}
