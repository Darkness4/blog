{{- range .Lvl0s }}
{{ $records := index $.GroupedRecords . }}
<section>
  <div>{{ . }}</div>
  <ul role="listbox" style="flex-direction: column; justify-content: start; align-items: start;">
    {{- range $records }}
    <li style="display: flex; list-style: none;">
      <a class="search-result" aria-label="Link to the result" href="{{ .URL }}">
      {{ .Formatted.HierarchyLvl1 | noescape }}{{- if .HierarchyLvl2 }} &rsaquo; {{ .Formatted.HierarchyLvl2 | noescape }}{{- end }}{{- if .HierarchyLvl3 }} &rsaquo; {{ .Formatted.HierarchyLvl3 | noescape }}{{- end }}{{- if .HierarchyLvl4 }} &rsaquo; {{ .Formatted.HierarchyLvl4 | noescape }}{{- end }}{{- if .HierarchyLvl5 }} &rsaquo; {{ .Formatted.HierarchyLvl5 | noescape }}{{- end }}{{- if .HierarchyLvl6 }} &rsaquo; {{ .Formatted.HierarchyLvl6 | noescape }}{{- end }}
      </a>
    </li>
    {{- end }}
  </ul>
</section>
{{- end }}
