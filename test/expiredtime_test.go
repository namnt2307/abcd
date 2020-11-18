package expiredtime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// func main() {

// 	abc := GetExpiredDateFromDuration(6, "months", "")

// 	fmt.Println("out", abc)

// }

func TestAdd(t *testing.T) {
	assert := assert.New(t)

	var tests = []struct {
		duration      int
		duration_type string
		currentDate   string
		expected      string
	}{
		{1, "months", "2020-01-31 11:30:00", "2020-02-29 23:59:59"},
		{1, "months", "2020-02-01 11:30:00", "2020-02-29 23:59:59"},
		{1, "months", "2021-02-01 11:30:00", "2021-02-28 23:59:59"},
		{1, "months", "2020-02-10 11:30:00", "2020-03-09 23:59:59"},
		{1, "months", "2020-03-01 11:30:00", "2020-03-31 23:59:59"},
		{1, "months", "2020-04-01 11:30:00", "2020-04-30 23:59:59"},
		{1, "months", "2020-02-29 11:30:00", "2020-03-28 23:59:59"},
		{1, "months", "2019-12-29 11:30:00", "2020-01-28 23:59:59"},
		{1, "months", "2019-01-31 11:30:00", "2019-02-28 23:59:59"},
		{3, "months", "2020-11-30 11:30:00", "2021-02-28 23:59:59"},
		{3, "months", "2020-12-01 11:30:00", "2021-02-28 23:59:59"},
		{3, "months", "2020-12-30 11:30:00", "2021-03-29 23:59:59"},
		{3, "months", "2020-10-02 11:30:00", "2021-01-01 23:59:59"},
		{3, "months", "2020-10-01 11:30:00", "2020-12-31 23:59:59"},
		{12, "months", "2020-05-01 11:30:00", "2021-04-30 23:59:59"},
		{12, "months", "2020-04-30 11:30:00", "2021-04-29 23:59:59"},
		{12, "months", "2020-02-29 11:30:00", "2021-02-28 23:59:59"},
		{12, "months", "2020-03-01 11:30:00", "2021-02-28 23:59:59"},
		{12, "months", "2020-01-02 11:30:00", "2021-01-01 23:59:59"},
		{12, "months", "2020-01-01 11:30:00", "2020-12-31 23:59:59"},
		{1, "months_end_day", "2020-02-29 11:30:00", "2020-03-31 23:59:59"},
		{1, "months_end_day", "2020-03-31 11:30:00", "2020-04-30 23:59:59"},
		{3, "months_end_day", "2020-03-31 11:30:00", "2020-06-30 23:59:59"},
		{12, "months_end_day", "2020-01-30 11:30:00", "2021-01-30 23:59:59"},
		{12, "months_end_day", "2020-01-31 11:30:00", "2021-01-31 23:59:59"},
		{11, "months_end_day", "2020-01-31 11:30:00", "2020-12-31 23:59:59"},
		{1, "months_end_day", "2020-11-30 11:30:00", "2020-12-31 23:59:59"},
		{3, "months_end_day", "2020-09-30 11:30:00", "2020-12-31 23:59:59"},
	}

	for _, test := range tests {
		assert.Equal(test.expected, GetExpiredDateFromDuration(test.duration, test.duration_type, test.currentDate))
	}
}

