package report

import (
	"bytes"
	"fmt"
	"html"
	"strings"

	"github.com/orisan/review/internal/app"
	"github.com/orisan/review/internal/model"
	"github.com/orisan/review/internal/route"
)

type HTML struct{}

func (HTML) Render(result model.ReviewResult) ([]byte, error) {
	var out bytes.Buffer
	fmt.Fprintln(&out, "<!doctype html>")
	fmt.Fprintln(&out, `<html lang="en">`)
	fmt.Fprintln(&out, "<head>")
	fmt.Fprintln(&out, `<meta charset="utf-8">`)
	fmt.Fprintln(&out, `<meta name="viewport" content="width=device-width, initial-scale=1">`)
	fmt.Fprintf(&out, "<title>%s Report</title>\n", esc(app.Name))
	writeHTMLStyle(&out)
	fmt.Fprintln(&out, "</head>")
	fmt.Fprintln(&out, "<body>")
	fmt.Fprintln(&out, `<main class="page">`)
	fmt.Fprintf(&out, "<p class=\"eyebrow\">Generated locally by %s %s</p>\n", esc(app.Name), esc(app.Version))
	fmt.Fprintln(&out, "<h1>Orisan Review</h1>")
	writeHTMLSummary(&out, result)
	writeHTMLRoutes(&out, result.Summary.Routes)
	writeHTMLFindings(&out, result)
	writeHTMLSafetyNote(&out)
	fmt.Fprintln(&out, "</main>")
	fmt.Fprintln(&out, "</body>")
	fmt.Fprintln(&out, "</html>")
	return out.Bytes(), nil
}

func writeHTMLStyle(out *bytes.Buffer) {
	fmt.Fprintln(out, `<style>
:root {
  color-scheme: light;
  --bg: #f7f8fa;
  --panel: #ffffff;
  --text: #18202a;
  --muted: #5d6978;
  --border: #d9dee7;
  --critical: #a21414;
  --high: #b94a00;
  --medium: #8a6500;
  --ok: #17623a;
}
* { box-sizing: border-box; }
body {
  margin: 0;
  background: var(--bg);
  color: var(--text);
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
  line-height: 1.45;
}
.page { max-width: 1120px; margin: 0 auto; padding: 40px 24px; }
.eyebrow { margin: 0 0 8px; color: var(--muted); font-size: 13px; }
h1 { margin: 0 0 24px; font-size: 34px; }
h2 { margin: 28px 0 12px; font-size: 20px; }
.summary, .panel, .finding {
  background: var(--panel);
  border: 1px solid var(--border);
  border-radius: 8px;
}
.summary {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 1px;
  overflow: hidden;
}
.metric { padding: 16px; border-right: 1px solid var(--border); border-bottom: 1px solid var(--border); }
.label { color: var(--muted); display: block; font-size: 12px; text-transform: uppercase; }
.value { display: block; font-size: 18px; font-weight: 700; margin-top: 4px; }
.yes { color: var(--critical); }
.no { color: var(--ok); }
.routes { display: flex; flex-wrap: wrap; gap: 8px; margin: 0; padding: 0; list-style: none; }
.route, .severity {
  display: inline-block;
  border: 1px solid var(--border);
  border-radius: 999px;
  padding: 4px 10px;
  font-size: 13px;
  font-weight: 650;
}
.finding { padding: 18px; margin-bottom: 14px; }
.finding h3 { margin: 0 0 12px; font-size: 18px; }
.meta {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 8px 14px;
  margin-bottom: 12px;
}
.meta div, .evidence, .reason { min-width: 0; }
code {
  background: #f0f2f5;
  border-radius: 4px;
  padding: 2px 5px;
  word-break: break-word;
}
.critical { color: var(--critical); }
.high { color: var(--high); }
.medium { color: var(--medium); }
.note { padding: 16px; color: var(--muted); }
</style>`)
}

func writeHTMLSummary(out *bytes.Buffer, result model.ReviewResult) {
	securityRequired := result.Summary.Decision != model.DecisionPass
	fmt.Fprintln(out, `<section class="summary" aria-label="Review summary">`)
	writeMetric(out, "Security review required", yesNo(securityRequired), boolClass(securityRequired))
	writeMetric(out, "Risk level", stringOrDefault(result.Summary.RiskLevel, "NONE"), "")
	writeMetric(out, "Review decision", string(result.Summary.Decision), "")
	writeMetric(out, "Findings", fmt.Sprintf("%d", len(result.Findings)), "")
	writeMetric(out, "Changed files", fmt.Sprintf("%d", result.Files.ChangedFiles), "")
	writeMetric(out, "payload_stored", "false", "no")
	fmt.Fprintln(out, "</section>")
}

