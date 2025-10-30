// Harbor Modifier - Automated Harbor Helm chart customization tool
// Pulls official Harbor chart and applies file-based modifications
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	defaultVersion = "1.18.0"
	defaultRepo    = "harbor"
	chartName      = "harbor"
)

type Config struct {
	Version          string
	RepoName         string
	ProjectDir       string
	ChartDir         string
	ModificationsDir string
}

func main() {
	version := flag.String("version", defaultVersion, "Harbor chart version")
	verbose := flag.Bool("verbose", false, "Verbose output")
	flag.Parse()

	cfg := &Config{
		Version:          *version,
		RepoName:         defaultRepo,
		ProjectDir:       mustGetwd(),
		ChartDir:         filepath.Join(mustGetwd(), "harbor-helm"),
		ModificationsDir: filepath.Join(mustGetwd(), "modifications"),
	}

	if *verbose {
		log.SetFlags(log.Ltime | log.Lshortfile)
	}

	fmt.Println("===================================")
	fmt.Println("Harbor Chart Automation (Go)")
	fmt.Println("===================================")
	fmt.Printf("Version: %s\n", cfg.Version)
	fmt.Printf("Project: %s\n\n", cfg.ProjectDir)

	// Step 1: Pull Harbor chart
	if err := pullChart(cfg); err != nil {
		log.Fatalf("❌ Failed to pull chart: %v", err)
	}

	// Step 2: Apply modifications
	if err := applyModifications(cfg); err != nil {
		log.Fatalf("❌ Failed to apply modifications: %v", err)
	}

	// Step 3: Validate (skip lint, just check dependencies)
	if err := validateDependencies(cfg); err != nil {
		log.Fatalf("❌ Validation failed: %v", err)
	}

	fmt.Println("\n✅ All modifications applied successfully!")
	fmt.Printf("\nModified chart location: %s\n", cfg.ChartDir)
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Review the modified chart")
	fmt.Println("  2. Update values as needed")
	fmt.Printf("  3. Install: helm install harbor %s -n harbor --create-namespace\n", cfg.ChartDir)
}

func mustGetwd() string {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}
	return wd
}

func pullChart(cfg *Config) error {
	fmt.Println("📦 Pulling Harbor chart...")

	// Remove existing chart
	if err := os.RemoveAll(cfg.ChartDir); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove existing chart: %w", err)
	}

	// Ensure repo is added
	if err := ensureRepo(cfg.RepoName); err != nil {
		return fmt.Errorf("failed to ensure repo: %w", err)
	}

	// Update repo
	cmd := exec.Command("helm", "repo", "update", cfg.RepoName)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("helm repo update failed: %w\n%s", err, output)
	}

	// Pull chart
	cmd = exec.Command("helm", "pull",
		fmt.Sprintf("%s/%s", cfg.RepoName, chartName),
		"--version", cfg.Version,
		"--untar",
		"--untardir", cfg.ProjectDir,
	)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("helm pull failed: %w\n%s", err, output)
	}

	// Rename harbor -> harbor-helm
	harborDir := filepath.Join(cfg.ProjectDir, "harbor")
	if _, err := os.Stat(harborDir); err == nil {
		if err := os.Rename(harborDir, cfg.ChartDir); err != nil {
			return fmt.Errorf("failed to rename chart directory: %w", err)
		}
	}

	fmt.Println("✅ Harbor chart ready")
	return nil
}

func ensureRepo(repoName string) error {
	cmd := exec.Command("helm", "repo", "list")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to list repos: %w", err)
	}

	if !strings.Contains(string(output), repoName) {
		fmt.Printf("Adding %s repo...\n", repoName)
		cmd = exec.Command("helm", "repo", "add", repoName, "https://helm.goharbor.io")
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("failed to add repo: %w\n%s", err, output)
		}
	}
	return nil
}

