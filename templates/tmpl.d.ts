{{range .Specs}}
declare interface {{.Name}} {
    {{range .Fields}}
        {{.Name}}: {{.Type}}
    {{end}}
}

{{end}}