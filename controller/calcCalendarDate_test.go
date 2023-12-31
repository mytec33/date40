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

func TestCalcHundreYearDate_ValidValues(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	router.POST("/api/CalcCalendarDate", CalcCalendarDate)

	testCases := []struct {
		name           string
		payload        string
		expectedStatus int
		expectedValues models.OutputResults
	}{
		{
			name:           "First day and month of a year",
			payload:        `{"date": "1/1/2023"}`,
			expectedStatus: http.StatusOK,
			expectedValues: models.OutputResults{
				AcscEuropean:          "01.01.23",
				AcscHundredYear:       "44926",
				AcscInternational:     "23-01-01",
				AcscJulian:            "23-001",
				AcscUsaStandard:       "  1/1/23",
				DayOfWeek:             "SUN.",
				ErrorFlag:             "0",
				ErrorText:             "",
				EuropeanStandard:      "01.01.2023",
				InternationalStandard: "2023-01-01",
				UsaStandard:           "  1/1/2023",
			},
		},
		{
			name:           "Last day and month of a year",
			payload:        `{"date": "12/31/1945"}`,
			expectedStatus: http.StatusOK,
			expectedValues: models.OutputResults{
				AcscEuropean:          "31.12.45",
				AcscHundredYear:       "16801",
				AcscInternational:     "45-12-31",
				AcscJulian:            "45-365",
				AcscUsaStandard:       "12/31/45",
				DayOfWeek:             "MON.",
				ErrorFlag:             "0",
				ErrorText:             "",
				EuropeanStandard:      "31.12.1945",
				InternationalStandard: "1945-12-31",
				UsaStandard:           "12/31/1945",
			},
		},
		{
			name:           "Leap Year",
			payload:        `{"date": "2/29/2020"}`,
			expectedStatus: http.StatusOK,
			expectedValues: models.OutputResults{
				AcscEuropean:          "29.02.20",
				AcscHundredYear:       "43889",
				AcscInternational:     "20-02-29",
				AcscJulian:            "20-060",
				AcscUsaStandard:       " 2/29/20",
				DayOfWeek:             "SAT.",
				ErrorFlag:             "0",
				ErrorText:             "",
				EuropeanStandard:      "29.02.2020",
				InternationalStandard: "2020-02-29",
				UsaStandard:           " 2/29/2020",
			},
		},
		{
			name:           "First Calendar in HYD",
			payload:        `{"date": "1/1/1900"}`,
			expectedStatus: http.StatusOK,
			expectedValues: models.OutputResults{
				AcscEuropean:          "01.01.00",
				AcscHundredYear:       "1",
				AcscInternational:     "00-01-01",
				AcscJulian:            "00-001",
				AcscUsaStandard:       "  1/1/00",
				DayOfWeek:             "MON.",
				ErrorFlag:             "0",
				ErrorText:             "",
				EuropeanStandard:      "01.01.1900",
				InternationalStandard: "1900-01-01",
				UsaStandard:           "  1/1/1900",
			},
		},
		{
			name:           "LAST Calendar in HYD",
			payload:        `{"date": "10/14/2173"}`,
			expectedStatus: http.StatusOK,
			expectedValues: models.OutputResults{
				AcscEuropean:          "14.10.73",
				AcscHundredYear:       "99999",
				AcscInternational:     "73-10-14",
				AcscJulian:            "73-287",
				AcscUsaStandard:       "10/14/73",
				DayOfWeek:             "THU.",
				ErrorFlag:             "0",
				ErrorText:             "",
				EuropeanStandard:      "14.10.2173",
				InternationalStandard: "2173-10-14",
				UsaStandard:           "10/14/2173",
			},
		},
		{
			name:           "Valid Date Used By HYD Test - 9/21/2023",
			payload:        `{"date": "9/21/2023"}`,
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
		{
			name:           "Invalid month",
			payload:        `{"date": "99/14/2173"}`,
			expectedStatus: http.StatusBadRequest,
			expectedValues: models.OutputResults{
				AcscEuropean:          "",
				AcscHundredYear:       "",
				AcscInternational:     "",
				AcscJulian:            "",
				AcscUsaStandard:       "",
				DayOfWeek:             "",
				ErrorFlag:             "HTTP 400",
				ErrorText:             "invalid date: 99/14/2173",
				EuropeanStandard:      "",
				InternationalStandard: "",
				UsaStandard:           "",
			},
		},
		{
			name:           "Invalid day",
			payload:        `{"date": "1/100/2173"}`,
			expectedStatus: http.StatusBadRequest,
			expectedValues: models.OutputResults{
				AcscEuropean:          "",
				AcscHundredYear:       "",
				AcscInternational:     "",
				AcscJulian:            "",
				AcscUsaStandard:       "",
				DayOfWeek:             "",
				ErrorFlag:             "HTTP 400",
				ErrorText:             "invalid date: 1/100/2173",
				EuropeanStandard:      "",
				InternationalStandard: "",
				UsaStandard:           "",
			},
		},
		{
			name:           "Invalid year",
			payload:        `{"date": "1/1/99999"}`,
			expectedStatus: http.StatusBadRequest,
			expectedValues: models.OutputResults{
				AcscEuropean:          "",
				AcscHundredYear:       "",
				AcscInternational:     "",
				AcscJulian:            "",
				AcscUsaStandard:       "",
				DayOfWeek:             "",
				ErrorFlag:             "HTTP 400",
				ErrorText:             "invalid date: 1/1/99999",
				EuropeanStandard:      "",
				InternationalStandard: "",
				UsaStandard:           "",
			},
		},
		{
			name:           "Missing date",
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
				ErrorText:             "invalid date: empty",
				EuropeanStandard:      "",
				InternationalStandard: "",
				UsaStandard:           "",
			},
		},
		{
			name:           "Missing day",
			payload:        `{"date": "/1/2023"}`,
			expectedStatus: http.StatusBadRequest,
			expectedValues: models.OutputResults{
				AcscEuropean:          "",
				AcscHundredYear:       "",
				AcscInternational:     "",
				AcscJulian:            "",
				AcscUsaStandard:       "",
				DayOfWeek:             "",
				ErrorFlag:             "HTTP 400",
				ErrorText:             "invalid date: /1/2023",
				EuropeanStandard:      "",
				InternationalStandard: "",
				UsaStandard:           "",
			},
		},
		{
			name:           "Missing month",
			payload:        `{"date": "1//2023"}`,
			expectedStatus: http.StatusBadRequest,
			expectedValues: models.OutputResults{
				AcscEuropean:          "",
				AcscHundredYear:       "",
				AcscInternational:     "",
				AcscJulian:            "",
				AcscUsaStandard:       "",
				DayOfWeek:             "",
				ErrorFlag:             "HTTP 400",
				ErrorText:             "invalid date: 1//2023",
				EuropeanStandard:      "",
				InternationalStandard: "",
				UsaStandard:           "",
			},
		},
		{
			name:           "Missing year",
			payload:        `{"date": "1/1/"}`,
			expectedStatus: http.StatusBadRequest,
			expectedValues: models.OutputResults{
				AcscEuropean:          "",
				AcscHundredYear:       "",
				AcscInternational:     "",
				AcscJulian:            "",
				AcscUsaStandard:       "",
				DayOfWeek:             "",
				ErrorFlag:             "HTTP 400",
				ErrorText:             "invalid date: 1/1/",
				EuropeanStandard:      "",
				InternationalStandard: "",
				UsaStandard:           "",
			},
		},
		{
			name:           "Dash separators",
			payload:        `{"date": "1-1-2023"}`,
			expectedStatus: http.StatusBadRequest,
			expectedValues: models.OutputResults{
				AcscEuropean:          "",
				AcscHundredYear:       "",
				AcscInternational:     "",
				AcscJulian:            "",
				AcscUsaStandard:       "",
				DayOfWeek:             "",
				ErrorFlag:             "HTTP 400",
				ErrorText:             "invalid separators: use / instead",
				EuropeanStandard:      "",
				InternationalStandard: "",
				UsaStandard:           "",
			},
		},
		{
			name:           "Dot separators",
			payload:        `{"date": "1.1.2023"}`,
			expectedStatus: http.StatusBadRequest,
			expectedValues: models.OutputResults{
				AcscEuropean:          "",
				AcscHundredYear:       "",
				AcscInternational:     "",
				AcscJulian:            "",
				AcscUsaStandard:       "",
				DayOfWeek:             "",
				ErrorFlag:             "HTTP 400",
				ErrorText:             "invalid separators: use / instead",
				EuropeanStandard:      "",
				InternationalStandard: "",
				UsaStandard:           "",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/api/CalcCalendarDate", strings.NewReader(tc.payload))
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