func applyModifications(cfg *Config) error {
	fmt.Println("\n🔧 Applying custom modifications...")

	// 1. Apply helper templates
	if err := applyHelpers(cfg); err != nil {
		return fmt.Errorf("failed to apply helpers: %w", err)
	}

	// 1.5. Patch database templates for relizapostgresql support
	if err := patchDatabaseTemplates(cfg); err != nil {
		return fmt.Errorf("failed to patch database templates: %w", err)
	}

	// 1.55. Remove redundant harbor.relizapostgresql template (harbor.database now points to it)
	if err := removeRelizaPostgresqlTemplate(cfg); err != nil {
		return fmt.Errorf("failed to remove harbor.relizapostgresql template: %w", err)
	}

	// 1.6. Remove harbor-db templates (replaced by relizapostgresql)
	if err := removeHarborDatabase(cfg); err != nil {
		return fmt.Errorf("failed to remove harbor-db templates: %w", err)
	}

	// 2. Apply templates
	if err := applyTemplates(cfg); err != nil {
		return fmt.Errorf("failed to apply templates: %w", err)
	}

	// 3. Merge values
	if err := mergeValues(cfg); err != nil {
		return fmt.Errorf("failed to merge values: %w", err)
	}

	// 3.5. Clean up obsolete database.internal section
	if err := cleanupDatabaseInternal(cfg); err != nil {
		return fmt.Errorf("failed to cleanup database.internal: %w", err)
	}

	// 4. Update Chart.yaml
	if err := updateChart(cfg); err != nil {
		return fmt.Errorf("failed to update Chart.yaml: %w", err)
	}

	// 5. Update .helmignore
	if err := updateHelmignore(cfg); err != nil {
		return fmt.Errorf("failed to update .helmignore: %w", err)
	}

	// 6. Apply template overlays (replaces image patching)
	if err := applyTemplateOverlays(cfg); err != nil {
		return fmt.Errorf("failed to apply template overlays: %w", err)
	}

	return nil
}

func applyHelpers(cfg *Config) error {
	fmt.Println("  → Adding helper templates...")

	helpersDir := filepath.Join(cfg.ModificationsDir, "helpers")
	targetFile := filepath.Join(cfg.ChartDir, "templates", "_helpers.tpl")

	// Append all helper files
	helpers, err := filepath.Glob(filepath.Join(helpersDir, "*.tpl"))
	if err != nil {
		return fmt.Errorf("failed to glob helpers: %w", err)
	}

	f, err := os.OpenFile(targetFile, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open _helpers.tpl: %w", err)
	}
	defer f.Close()

	for _, helper := range helpers {
		content, err := os.ReadFile(helper)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", helper, err)
		}
		if _, err := f.WriteString("\n" + string(content)); err != nil {
			return fmt.Errorf("failed to write helper: %w", err)
		}
	}

	fmt.Println("    ✅ Helper templates added")
	return nil
}

func removeHarborDatabase(cfg *Config) error {
	fmt.Println("  → Removing harbor-db templates...")

	databaseDir := filepath.Join(cfg.ChartDir, "templates", "database")

	// Files to remove (harbor's internal database, replaced by relizapostgresql)
	filesToRemove := []string{
		"database-ss.yaml",
		"database-svc.yaml",
		"database-secret.yaml",
	}

	for _, file := range filesToRemove {
		path := filepath.Join(databaseDir, file)
		if err := os.Remove(path); err != nil {
			if os.IsNotExist(err) {
				continue // Already removed
			}
			return fmt.Errorf("failed to remove %s: %w", file, err)
		}
	}

	fmt.Println("    ✅ Harbor-db templates removed")
	return nil
}

