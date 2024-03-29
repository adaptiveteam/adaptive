package strategy

const (
	CreateVision   = "create_vision"
	ViewVision     = "view_vision"
	ViewEditVision = "view_edit_vision"

	CreateStrategyObjective = "create_strategy_objective"
	ViewStrategyObjectives  = "view_strategy_objectives"
	ViewInitiatives         = "view_capability_community_initiatives" // ReadOnly

	CreateCapabilityCommunity = "create_capability_community"
	ViewCapabilityCommunities = "view_capability_communities"

	CreateInitiative          = "create_initiative"
	CreateInitiativeCommunity = "create_initiative_community"

	ViewCapabilityCommunityObjectives = "view_capability_community_objectives"
	ViewAdvocacyObjectives            = "view_advocacy_objectives"
	ViewAdvocacyInitiatives           = "view_advocacy_initiatives"

	AssociateStrategyObjectiveToCapabilityCommunity = "associate_objective_to_capability_community"
	AssociateInitiativeWithInitiativeCommunity      = "associate_initiative_with_initiative_community"

	VisionEvent = "vision"
)
