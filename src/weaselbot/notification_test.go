package weaselbot

import (
	"testing"
)

func TestNotificationString(t *testing.T) {
	cases := []struct {
		n      Notification
		expect string
	}{
		{
			Notification{"friend", "general", Words{"strong", "good"}},
			`Hey friend, you used some language in #general matching these weasel expressions:
* strong
* good
`,
		},
	}

	for _, c := range cases {
		if s := c.n.String(); s != c.expect {
			t.Errorf("Bad string for notification. Expected:\n%s\n\nGot:\n%s\n", c.expect, s)
			t.Fail()
		}
	}
}