func GetExpiredDateFromDuration(duration int, duration_type, currentDate string) string {
	layout := "2006-01-02 15:04:05"

	//default in hours
	if duration_type == "" {
		duration_type = "hours"
	}

	//get current time
	// currentTime := time.Now()
	currentTime, _ := time.Parse(layout, currentDate)
	currentTime = currentTime.Local()

	//Ngày dời lại mặc định là -1
	var dayToBack = -1
	if duration_type == "months_end_day" {
		dayToBack = 0
	}

	var expiredTime time.Time

	switch duration_type {
	case "hours":
		expiredTime = currentTime.Add(time.Duration(duration) * time.Hour)
	case "days":
		expiredTime = currentTime.AddDate(0, 0, duration)
	case "months", "months_end_day":

		_, currentMonth, _ := currentTime.Date()

		//Tính tháng hết hạn
		var expireMonth = (int(currentMonth) + duration)
		if expireMonth > 12 {
			expireMonth = expireMonth % 12
		}

		//trường hợp bình thường + theo công thức d-1/m+1
		expiredTime = currentTime.AddDate(0, duration, dayToBack)
		expiredYear, expiredMonth, _ := expiredTime.Date()
		expiredHours, expiredMinutes, expiredSeconds := expiredTime.Clock()

		//trường hợp tháng hiện tại + duration rơi vào tháng 2 và tháng hết hạn sau khi tính toán lại nhảy tới tháng 3 thì ngày hết hạn cần back lại về cuối tháng 2
		//trường hợp ngày mua là ngày cuối tháng thì ngày hết hạn và gia hạn cũng là ngày cuối tháng (chỉ áp dụng đối với months_end_day)
		if (expireMonth == 2 && expireMonth < int(expiredMonth)) || (duration_type == "months_end_day" && IsEndOfMonth(currentTime)) {
			expiredTime = time.Date(expiredYear, time.Month(expireMonth)+1, 0, expiredHours, expiredMinutes, expiredSeconds, 0, currentTime.Location())
		}

	case "years":
		expiredTime = currentTime.AddDate(duration, 0, dayToBack)
	}

	year, month, day := expiredTime.Date()
	expiredTime = time.Date(year, month, day, 23, 59, 59, 0, time.Local)

	outPutFormat := "2006-01-02 15:04:05"
	return expiredTime.Format(outPutFormat)
}

func IsEndOfMonth(t time.Time) bool {
	_, m, _ := t.Date()
	_, month, _ := t.AddDate(0, 0, 1).Date()
	return int(month) != int(m)
}

// func main() {
// 	GetExpiredDateFromDuration(1, "months", "2019-01-31 11:30:00")
// }

// func GetExpiredDateFromDuration(duration int, duration_type string, currentDate string) string {
// 	// if outPutFormat == "" {
// 	outPutFormat := "2006-01-02 15:04:05"
// 	// 15:04:05
// 	// }

// 	//default in hours
// 	if duration_type == "" {
// 		duration_type = "hours"
// 	}

// 	//get current time
// 	// currentTime := time.Now()
// 	currentTime, _ := time.Parse(outPutFormat, currentDate)

// 	// fmt.Println("currentTime", currentTime)

// 	var expiredTime time.Time

// 	switch duration_type {
// 	case "hours":
// 		expiredTime = currentTime.Add(time.Duration(duration) * time.Hour)
// 	case "days":
// 		expiredTime = currentTime.AddDate(0, 0, duration)
// 	case "months":

// 		_, curr_month, _ := currentTime.Date()
// 		expiredTime = currentTime.AddDate(0, duration, -1)
// 		year, month, _ := expiredTime.Date()
// 		h, i, s := expiredTime.Clock()

// 		//trường hợp tháng hiện tại + duration rơi vào tháng 2 và tháng sau khi tính được tăng lên tháng 3
// 		if (int(curr_month)+duration)%12 == 2 && int(month) == 3 {
// 			expiredTime = time.Date(year, month, 0, h, i, s, 0, currentTime.Location())
// 		}

// 		// _, curr_month, curr_day := currentTime.Date()

// 		// if curr_day == 1 && duration == 1 {
// 		// 	expiredTime = currentTime.AddDate(0, 1, -1)
// 		// } else {
// 		// 	expiredTime = currentTime.AddDate(0, duration, -1)
// 		// 	year, month, _ := expiredTime.Date()
// 		// 	h, i, s := expiredTime.Clock()

// 		// 	//trường hợp tháng hiện tại + duration rơi vào tháng 2 và
// 		// 	if (int(curr_month)+duration)%12 == 2 && int(month) == 3 {
// 		// 		expiredTime = time.Date(year, month, 0, h, i, s, 0, currentTime.Location())
// 		// 	}
// 		// }

// 	case "years":
// 		expiredTime = currentTime.AddDate(duration, 0, -1)
// 	}

// 	return expiredTime.Format(outPutFormat)
// }
