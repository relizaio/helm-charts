{{/*
Reliza customization: Smart image reference builder
Handles both standard deployments and reliza-cd tag replacement

Usage: {{ include "harbor.imageRef" (dict "repository" .Values.core.image.repository "tag" .Values.core.image.tag) }}
With digest: {{ include "harbor.imageRef" (dict "repository" .Values.core.image.repository "tag" .Values.core.image.tag "digest" .Values.imageDigests.core) }}
*/}}
{{- define "harbor.imageRef" -}}
{{- $repo := .repository -}}
{{- $tag := .tag -}}
{{- $digest := .digest | default "" -}}
{{- if contains ":" $repo -}}
  {{/* Repository already contains tag/digest (reliza-cd format), use as-is */}}
  {{- $repo -}}
{{- else -}}
  {{/* Standard format: build from repository + tag + optional digest */}}
  {{- $repo -}}:{{- $tag -}}
  {{- if $digest -}}@{{- $digest -}}{{- end -}}
{{- end -}}
{{- end -}}
