package backend

import (
	"bytes"
	"encoding/json"
	"net/textproto"
	"text/template"

	"github.com/kofuk/spaghettini/server/types"
)

type TemplateFuncsCollections struct{}

func (f TemplateFuncsCollections) Map(values ...any) map[string]any {
	if len(values)%2 != 0 {
		return nil
	}
	m := make(map[string]any)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil
		}
		m[key] = values[i+1]
	}
	return m
}

type TemplateFuncsEncoding struct{}

func (f TemplateFuncsEncoding) JSON(v any) string {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		return ""
	}
	return buf.String()
}

type TemplateFuncsRequests struct{}

func (f TemplateFuncsRequests) GetHeader(request *types.Request, key string) string {
	return textproto.MIMEHeader(request.Header).Get(key)
}

var funcMap = template.FuncMap{
	"collections": func() TemplateFuncsCollections {
		return TemplateFuncsCollections{}
	},
	"encoding": func() TemplateFuncsEncoding {
		return TemplateFuncsEncoding{}
	},
	"requests": func() TemplateFuncsRequests {
		return TemplateFuncsRequests{}
	},
}

type GoTemplateBackend struct {
	template *template.Template
}

func NewGoTemplateBackend(source string) (*GoTemplateBackend, error) {
	tmpl := template.New("template")
	tmpl.Funcs(funcMap)

	if _, err := tmpl.Parse(source); err != nil {
		return nil, err
	}

	return &GoTemplateBackend{template: tmpl}, nil
}

type goTemplateContext struct {
	Request *types.Request
}

func (b *GoTemplateBackend) Handle(request *types.Request) ([]byte, error) {
	context := &goTemplateContext{Request: request}

	var buffer bytes.Buffer
	if err := b.template.Execute(&buffer, context); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
