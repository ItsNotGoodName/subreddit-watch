package templater

import (
	"bytes"
	"log"
	"text/template"

	"github.com/turnage/graw/reddit"
)

const DefaultTitleTemplate string = "{{ .Post.Title }}"
const DefaultMessageTemplate string = `{{ if .Post.IsSelf }}{{ .Post.URL }}{{ else }}https://reddit.com{{ .Post.Permalink }}
{{ .Post.URL }}{{ end }}`

type Data struct {
	Post *reddit.Post
}

type Templater struct {
	title   *template.Template
	message *template.Template
}

func New(titleTemplate, messageTemplate *template.Template) *Templater {
	return &Templater{
		title:   titleTemplate,
		message: messageTemplate,
	}
}

func run(t *template.Template, p *reddit.Post) string {
	buffer := bytes.NewBuffer([]byte{})
	if err := t.Execute(buffer, Data{Post: p}); err != nil {
		log.Println("templater.run:", err)
		return ""
	}

	return buffer.String()
}

func (t *Templater) GetTitle(p *reddit.Post) string {
	return run(t.title, p)
}

func (t *Templater) GetMessage(p *reddit.Post) string {
	return run(t.message, p)
}
