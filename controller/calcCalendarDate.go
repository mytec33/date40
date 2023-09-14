package controller

import (
	"date_calculation/models"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Green screen places a dot at the end
var weekdayAbbreviations = map[time.Weekday]string{
	time.Sunday:    "SUN.",
	time.Monday:    "MON.",
	time.Tuesday:   "TUE.",
	time.Wednesday: "WED.",
	time.Thursday:  "THU.",
	time.Friday:    "FRI.",
	time.Saturday:  "SAT.",
}

func CalcCalendarDate(context *gin.Context) {
	var input models.InputCalendarDate
	var output models.OutputResults

	handleError := func(status int, errorString string) {
		output.ErrorFlag = "HTTP: " + strconv.Itoa(status)
		output.ErrorText = errorString
		context.JSON(status, gin.H{"results": output})
	}

	if err := context.ShouldBindJSON(&input); err != nil {
		handleError(http.StatusBadRequest, err.Error())
		return
	}

	if input.Date == "" {
		output.ErrorFlag = "HTTP: " + strconv.Itoa(http.StatusBadRequest)
		output.ErrorText = "invalid date: empty"
		context.JSON(http.StatusBadRequest, gin.H{"results": output})

		return
	}

	parsedDate, err := time.Parse("1/2/2006", input.Date)
	if err != nil {
		output.ErrorFlag = "HTTP: " + strconv.Itoa(http.StatusBadRequest)
		output.ErrorText = "invalid date: " + input.Date
		context.JSON(http.StatusBadRequest, gin.H{"results": output})

		return
	}

	output = calcDatesByCalendarDate(parsedDate.String())

	context.JSON(http.StatusOK, gin.H{"results": output})
}

func isHydInRange(number int) bool {
	return number >= 0 && number <= 99999
}

func CalcHundreYearDate(context *gin.Context) {
	var input models.InputHundredYearDate
	var output models.OutputResults

	handleError := func(status int, errorString string) {
		output.ErrorFlag = "HTTP: " + strconv.Itoa(status)
		output.ErrorText = errorString
		context.JSON(status, gin.H{"results": output})
	}

	if err := context.ShouldBindJSON(&input); err != nil {
		handleError(http.StatusBadRequest, err.Error())
		return
	}

	if input.HundredYear == "" {
		handleError(http.StatusBadRequest, "invalid 100 year date: empty")
		return
	}

	hundredYear, err := strconv.Atoi(input.HundredYear)
	if err != nil {
		handleError(http.StatusBadRequest, "invalid 100 year date: must be a positive number")
		return
	}

	if !isHydInRange(hundredYear) {
		handleError(http.StatusBadRequest, "100 year date out of range: must be between 0 and 99999")
		return
	}

	inputDate, err := calcCalendarDateByHundredYear(hundredYear)
	if err != nil {
		handleError(http.StatusBadRequest, "error converting dates")
		return
	}

	output = calcDatesByCalendarDate(inputDate)
	context.JSON(http.StatusOK, gin.H{"results": output})
}

func calcCalendarDateByHundredYear(inputDate int) (string, error) {
	referenceDate := time.Date(1899, 12, 31, 0, 0, 0, 0, time.UTC)

	daysSinceReference := inputDate
	calculatedDate := referenceDate.Add(time.Duration(daysSinceReference) * 24 * time.Hour)

	return calculatedDate.Format("1/02/2006"), nil
}

func calcDatesByCalendarDate(inputDate string) models.OutputResults {
	var output models.OutputResults

	output.AcscEuropean = calcAcscEuropean(inputDate)
	output.AcscHundredYear = calcAcscHundredYear(inputDate)
	output.AcscInternational = calcAcscInternationalStandard(inputDate)
	output.AcscJulian = calcAcscJulian(inputDate)
	output.AcscUsaStandard = calcAcscUsaStandard(inputDate)
	output.DayOfWeek = calcDayOfWeek(inputDate)
	output.EuropeanStandard = calcEuropeanStandard(inputDate)
	output.InternationalStandard = calcInternationalStandard(inputDate)
	output.UsaStandard = padUsaStandard(inputDate)
	output.ErrorFlag = "0"

	return output
}

func calcAcscEuropean(inputDate string) string {
	return formatUsaStandard(inputDate, "02.01.06")
}

func calcAcscHundredYear(inputDate string) string {
	parsedDate, err := time.Parse("1/2/2006", inputDate)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	referenceDate := time.Date(1899, 12, 31, 0, 0, 0, 0, time.UTC)
	days := parsedDate.Sub(referenceDate).Hours() / 24

	return strconv.Itoa(int(days))
}

func calcEuropeanStandard(inputDate string) string {
	return formatUsaStandard(inputDate, "02.01.2006")
}

func calcAcscInternationalStandard(inputDate string) string {
	return formatUsaStandard(inputDate, "06-01-02")
}

func calcAcscJulian(inputDate string) string {
	parsedDate, err := time.Parse("1/2/2006", inputDate)
	if err != nil {
		fmt.Println("Error:", err)
		return "unable to parse Julian"
	}

	as400JulianYear := parsedDate.Year() % 100
	startOfYear := time.Date(parsedDate.Year(), 1, 1, 0, 0, 0, 0, time.UTC)
	daysInto2023 := int(parsedDate.Sub(startOfYear).Hours()/24) + 1

	return fmt.Sprintf("%02d-%03d", as400JulianYear, daysInto2023)
}

func calcAcscUsaStandard(inputDate string) string {
	date := formatUsaStandard(inputDate, "1/2/06")
	if len(date) == 6 {
		date = "  " + date
	} else if len(date) == 7 {
		date = " " + date
	}

	return date
}

func calcDayOfWeek(inputDate string) string {
	parsedDate, err := time.Parse("1/2/2006", inputDate)
	if err != nil {
		fmt.Println("Error parsing date:", err)
		return ""
	}

	dayOfWeek := parsedDate.Weekday()
	return weekdayAbbreviations[dayOfWeek]
}

func calcInternationalStandard(inputDate string) string {
	return formatUsaStandard(inputDate, "2006-01-02")
}

func formatUsaStandard(inputDate string, format string) string {
	parsedDate, err := time.Parse("1/2/2006", inputDate)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	return parsedDate.Format(format)
}

func padUsaStandard(inputDate string) string {
	date := strings.TrimSpace(inputDate)

	if len(date) == 8 {
		date = "  " + date
	} else if len(date) == 9 {
		date = " " + date
	}

	return date
}