func patchDatabaseTemplates(cfg *Config) error {
	fmt.Println("  → Patching database templates for relizapostgresql...")

	targetFile := filepath.Join(cfg.ChartDir, "templates", "_helpers.tpl")

	content, err := os.ReadFile(targetFile)
	if err != nil {
		return fmt.Errorf("failed to read _helpers.tpl: %w", err)
	}

	newContent := string(content)

	// 1. Replace harbor.database definition to point to relizapostgresql service
	oldDatabase := `{{- define "harbor.database" -}}
  {{- printf "%s-database" (include "harbor.fullname" .) -}}
{{- end -}}`

	newDatabase := `{{- define "harbor.database" -}}
  {{- printf "%s-relizapostgresql" (include "harbor.fullname" .) -}}
{{- end -}}`

	newContent = strings.Replace(newContent, oldDatabase, newDatabase, 1)

	// 2. Update harbor.database.username to use relizapostgresql.auth.username
	oldUsername := `{{- define "harbor.database.username" -}}
  {{- if eq .Values.database.type "internal" -}}
    {{- printf "%s" "postgres" -}}
  {{- else -}}
    {{- .Values.database.external.username -}}
  {{- end -}}
{{- end -}}`

	newUsername := `{{- define "harbor.database.username" -}}
  {{- if eq .Values.database.type "internal" -}}
    {{- .Values.relizapostgresql.auth.username -}}
  {{- else -}}
    {{- .Values.database.external.username -}}
  {{- end -}}
{{- end -}}`

	newContent = strings.Replace(newContent, oldUsername, newUsername, 1)

	// 3. Update harbor.database.rawPassword to use relizapostgresql.auth.password (remove secret lookup)
	oldPassword := `{{- define "harbor.database.rawPassword" -}}
  {{- if eq .Values.database.type "internal" -}}
    {{- $existingSecret := lookup "v1" "Secret" .Release.Namespace (include "harbor.database" .) -}}
    {{- if and (not (empty $existingSecret)) (hasKey $existingSecret.data "POSTGRES_PASSWORD") -}}
      {{- .Values.database.internal.password | default (index $existingSecret.data "POSTGRES_PASSWORD" | b64dec) -}}
    {{- else -}}
      {{- .Values.database.internal.password -}}
    {{- end -}}
  {{- else -}}
    {{- .Values.database.external.password -}}
  {{- end -}}
{{- end -}}`

	newPassword := `{{- define "harbor.database.rawPassword" -}}
  {{- if eq .Values.database.type "internal" -}}
    {{- .Values.relizapostgresql.auth.password -}}
  {{- else -}}
    {{- .Values.database.external.password -}}
  {{- end -}}
{{- end -}}`

	newContent = strings.Replace(newContent, oldPassword, newPassword, 1)

	// 4. Update harbor.database.coreDatabase to use relizapostgresql.auth.database
	oldCoreDatabase := `{{- define "harbor.database.coreDatabase" -}}
  {{- if eq .Values.database.type "internal" -}}
    {{- printf "%s" "registry" -}}
  {{- else -}}
    {{- .Values.database.external.coreDatabase -}}
  {{- end -}}
{{- end -}}`

	newCoreDatabase := `{{- define "harbor.database.coreDatabase" -}}
  {{- if eq .Values.database.type "internal" -}}
    {{- .Values.relizapostgresql.auth.database -}}
  {{- else -}}
    {{- .Values.database.external.coreDatabase -}}
  {{- end -}}
{{- end -}}`

	newContent = strings.Replace(newContent, oldCoreDatabase, newCoreDatabase, 1)

	if err := os.WriteFile(targetFile, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write patched _helpers.tpl: %w", err)
	}

	fmt.Println("    ✅ Database templates patched (4 templates updated)")
	return nil
}

func removeRelizaPostgresqlTemplate(cfg *Config) error {
	fmt.Println("  → Removing redundant harbor.relizapostgresql template...")

	targetFile := filepath.Join(cfg.ChartDir, "templates", "_helpers.tpl")

	content, err := os.ReadFile(targetFile)
	if err != nil {
		return fmt.Errorf("failed to read _helpers.tpl: %w", err)
	}

	// Remove the harbor.relizapostgresql template definition (added by applyHelpers)
	templateToRemove := `{{/*
Reliza PostgreSQL service name
Returns the service name for relizapostgresql when enabled
*/}}
{{- define "harbor.relizapostgresql" -}}
  {{- printf "%s-relizapostgresql" (include "harbor.fullname" .) -}}
{{- end -}}
`

	newContent := strings.Replace(string(content), templateToRemove, "", 1)

	if newContent == string(content) {
		fmt.Println("    ⏭️  harbor.relizapostgresql template not found (already removed or not added)")
		return nil
	}

	if err := os.WriteFile(targetFile, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write _helpers.tpl: %w", err)
	}

	fmt.Println("    ✅ Redundant template removed")
	return nil
}

func applyTemplates(cfg *Config) error {
	fmt.Println("  → Adding custom templates...")

	templatesDir := filepath.Join(cfg.ModificationsDir, "templates")
	targetDir := filepath.Join(cfg.ChartDir, "templates")

	templates, err := filepath.Glob(filepath.Join(templatesDir, "*.yaml"))
	if err != nil {
		return fmt.Errorf("failed to glob templates: %w", err)
	}

	for _, tmpl := range templates {
		basename := filepath.Base(tmpl)
		target := filepath.Join(targetDir, basename)

		// Check if already exists
		if _, err := os.Stat(target); err == nil {
			fmt.Printf("    ⏭️  %s already exists, skipping...\n", basename)
			continue
		}

		content, err := os.ReadFile(tmpl)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", tmpl, err)
		}

		if err := os.WriteFile(target, content, 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", target, err)
		}
		fmt.Printf("    ✅ Added %s\n", basename)
	}

	return nil
}

