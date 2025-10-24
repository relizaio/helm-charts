{{/*
Reliza customization: Standardized common labels
These labels work correctly whether Harbor is used standalone or as a subchart
*/}}
{{- define "harbor.common.labels" -}}
app.kubernetes.io/name: {{ include "harbor.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
helm.sh/chart: {{ include "harbor.chart" . }}
{{- end -}}

{{/*
Reliza customization: Component-specific labels
Adds component label while maintaining compatibility
*/}}
{{- define "harbor.component.labels" -}}
{{ include "harbor.common.labels" . }}
app.kubernetes.io/component: {{ .component }}
{{- end -}}

{{/*
Reliza customization: Selector labels
Use only stable labels for selectors (no version/chart)
*/}}
{{- define "harbor.selector.labels" -}}
app.kubernetes.io/name: {{ include "harbor.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- if .component }}
app.kubernetes.io/component: {{ .component }}
{{- end }}
{{- end -}}