func writeMetric(out *bytes.Buffer, label, value, class string) {
	classAttr := "value"
	if class != "" {
		classAttr += " " + class
	}
	fmt.Fprintf(out, "<div class=\"metric\"><span class=\"label\">%s</span><span class=\"%s\">%s</span></div>\n", esc(label), esc(classAttr), esc(value))
}

func writeHTMLRoutes(out *bytes.Buffer, routes []model.ReviewRoute) {
	fmt.Fprintln(out, "<h2>Reviewer routing</h2>")
	if len(routes) == 0 {
		fmt.Fprintln(out, `<p class="panel note">No reviewer route required.</p>`)
		return
	}
	fmt.Fprintln(out, `<ul class="routes">`)
	for _, route := range routes {
		fmt.Fprintf(out, "<li class=\"route\">%s</li>\n", esc(routeLabel(route)))
	}
	fmt.Fprintln(out, "</ul>")
}

func writeHTMLFindings(out *bytes.Buffer, result model.ReviewResult) {
	fmt.Fprintln(out, "<h2>Findings</h2>")
	if len(result.Findings) == 0 {
		fmt.Fprintln(out, `<section class="panel note"><strong>No findings.</strong> Security review required: NO.</section>`)
		return
	}
	for i, finding := range result.Findings {
		fmt.Fprintln(out, `<article class="finding">`)
		fmt.Fprintf(out, "<h3>%d. %s - %s</h3>\n", i+1, esc(finding.ID), esc(finding.Title))
		fmt.Fprintln(out, `<div class="meta">`)
		fmt.Fprintf(out, "<div><span class=\"label\">Severity</span><span class=\"severity %s\">%s</span></div>\n", severityClass(finding.Severity), esc(strings.ToUpper(string(finding.Severity))))
		fmt.Fprintf(out, "<div><span class=\"label\">Category</span>%s</div>\n", esc(finding.Category))
		fmt.Fprintf(out, "<div><span class=\"label\">File</span><code>%s</code></div>\n", esc(finding.Location.Path))
		if finding.Location.StartLine != 0 {
			fmt.Fprintf(out, "<div><span class=\"label\">Line</span>%d</div>\n", finding.Location.StartLine)
		}
		fmt.Fprintf(out, "<div><span class=\"label\">Reviewer</span>%s</div>\n", esc(labelsForFinding(finding)))
		fmt.Fprintf(out, "<div><span class=\"label\">payload_stored</span>false</div>\n")
		fmt.Fprintln(out, "</div>")
		fmt.Fprintf(out, "<p class=\"reason\"><span class=\"label\">Reason</span>%s</p>\n", esc(reasonSentence(finding)))
		if finding.Evidence != "" {
			fmt.Fprintf(out, "<p class=\"evidence\"><span class=\"label\">Safe evidence</span><code>%s</code></p>\n", esc(finding.Evidence))
		} else {
			fmt.Fprintln(out, `<p class="evidence"><span class="label">Safe evidence</span><code>redacted snippet only</code></p>`)
		}
		fmt.Fprintln(out, "</article>")
	}
}

func writeHTMLSafetyNote(out *bytes.Buffer) {
	fmt.Fprintln(out, "<h2>Safety note</h2>")
	fmt.Fprintln(out, `<section class="panel note">Generated locally. No cloud calls. No source upload. No full diff stored. Evidence is redacted and every finding has <code>payload_stored=false</code>.</section>`)
}

func labelsForFinding(finding model.Finding) string {
	routes := route.RoutesForFindings([]model.Finding{finding})
	if len(routes) == 0 {
		return "Human review"
	}
	labels := make([]string, 0, len(routes))
	for _, route := range routes {
		labels = append(labels, routeLabel(route))
	}
	return strings.Join(labels, ", ")
}

func boolClass(v bool) string {
	if v {
		return "yes"
	}
	return "no"
}

func severityClass(severity model.Severity) string {
	switch severity {
	case model.SeverityCritical:
		return "critical"
	case model.SeverityHigh:
		return "high"
	case model.SeverityMedium:
		return "medium"
	default:
		return ""
	}
}

func stringOrDefault(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}

func esc(value string) string {
	return html.EscapeString(value)
}
