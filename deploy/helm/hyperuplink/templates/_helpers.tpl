{{/*
Expand the name of the chart.
*/}}
{{- define "hyperuplink.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this
(by the DNS naming spec).
*/}}
{{- define "hyperuplink.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "hyperuplink.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels.
*/}}
{{- define "hyperuplink.labels" -}}
helm.sh/chart: {{ include "hyperuplink.chart" . }}
{{ include "hyperuplink.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- with .Values.commonLabels }}
{{ toYaml . }}
{{- end }}
{{- end }}

{{/*
Selector labels.
*/}}
{{- define "hyperuplink.selectorLabels" -}}
app.kubernetes.io/name: {{ include "hyperuplink.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use.
*/}}
{{- define "hyperuplink.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "hyperuplink.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Fully qualified container image reference.
*/}}
{{- define "hyperuplink.image" -}}
{{- $tag := default .Chart.AppVersion .Values.image.tag -}}
{{- printf "%s:%s" .Values.image.repository $tag -}}
{{- end }}

{{/*
Name of the Secret that holds the config file. Points at an out-of-band Secret
when `existingConfigSecret` is set, otherwise at the chart-rendered one.
*/}}
{{- define "hyperuplink.configSecretName" -}}
{{- if .Values.existingConfigSecret }}
{{- .Values.existingConfigSecret }}
{{- else }}
{{- printf "%s-config" (include "hyperuplink.fullname" .) }}
{{- end }}
{{- end }}

{{/*
Name of the PVC backing media storage. Points at an existing claim when set.
*/}}
{{- define "hyperuplink.pvcName" -}}
{{- if .Values.persistence.existingClaim }}
{{- .Values.persistence.existingClaim }}
{{- else }}
{{- printf "%s-media" (include "hyperuplink.fullname" .) }}
{{- end }}
{{- end }}

{{/*
Absolute path to the config file inside the container.
*/}}
{{- define "hyperuplink.configFilePath" -}}
{{- printf "%s/%s" (.Values.configMountPath | trimSuffix "/") .Values.configFileKey }}
{{- end }}
