package main

import (
	"bytes"
	"strings"
	"text/template"
)

var httpTemplate = `
{{$svrType := .ServiceType}}
{{$svrName := .ServiceName}}
type {{.ServiceType}}HTTPServer interface {
{{- range .MethodSets}}
	{{.Name}}(context.Context, *{{.Request}}) (*{{.Reply}}, error)
{{- end}}
}

type {{.ServiceType}}HTTPHandler struct {
	srv {{.ServiceType}}HTTPServer
}

func new{{.ServiceType}}HTTPHandler(s {{.ServiceType}}HTTPServer) *{{.ServiceType}}HTTPHandler {
	return &{{.ServiceType}}HTTPHandler{srv: s}
}

{{- range .MethodSets}}

func (h *{{$svrType}}HTTPHandler) {{.Name}}(req *go_restful.Request, resp *go_restful.Response) {
	in := &{{.Request}}{}

	{{- if .HasBody}}
	if err := transportHTTP.GetBody(req, &in); err != nil {
		resp.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
	{{- else}}
	if err := transportHTTP.GetQuery(req, &in); err != nil {
		resp.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
	{{- end}}
	{{- if .HasVars}}
	if err := transportHTTP.GetPathValue(req, &in); err != nil {
		resp.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
	{{- end}}

	out,err := h.srv.{{.Name}}(req.Request.Context(),in)
	if err != nil {
		resp.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}

	result, err := json.Marshal(out)
	if err != nil {
		resp.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
	_, err = resp.Write(result)
	if err != nil {
		resp.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
}

{{- end}}

func Register{{.ServiceType}}HTTPServer(container *go_restful.Container, srv {{.ServiceType}}HTTPServer) {
	handler := new{{.ServiceType}}HTTPHandler(srv)

	ws := new(go_restful.WebService)
	ws.ApiVersion("{{.ServiceVersion}}")
	ws.Path("/{{.ServiceVersion}}").Produces(go_restful.MIME_JSON)

	{{- range .Methods}}
	ws.Route(ws.{{.Method}}("{{.Path}}").
		To(handler.{{.Name}}))
	{{- end}}

	container.Add(ws)
}
`

type serviceDesc struct {
	ServiceType    string // Greeter
	ServiceName    string // helloworld.Greeter
	Metadata       string // api/helloworld/helloworld.proto
	ServiceVersion string // v1
	Methods        []*methodDesc
	MethodSets     map[string]*methodDesc
}

type methodDesc struct {
	// method
	Name    string
	Num     int
	Request string
	Reply   string
	// http_rule
	Path         string
	Method       string
	HasVars      bool
	HasBody      bool
	Body         string
	ResponseBody string
}

func (s *serviceDesc) execute() string {
	s.MethodSets = make(map[string]*methodDesc)
	for _, m := range s.Methods {
		s.MethodSets[m.Name] = m
	}
	buf := new(bytes.Buffer)
	tmpl, err := template.New("http").Parse(strings.TrimSpace(httpTemplate))
	if err != nil {
		panic(err)
	}
	if err := tmpl.Execute(buf, s); err != nil {
		panic(err)
	}
	return strings.Trim(buf.String(), "\r\n")
}
