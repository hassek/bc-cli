package templates

const ProductDetailsTemplate = `
{{repeat "═" 60}}
{{.Name}}
{{repeat "═" 60}}

Price: {{.Currency}} {{.Price}}

Description:
{{wrapAuto .Description}}

{{repeat "═" 60}}
`
