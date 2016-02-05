package weaselbot

import (
	"bytes"
	"text/template"
)

const notificationTmplS = `Hey {{.User_Name }}, you used some language in #{{ .Channel }} matching these weasel expressions:
{{ range .Words }}* {{ . }}
{{ end }}`

var (
	notificationTmpl *template.Template
)

type Notification struct {
	User_Name string
	Channel   string
	Words     Words
}

func (n Notification) String() string {
	buf := bytes.NewBuffer(nil)

	if err := notificationTmpl.Execute(buf, n); err != nil {
		panic(err)
	}

	return string(buf.Bytes())
}

func init() {
	notificationTmpl = template.Must(template.New("MachineList").Parse(notificationTmplS))
}
