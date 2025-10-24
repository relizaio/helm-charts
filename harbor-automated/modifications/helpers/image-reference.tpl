{{/*
Reliza customization: Generate image reference with optional digest support
Usage: {{ include "harbor.imageReference" (dict "imageRoot" .Values.core.image "digest" .Values.imageDigests.core.digest "context" $) }}
*/}}
{{- define "harbor.imageReference" -}}
{{- $imageRoot := .imageRoot -}}
{{- $digest := .digest -}}
{{- $registry := $imageRoot.registry | default .context.Values.global.registry -}}
{{- $repository := $imageRoot.repository -}}
{{- $tag := $imageRoot.tag | default .context.Chart.AppVersion -}}
{{- if $registry -}}
{{ $registry }}/{{ $repository }}:{{ $tag }}{{- if $digest -}}@{{ $digest }}{{- end -}}
{{- else -}}
{{ $repository }}:{{ $tag }}{{- if $digest -}}@{{ $digest }}{{- end -}}
{{- end -}}
{{- end -}}