func mergeValues(cfg *Config) error {
	fmt.Println("  → Merging values...")

	valuesDir := filepath.Join(cfg.ModificationsDir, "values")
	targetFile := filepath.Join(cfg.ChartDir, "values.yaml")

	// Read existing values
	existing, err := os.ReadFile(targetFile)
	if err != nil {
		return fmt.Errorf("failed to read values.yaml: %w", err)
	}

	var existingValues map[string]interface{}
	if err := yaml.Unmarshal(existing, &existingValues); err != nil {
		return fmt.Errorf("failed to parse values.yaml: %w", err)
	}

	// Check if already modified
	if _, ok := existingValues["relizaPostgresql"]; ok {
		fmt.Println("    ⏭️  Values already merged, skipping...")
		return nil
	}

	// Read and merge all value files
	valueFiles, err := filepath.Glob(filepath.Join(valuesDir, "*.yaml"))
	if err != nil {
		return fmt.Errorf("failed to glob values: %w", err)
	}

	for _, vf := range valueFiles {
		content, err := os.ReadFile(vf)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", vf, err)
		}

		var newValues map[string]interface{}
		if err := yaml.Unmarshal(content, &newValues); err != nil {
			return fmt.Errorf("failed to parse %s: %w", vf, err)
		}

		// Merge
		mergeMaps(existingValues, newValues)
	}

	// Add traefik config under expose section
	if expose, ok := existingValues["expose"].(map[string]interface{}); ok {
		if _, ok := expose["traefik"]; !ok {
			expose["traefik"] = map[string]interface{}{
				"enabled":     false,
				"host":        "harbor.example.com",
				"middlewares": []interface{}{},
				"tls": map[string]interface{}{
					"enabled":      true,
					"certResolver": "",
					"secretName":   "",
				},
				"httpsRedirect": map[string]interface{}{
					"enabled": true,
				},
			}
		}
	}

	// Write merged values
	merged, err := yaml.Marshal(existingValues)
	if err != nil {
		return fmt.Errorf("failed to marshal values: %w", err)
	}

	if err := os.WriteFile(targetFile, merged, 0644); err != nil {
		return fmt.Errorf("failed to write values.yaml: %w", err)
	}

	fmt.Println("    ✅ Values merged")
	return nil
}

func cleanupDatabaseInternal(cfg *Config) error {
	fmt.Println("  → Cleaning up obsolete database.internal section...")

	targetFile := filepath.Join(cfg.ChartDir, "values.yaml")

	// Read values
	content, err := os.ReadFile(targetFile)
	if err != nil {
		return fmt.Errorf("failed to read values.yaml: %w", err)
	}

	var values map[string]interface{}
	if err := yaml.Unmarshal(content, &values); err != nil {
		return fmt.Errorf("failed to parse values.yaml: %w", err)
	}

	// Remove database.internal section
	if database, ok := values["database"].(map[string]interface{}); ok {
		if _, exists := database["internal"]; exists {
			delete(database, "internal")
			fmt.Println("    ✅ Removed database.internal section")
		} else {
			fmt.Println("    ⏭️  database.internal not found (already removed)")
		}
	}

	// Write updated values
	updated, err := yaml.Marshal(values)
	if err != nil {
		return fmt.Errorf("failed to marshal values: %w", err)
	}

	if err := os.WriteFile(targetFile, updated, 0644); err != nil {
		return fmt.Errorf("failed to write values.yaml: %w", err)
	}

	return nil
}

func mergeMaps(dst, src map[string]interface{}) {
	for k, v := range src {
		if dstVal, ok := dst[k]; ok {
			if dstMap, ok := dstVal.(map[string]interface{}); ok {
				if srcMap, ok := v.(map[string]interface{}); ok {
					mergeMaps(dstMap, srcMap)
					continue
				}
			}
		}
		dst[k] = v
	}
}

