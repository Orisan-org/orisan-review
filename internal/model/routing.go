package model

type ReviewRoute string

const (
	RouteHuman      ReviewRoute = "human_review_required"
	RouteAppSec     ReviewRoute = "appsec_review_required"
	RouteInfra      ReviewRoute = "infra_review_required"
	RouteCICD       ReviewRoute = "ci_cd_review_required"
	RouteDependency ReviewRoute = "dependency_review_required"
	RouteData       ReviewRoute = "data_review_required"
)

type ReviewDecision string

const (
	DecisionPass               ReviewDecision = "pass"
	DecisionReviewRequired     ReviewDecision = "review_required"
	DecisionSecurityReview     ReviewDecision = "security_review_required"
	DecisionBlockUntilReviewed ReviewDecision = "block_until_reviewed"
)
