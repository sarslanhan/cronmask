package conditions

import (
	"testing"
	"time"
)

func TestParseCronMaskExpression(t *testing.T) {
	testCases := []struct {
		expr    string
		ts      string
		matches bool
		isError bool
	}{
		// test minute field with different combinations
		{expr: "* * * * *", ts: "2019-07-04T07:00:00Z", matches: true},
		{expr: "2 * * * *", ts: "2019-07-04T07:00:00Z", matches: false},
		{expr: "2 * * * *", ts: "2019-07-04T07:02:00Z", matches: true},

		{expr: "20 * * * *", ts: "2019-07-04T07:02:00Z", matches: false},
		{expr: "20 * * * *", ts: "2019-07-04T07:20:00Z", matches: true},

		{expr: "20-40 * * * *", ts: "2019-07-04T07:00:00Z", matches: false},
		{expr: "20-40 * * * *", ts: "2019-07-04T07:20:00Z", matches: true},
		{expr: "20-40 * * * *", ts: "2019-07-04T07:21:00Z", matches: true},
		{expr: "20-40 * * * *", ts: "2019-07-04T07:40:00Z", matches: true},
		{expr: "20-40 * * * *", ts: "2019-07-04T07:41:00Z", matches: false},

		{expr: "20-40,50 * * * *", ts: "2019-07-04T07:50:00Z", matches: true},

		{expr: "20,40,50 * * * *", ts: "2019-07-04T07:20:00Z", matches: true},
		{expr: "20,40,50 * * * *", ts: "2019-07-04T07:20:30Z", matches: true},
		{expr: "20,40,50 * * * *", ts: "2019-07-04T07:21:00Z", matches: false},
		{expr: "20,40,50 * * * *", ts: "2019-07-04T07:40:00Z", matches: true},
		{expr: "20,40,50 * * * *", ts: "2019-07-04T07:50:00Z", matches: true},
		{expr: "20,40,50 * * * *", ts: "2019-07-04T07:51:00Z", matches: false},

		{expr: "10-20,30,40-59 * * * *", ts: "2019-07-04T07:10:00Z", matches: true},
		{expr: "10-20,30,40-59 * * * *", ts: "2019-07-04T07:10:30Z", matches: true},
		{expr: "10-20,30,40-59 * * * *", ts: "2019-07-04T07:11:30Z", matches: true},
		{expr: "10-20,30,40-59 * * * *", ts: "2019-07-04T07:20:30Z", matches: true},
		{expr: "10-20,30,40-59 * * * *", ts: "2019-07-04T07:21:30Z", matches: false},
		{expr: "10-20,30,40-59 * * * *", ts: "2019-07-04T07:30:30Z", matches: true},
		{expr: "10-20,30,40-59 * * * *", ts: "2019-07-04T07:31:40Z", matches: false},
		{expr: "10-20,30,40-59 * * * *", ts: "2019-07-04T07:41:40Z", matches: true},
		{expr: "10-20,30,40-59 * * * *", ts: "2019-07-04T07:59:40Z", matches: true},
		{expr: "10-20,30,40-59 * * * *", ts: "2019-07-04T08:00:00Z", matches: false},

		// test hour field with different combinations
		{expr: "* * * * *", ts: "2019-07-04T07:00:00Z", matches: true},
		{expr: "* 2 * * *", ts: "2019-07-04T07:00:00Z", matches: false},
		{expr: "* 2 * * *", ts: "2019-07-04T02:00:00Z", matches: true},

		{expr: "* 20 * * *", ts: "2019-07-04T10:00:00Z", matches: false},
		{expr: "* 20 * * *", ts: "2019-07-04T20:00:00Z", matches: true},

		{expr: "* 10-20 * * *", ts: "2019-07-04T07:00:00Z", matches: false},
		{expr: "* 10-20 * * *", ts: "2019-07-04T10:00:00Z", matches: true},
		{expr: "* 10-20 * * *", ts: "2019-07-04T11:00:00Z", matches: true},
		{expr: "* 10-20 * * *", ts: "2019-07-04T20:00:00Z", matches: true},
		{expr: "* 10-20 * * *", ts: "2019-07-04T21:00:00Z", matches: false},

		{expr: "* 10-20,22 * * *", ts: "2019-07-04T22:00:00Z", matches: true},

		{expr: "* 10,15,20 * * *", ts: "2019-07-04T10:00:00Z", matches: true},
		{expr: "* 10,15,20 * * *", ts: "2019-07-04T15:00:30Z", matches: true},
		{expr: "* 10,15,20 * * *", ts: "2019-07-04T17:00:00Z", matches: false},
		{expr: "* 10,15,20 * * *", ts: "2019-07-04T20:00:00Z", matches: true},
		{expr: "* 10,15,20 * * *", ts: "2019-07-04T23:00:00Z", matches: false},

		{expr: "* 10-15,17,20-23 * * *", ts: "2019-07-04T10:00:00Z", matches: true},
		{expr: "* 10-15,17,20-23 * * *", ts: "2019-07-04T11:00:30Z", matches: true},
		{expr: "* 10-15,17,20-23 * * *", ts: "2019-07-04T15:00:30Z", matches: true},
		{expr: "* 10-15,17,20-23 * * *", ts: "2019-07-04T16:00:30Z", matches: false},
		{expr: "* 10-15,17,20-23 * * *", ts: "2019-07-04T17:00:30Z", matches: true},
		{expr: "* 10-15,17,20-23 * * *", ts: "2019-07-04T20:00:30Z", matches: true},
		{expr: "* 10-15,17,20-23 * * *", ts: "2019-07-04T19:00:40Z", matches: false},
		{expr: "* 10-15,17,20-23 * * *", ts: "2019-07-04T21:00:40Z", matches: true},
		{expr: "* 10-15,17,20-23 * * *", ts: "2019-07-04T23:00:40Z", matches: true},

		// test day of month field with different combinations
		{expr: "* * * * *", ts: "2019-07-04T07:00:00Z", matches: true},
		{expr: "* * 5 * *", ts: "2019-07-04T07:00:00Z", matches: false},
		{expr: "* * 4 * *", ts: "2019-07-04T07:00:00Z", matches: true},

		{expr: "* * 20 * *", ts: "2019-07-10T07:00:00Z", matches: false},
		{expr: "* * 20 * *", ts: "2019-07-20T07:00:00Z", matches: true},

		{expr: "* * 10-20 * *", ts: "2019-07-04T07:00:00Z", matches: false},
		{expr: "* * 10-20 * *", ts: "2019-07-10T07:00:00Z", matches: true},
		{expr: "* * 10-20 * *", ts: "2019-07-15T07:00:00Z", matches: true},
		{expr: "* * 10-20 * *", ts: "2019-07-20T07:00:00Z", matches: true},
		{expr: "* * 10-20 * *", ts: "2019-07-21T07:00:00Z", matches: false},

		{expr: "* * 10-20,22 * *", ts: "2019-07-04T07:00:00Z", matches: false},

		{expr: "* * 10,15,20 * *", ts: "2019-07-10T07:00:00Z", matches: true},
		{expr: "* * 10,15,20 * *", ts: "2019-07-15T07:00:30Z", matches: true},
		{expr: "* * 10,15,20 * *", ts: "2019-07-17T07:00:00Z", matches: false},
		{expr: "* * 10,15,20 * *", ts: "2019-07-20T07:00:00Z", matches: true},
		{expr: "* * 10,15,20 * *", ts: "2019-07-22T07:00:00Z", matches: false},

		{expr: "* * 10-15,17,20-22 * *", ts: "2019-07-04T07:00:00Z", matches: false},
		{expr: "* * 10-15,17,20-22 * *", ts: "2019-07-10T07:00:30Z", matches: true},
		{expr: "* * 10-15,17,20-22 * *", ts: "2019-07-15T07:00:30Z", matches: true},
		{expr: "* * 10-15,17,20-22 * *", ts: "2019-07-16T07:00:30Z", matches: false},
		{expr: "* * 10-15,17,20-22 * *", ts: "2019-07-17T07:00:30Z", matches: true},
		{expr: "* * 10-15,17,20-22 * *", ts: "2019-07-18T07:00:30Z", matches: false},
		{expr: "* * 10-15,17,20-22 * *", ts: "2019-07-20T07:00:40Z", matches: true},
		{expr: "* * 10-15,17,20-22 * *", ts: "2019-07-22T07:00:40Z", matches: true},
		{expr: "* * 10-15,17,20-22 * *", ts: "2019-07-23T07:00:40Z", matches: false},

		// test month field with different combinations
		{expr: "* * * * *", ts: "2019-07-04T07:00:00Z", matches: true},
		{expr: "* * * 5 *", ts: "2019-07-04T07:00:00Z", matches: false},
		{expr: "* * * 7 *", ts: "2019-07-04T07:00:00Z", matches: true},

		{expr: "* * * 10 *", ts: "2019-11-20T07:00:00Z", matches: false},
		{expr: "* * * 11 *", ts: "2019-11-20T07:00:00Z", matches: true},

		{expr: "* * * 5-10 *", ts: "2019-03-04T07:00:00Z", matches: false},
		{expr: "* * * 5-10 *", ts: "2019-05-04T07:00:00Z", matches: true},
		{expr: "* * * 5-10 *", ts: "2019-07-04T07:00:00Z", matches: true},
		{expr: "* * * 5-10 *", ts: "2019-10-04T07:00:00Z", matches: true},
		{expr: "* * * 5-10 *", ts: "2019-11-04T07:00:00Z", matches: false},

		{expr: "* * * 3-7,10 *", ts: "2019-01-04T07:00:00Z", matches: false},

		{expr: "* * * 3,5,10 *", ts: "2019-01-04T07:00:00Z", matches: false},
		{expr: "* * * 3,5,10 *", ts: "2019-03-04T07:00:00Z", matches: true},
		{expr: "* * * 3,5,10 *", ts: "2019-04-04T07:00:30Z", matches: false},
		{expr: "* * * 3,5,10 *", ts: "2019-05-04T07:00:00Z", matches: true},
		{expr: "* * * 3,5,10 *", ts: "2019-10-04T07:00:00Z", matches: true},
		{expr: "* * * 3,5,10 *", ts: "2019-11-04T07:00:00Z", matches: false},

		{expr: "* * * 3-5,7,9-11 *", ts: "2019-01-04T07:00:00Z", matches: false},
		{expr: "* * * 3-5,7,9-11 *", ts: "2019-03-04T07:00:30Z", matches: true},
		{expr: "* * * 3-5,7,9-11 *", ts: "2019-04-04T07:00:30Z", matches: true},
		{expr: "* * * 3-5,7,9-11 *", ts: "2019-05-04T07:00:30Z", matches: true},
		{expr: "* * * 3-5,7,9-11 *", ts: "2019-06-04T07:00:30Z", matches: false},
		{expr: "* * * 3-5,7,9-11 *", ts: "2019-07-04T07:00:30Z", matches: true},
		{expr: "* * * 3-5,7,9-11 *", ts: "2019-08-04T07:00:30Z", matches: false},
		{expr: "* * * 3-5,7,9-11 *", ts: "2019-09-04T07:00:40Z", matches: true},
		{expr: "* * * 3-5,7,9-11 *", ts: "2019-11-04T07:00:40Z", matches: true},
		{expr: "* * * 3-5,7,9-11 *", ts: "2019-12-04T07:00:40Z", matches: false},

		// test day of week field with different combinations
		{expr: "* * * * *", ts: "2019-07-04T07:00:00Z", matches: true},
		{expr: "* * * * 0", ts: "2019-07-07T07:00:00Z", matches: true},
		{expr: "* * * * 1", ts: "2019-07-07T07:00:00Z", matches: false},

		{expr: "* * * * 2-4", ts: "2019-07-08T07:00:00Z", matches: false},
		{expr: "* * * * 2-4", ts: "2019-07-09T07:00:00Z", matches: true},
		{expr: "* * * * 2-4", ts: "2019-07-10T07:00:00Z", matches: true},
		{expr: "* * * * 2-4", ts: "2019-07-11T07:00:00Z", matches: true},
		{expr: "* * * * 2-4", ts: "2019-07-12T07:00:00Z", matches: false},

		{expr: "* * * * 2,4,5", ts: "2019-07-08T07:00:00Z", matches: false},
		{expr: "* * * * 2,4,5", ts: "2019-07-09T07:00:00Z", matches: true},
		{expr: "* * * * 2,4,5", ts: "2019-07-10T07:00:30Z", matches: false},
		{expr: "* * * * 2,4,5", ts: "2019-07-11T07:00:00Z", matches: true},
		{expr: "* * * * 2,4,5", ts: "2019-07-12T07:00:00Z", matches: true},
		{expr: "* * * * 2,4,5", ts: "2019-07-13T07:00:00Z", matches: false},

		{expr: "* * * * 0-1,3,4-5", ts: "2019-07-07T07:00:30Z", matches: true},
		{expr: "* * * * 0-1,3,4-5", ts: "2019-07-08T07:00:30Z", matches: true},
		{expr: "* * * * 0-1,3,4-5", ts: "2019-07-09T07:00:30Z", matches: false},
		{expr: "* * * * 0-1,3,4-5", ts: "2019-07-10T07:00:30Z", matches: true},
		{expr: "* * * * 0-1,3,4-5", ts: "2019-07-11T07:00:30Z", matches: true},
		{expr: "* * * * 0-1,3,4-5", ts: "2019-07-12T07:00:40Z", matches: true},
		{expr: "* * * * 0-1,3,4-5", ts: "2019-07-13T07:00:40Z", matches: false},

		// invalid cases
		{expr: "* * * * * *", ts: "2019-07-07T07:00:00Z", isError: true},
		{expr: "* 1,,3 * * *", ts: "2019-07-07T07:00:00Z", isError: true},
		{expr: "* a * * *", ts: "2019-07-07T07:00:00Z", isError: true},
		{expr: "* a-3 * * *", ts: "2019-07-07T07:00:00Z", isError: true},
		{expr: "* 2-a * * *", ts: "2019-07-07T07:00:00Z", isError: true},
		{expr: "* 2-3-5 * * *", ts: "2019-07-07T07:00:00Z", isError: true},
	}

	for _, tc := range testCases {
		ts, err := time.Parse(time.RFC3339Nano, tc.ts)
		if err != nil {
			t.Error(err)
			continue
		}

		cron, err := New(tc.expr)
		if err != nil {
			if !tc.isError {
				t.Error(err)
			}
			continue
		}

		if result := cron.Matches(ts); result != tc.matches {
			t.Errorf("Expected result [%v] but got [%v] for expression %s and timestamp %s", tc.matches, result, tc.expr, tc.ts)
		}
	}
}
