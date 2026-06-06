package route

import "github.com/orisan/review/internal/model"

func RoutesForFindings(findings []model.Finding) []model.ReviewRoute {
	seen := map[model.ReviewRoute]bool{}
	var routes []model.ReviewRoute
	add := func(route model.ReviewRoute) {
		if !seen[route] {
			seen[route] = true
			routes = append(routes, route)
		}
	}

	for _, finding := range findings {
		switch finding.Category {
		case "auth_logic_changed", "authorization_weakened", "validation_removed", "tls_verification_disabled":
			add(model.RouteAppSec)
		case "secret_like_value_added":
			add(model.RouteAppSec)
			add(model.RouteHuman)
		case "ci_permissions_broadened":
			add(model.RouteCICD)
			add(model.RouteAppSec)
		case "unpinned_github_action":
			add(model.RouteCICD)
		case "dependency_manifest_changed":
			add(model.RouteDependency)
		case "infra_public_exposure":
			add(model.RouteInfra)
			add(model.RouteAppSec)
		case "destructive_migration", "tests_skipped", "ai_generated_marker", "binary_file_change":
			add(model.RouteHuman)
		}
	}
	return routes
}
