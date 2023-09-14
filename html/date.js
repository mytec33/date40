document.addEventListener('DOMContentLoaded', function () {
  const calculateButton = document.getElementById('calculateButton');
  const dateInput = document.getElementById('dateInput');
  const prompt = document.querySelector('.prompt');

  function formatNumberWithLeadingZeros(number, length) {
    return String(number).padStart(length, '0');
  }

  function formatDateInCustomFormat(date) {
    const day = formatNumberWithLeadingZeros(date.getDate(), 2);
    const month = formatNumberWithLeadingZeros(date.getMonth() + 1, 2); // Months are zero-based
    const year = date.getFullYear();

    return `${month}${day}${year}`;
  }

  function formatDateInCustomFormatUTC(date) {
    const utcDay = formatNumberWithLeadingZeros(date.getUTCDate(), 2);
    const utcMonth = formatNumberWithLeadingZeros(date.getUTCMonth() + 1, 2); // Months are zero-based
    const utcYear = date.getUTCFullYear();

    return `${utcMonth}${utcDay}${utcYear}`;
  }

  function formatTimeMiltary(date) {
    const hours = formatNumberWithLeadingZeros(date.getHours(), 2);
    const minutes = formatNumberWithLeadingZeros(date.getMinutes(), 2);
    const seconds = formatNumberWithLeadingZeros(date.getSeconds(), 2);

    return `${hours}:${minutes}:${seconds}`;
  }

  function formatInput() {
    if (dateInput.value === 7) {
      dateInput.value = dateInput.value.replace(/(\d{1,2})(\d{2})(\d{4})/, '$1/$2/$3')
    }

    const inputValue = dateInput.value;
    const formattedValue = inputValue.replace(/(\d{1,2})(\d{2})(\d{4})/, '$1/$2/$3');
    dateInput.value = formattedValue;
  }

  function isValidDate(inputDate) {
    const dateFormats = [
      /^\d{4}-\d{2}-\d{2}$/,    // YYYY-MM-DD
      /^\d{2}\/\d{2}\/\d{4}$/,  // MM/DD/YYYY
      /^\d{2}-\d{2}-\d{4}$/,    // MM-DD-YYYY
    ];

    const matchedFormat = dateFormats.find(format => format.test(inputDate));
    if (!matchedFormat) {
      return false;
    }

    const parsedDate = new Date(inputDate);

    return !isNaN(parsedDate) && parsedDate instanceof Date;
  }


  dateInput.addEventListener('blur', () => {
    if (!dateInput.value) {
      prompt.style.display = 'inline';
    }
  });

  dateInput.addEventListener('input', () => {
    formatInput();
  });

  calculateButton.addEventListener('click', function () {
    makeAPICall('https://dolotsoflittlethings.com:8010/api/CalcCalendarDate', dateInput.value);
  });

  calculateHundredYear.addEventListener('click', function () {
    makeAPICall('https://dolotsoflittlethings.com:8010/api/CalcHundredYearDate', hundredYearInput.value);
  });

  dateInput.addEventListener('keyup', function (event) {
    if (event.key === 'Enter') {
      makeAPICall('https://dolotsoflittlethings.com:8010/api/CalcCalendarDate', dateInput.value);
    }
  });

  hundredYearInput.addEventListener('keyup', function (event) {
    if (event.key === 'Enter') {
      makeAPICall('https://dolotsoflittlethings.com:8010/api/CalcHundredYearDate', hundredYearInput.value);
    }
  });

  function makeAPICall(url, input) {
    fetch(url, {
      method: 'POST',
      body: JSON.stringify({ date: input }),
      headers: {
        'Content-Type': 'application/json'
      }
    })
      .then(response => response.json())
      .then(data => {
        const results = data.results;

        document.getElementById('acscUsaStandard').innerText = results.AcscUsaStandard;
        document.getElementById('acscInternational').innerText = results.AcscInternational;
        document.getElementById('acscEuropean').innerText = results.AcscEuropean;
        document.getElementById('acscHundredYear').innerText = results.AcscHundredYear;
        document.getElementById('acscJulian').innerText = results.AcscJulian;
        document.getElementById('usaStandard').innerText = results.UsaStandard;
        document.getElementById('internationalStandard').innerText = results.InternationalStandard;
        document.getElementById('europeanStandard').innerText = results.EuropeanStandard;
        document.getElementById('dayOfWeek').innerText = results.DayOfWeek;
        document.getElementById('errorFlag').innerText = results.ErrorFlag;
        document.getElementById('errorText').innerText = results.ErrorText;
      })
      .catch(error => {
        console.error(error);
      });
  }

  window.onload = function () {
    const currentDate = new Date();
    const formattedDate = formatDateInCustomFormat(currentDate);
    const formattedDateUTC = formatDateInCustomFormatUTC(currentDate);
    const formattedPageTime = formatTimeMiltary(currentDate)

    const formattedDateElement = document.getElementById('sysdate');
    const formattedDateElementUTC = document.getElementById('udate');
    const formattedTimeElement = document.getElementById("pageTime")

    formattedDateElement.textContent = formattedDate;
    formattedDateElementUTC.textContent = formattedDateUTC;
    formattedTimeElement.textContent = formattedPageTime;
  }
});