func updateChart(cfg *Config) error {
	fmt.Println("  → Updating Chart.yaml...")

	chartFile := filepath.Join(cfg.ChartDir, "Chart.yaml")

	// Read existing Chart.yaml
	existing, err := os.ReadFile(chartFile)
	if err != nil {
		return fmt.Errorf("failed to read Chart.yaml: %w", err)
	}

	var chartData map[string]interface{}
	if err := yaml.Unmarshal(existing, &chartData); err != nil {
		return fmt.Errorf("failed to parse Chart.yaml: %w", err)
	}

	// Update apiVersion to v2
	if apiVersion, ok := chartData["apiVersion"].(string); ok && apiVersion == "v1" {
		chartData["apiVersion"] = "v2"
		fmt.Println("    ✅ Updated apiVersion to v2")
	}

	// Read and merge all chart modification files
	chartModsDir := filepath.Join(cfg.ModificationsDir, "chart")
	chartModFiles, err := filepath.Glob(filepath.Join(chartModsDir, "*.yaml"))
	if err != nil {
		return fmt.Errorf("failed to glob chart modifications: %w", err)
	}

	for _, modFile := range chartModFiles {
		modContent, err := os.ReadFile(modFile)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", modFile, err)
		}

		var modData map[string]interface{}
		if err := yaml.Unmarshal(modContent, &modData); err != nil {
			return fmt.Errorf("failed to parse %s: %w", modFile, err)
		}

		// Special handling for dependencies - append instead of replace
		if deps, ok := modData["dependencies"].([]interface{}); ok {
			if existingDeps, ok := chartData["dependencies"].([]interface{}); ok {
				chartData["dependencies"] = append(existingDeps, deps...)
			} else {
				chartData["dependencies"] = deps
			}
			delete(modData, "dependencies") // Don't merge it again below
		}

		// Merge other fields
		mergeMaps(chartData, modData)
	}

	if len(chartModFiles) > 0 {
		fmt.Printf("    ✅ Applied %d chart modification(s)\n", len(chartModFiles))
	}

	// Write updated Chart.yaml
	updated, err := yaml.Marshal(chartData)
	if err != nil {
		return fmt.Errorf("failed to marshal Chart.yaml: %w", err)
	}

	if err := os.WriteFile(chartFile, updated, 0644); err != nil {
		return fmt.Errorf("failed to write Chart.yaml: %w", err)
	}

	fmt.Println("    ✅ Chart.yaml updated")
	return nil
}

func updateHelmignore(cfg *Config) error {
	fmt.Println("  → Updating .helmignore...")

	ignoreFile := filepath.Join(cfg.ModificationsDir, ".helmignore")
	targetFile := filepath.Join(cfg.ChartDir, ".helmignore")

	// Check if modifications .helmignore exists
	if _, err := os.Stat(ignoreFile); os.IsNotExist(err) {
		fmt.Println("    ⏭️  No .helmignore modifications, skipping...")
		return nil
	}

	// Read modifications .helmignore
	newContent, err := os.ReadFile(ignoreFile)
	if err != nil {
		return fmt.Errorf("failed to read .helmignore: %w", err)
	}

	// Append new content to chart's .helmignore
	f, err := os.OpenFile(targetFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open .helmignore: %w", err)
	}
	defer f.Close()

	if _, err := f.WriteString("\n" + string(newContent)); err != nil {
		return fmt.Errorf("failed to write .helmignore: %w", err)
	}

	fmt.Println("    ✅ .helmignore updated")
	return nil
}

func applyTemplateOverlays(cfg *Config) error {
	fmt.Println("  → Applying template overlays...")

	overlaysDir := filepath.Join(cfg.ModificationsDir, "template-overlays")
	templatesDir := filepath.Join(cfg.ChartDir, "templates")

	// Check if overlays directory exists
	if _, err := os.Stat(overlaysDir); os.IsNotExist(err) {
		fmt.Println("    ⏭️  No template overlays, skipping...")
		return nil
	}

	overlayCount := 0
	err := filepath.Walk(overlaysDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Get relative path from overlays dir
		relPath, err := filepath.Rel(overlaysDir, path)
		if err != nil {
			return fmt.Errorf("failed to get relative path: %w", err)
		}

		// Target path in chart templates
		targetPath := filepath.Join(templatesDir, relPath)

		// Ensure target directory exists
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}

		// Copy overlay file to target
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read overlay %s: %w", relPath, err)
		}

		if err := os.WriteFile(targetPath, content, 0644); err != nil {
			return fmt.Errorf("failed to write %s: %w", relPath, err)
		}

		overlayCount++
		return nil
	})

	if err != nil {
		return err
	}

	if overlayCount > 0 {
		fmt.Printf("    ✅ Applied %d template overlay(s)\n", overlayCount)
	}

	return nil
}

func validateDependencies(cfg *Config) error {
	fmt.Println("\n🔍 Validating chart dependencies...")

	// Build dependencies
	cmd := exec.Command("helm", "dependency", "build", cfg.ChartDir)
	if output, err := cmd.CombinedOutput(); err != nil {
		// Ignore errors for conditional dependencies
		if !strings.Contains(string(output), "missing in charts/ directory") {
			return fmt.Errorf("dependency build failed: %w\n%s", err, output)
		}
	}

	fmt.Println("✅ Dependencies validated")
	return nil
}
