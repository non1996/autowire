package autowire

import (
	"text/template"
)

var (
	factoryTemplate *template.Template
)

const tmplGenFile = `
package {{$.Package}}

import (
	{{range $i := $.Imports.List -}}
	{{if $i.HasAlias}} {{$i.Alias}} {{end}}  "{{$i.PackagePath}}"
	{{end}}
)

{{if $.Components}}
var _ = autowire.Register(
{{- range $c := $.Components -}}
{{template "tmplComponentFactory" $c}}
{{- end}}
)
{{end}}

{{if $.Applications}}
{{range $a := $.Applications}}
{{template "tmplApp" $a}}
{{end}}
{{end}}
`

const tmplComponentFactory = `
autowire.ComponentFactory[{{$.Type.NameWithPkg}}]{
	Name: "{{$.Alias}}",
    Ptr: {{$.Type.Ptr}},
	Primary: {{$.Primary}},
	Configuration: {{$.IsConfiguration}},
    {{if $.Implements -}}
	Implement: []autowire.Type{
		{{range $impl := $.Implements -}}
			autowire.TypeOf[{{$impl.NameWithPkg}}](),	
		{{end}}
	},
    {{- end}}
	{{if $.Condition -}}
	Condition: &autowire.Condition{
		Key: "{{$.Condition.Key}}",
		Value: "{{$.Condition.Value}}",
	},
	{{- end}}
	{{if $.Injectors -}}
	FieldInjectors: []autowire.Injector[{{$.Type.NameWithPkg}}]{
		{{range $ij := $.Injectors -}}{{template "tmplInjector" $ij}}{{- end}}
	},
	{{- end}}
	{{if $.PostConstruct -}}
	PostConstruct: {{if $.PostConstruct.IsMethod}}(*{{$.Type.NameWithPkg}}).{{$.PostConstruct.FuncName}}{{else}}{{$.PostConstruct.FuncName}}{{end}},
	{{end}}
},
{{if $.IsConfiguration}}
{{range $bean := $.Beans}}
autowire.BeanFactory[{{$.Type.NameComplete}}, {{$bean.Type.NameComplete}}]{
	Name: "{{$bean.Alias}}",
	ComponentName: "{{$.Alias}}",
	BuildFunc: func(comp {{$.Type.NameComplete}}) {{$bean.Type.NameComplete}} {
		return comp.{{$bean.Method}}()			
	},
},
{{end}}
{{range $property := $.Properties}}
autowire.PropertyFactory[{{$.Type.NameComplete}}]{
	Scope: "{{$property.Scope}}",
	ComponentName: "{{$.Alias}}",
	BuildFunc: func(comp {{$.Type.NameComplete}}) any {
		return comp.{{$property.Field}}
	},
},
{{end}}
{{end}}
`

const tmplInjector = `
{{- if eq $.Kind "Component" -}}
autowire.ComponentInjector[{{$.CompType.NameWithPkg}}, {{$.Type.NameComplete}}]{
	{{- if $.Qualifier}}Qualifier: "{{$.Qualifier}}",{{end}}
	Required: {{$.Required}},
	InjectFn: func(c *{{$.CompType.NameWithPkg}}, d {{$.Type.NameComplete}}) {
		c.{{$.FieldName}} = d
	},
},
{{- else if eq $.Kind "Value" -}}
autowire.ValueInjector[{{$.CompType.NameWithPkg}}]{
	Key: "{{$.Key}}",
	Scope: "{{$.Scope}}",
	Required: {{$.Required}},
	InjectFn: func(c *{{$.CompType.NameWithPkg}}, value any) {
		autowire.SetValue(&c.{{$.FieldName}}, value)
	},
},
{{- else if eq $.Kind "Env" -}}
autowire.EnvInjector[{{$.CompType.NameWithPkg}}]{
	Key: "{{$.Key}}",
	Required: {{$.Required}},
	InjectFn: func(c *{{$.CompType.NameWithPkg}}, value string) {
		c.{{$.FieldName}} = value
	},
},
{{- end}}
`

const tmplApp = `
func (a *{{$.Type.NameWithPkg}}) Autowire() {
	 autowire.Context().Inject(autowire.ApplicationFactory[{{$.Type.NameWithPkg}}]{
		App: a,
		{{if $.Injectors}}
		Injectors: []autowire.Injector[{{$.Type.NameWithPkg}}]{
			{{range $ij := $.Injectors}}
			{{template "tmplInjector" $ij}}
			{{end}}
		},
		{{end}}
	})
}
`
