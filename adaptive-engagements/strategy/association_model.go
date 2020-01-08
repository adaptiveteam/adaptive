package strategy

import "fmt"

type StrategyEntityType string

const (
	StrategyObjectiveEntity           StrategyEntityType = "objective"
	StrategyCapabilityCommunityEntity StrategyEntityType = "capability_community"
	StrategyInitiativeEntity          StrategyEntityType = "initiative"
	StrategyInitiativeCommunityEntity StrategyEntityType = "initiative_community"
	StrategyNoEntity                  StrategyEntityType = ""

	StrategyObjectiveCapabilityCommunityAssociationSourceLabel      = "Capability Community"
	StrategyObjectiveCapabilityCommunityAssociationTargetLabel      = "Allocated Objective"
	StrategyObjectiveCapabilityCommunityAssociationDescriptionLabel = "Allocation Description"

	StrategyObjectivesAssociationSourceLabel      = "What objective does this support?"
	StrategyObjectivesAssociationTargetLabel      = "What objective supports this one?"
	StrategyObjectivesAssociationDescriptionLabel = "How are these related?"

	StrategyObjectiveInitiativeAssociationSourceLabel      = "Objective Enabled"
	StrategyObjectiveInitiativeAssociationTargetLabel      = "Enabling Initiative"
	StrategyObjectiveInitiativeAssociationDescriptionLabel = "Relationship Description"

	StrategyInitiativeInitiativeCommunityAssociationSourceLabel      = "Initiative Community"
	StrategyInitiativeInitiativeCommunityAssociationTargetLabel      = "Allocated Initiative"
	StrategyInitiativeInitiativeCommunityAssociationDescriptionLabel = "Allocation Description"
)

var (
	strategyEntityLabelMap = map[string]StrategyEntityDialog{
		fmt.Sprintf("%s:%s", StrategyCapabilityCommunityEntity, StrategyObjectiveEntity): {
			SourceLabel:      StrategyObjectiveCapabilityCommunityAssociationSourceLabel,
			TargetLabel:      StrategyObjectiveCapabilityCommunityAssociationTargetLabel,
			DescriptionLabel: StrategyObjectiveCapabilityCommunityAssociationDescriptionLabel,
		},
		fmt.Sprintf("%s:%s", StrategyObjectiveEntity, StrategyObjectiveEntity): {
			SourceLabel:      StrategyObjectivesAssociationSourceLabel,
			TargetLabel:      StrategyObjectivesAssociationTargetLabel,
			DescriptionLabel: StrategyObjectivesAssociationDescriptionLabel,
		},
		fmt.Sprintf("%s:%s", StrategyObjectiveEntity, StrategyInitiativeEntity): {
			SourceLabel:      StrategyObjectiveInitiativeAssociationSourceLabel,
			TargetLabel:      StrategyObjectiveInitiativeAssociationTargetLabel,
			DescriptionLabel: StrategyObjectiveInitiativeAssociationDescriptionLabel,
		},
		fmt.Sprintf("%s:%s", StrategyInitiativeCommunityEntity, StrategyInitiativeEntity): {
			SourceLabel:      StrategyInitiativeInitiativeCommunityAssociationSourceLabel,
			TargetLabel:      StrategyInitiativeInitiativeCommunityAssociationTargetLabel,
			DescriptionLabel: StrategyInitiativeInitiativeCommunityAssociationDescriptionLabel,
		},
	}

	strategyEntityTextMapping = map[StrategyEntityType]string{
		StrategyObjectiveEntity:           "strategy objectives",
		StrategyCapabilityCommunityEntity: "capability communities",
		StrategyInitiativeEntity:          "strategy initiatives",
		StrategyInitiativeCommunityEntity: "initiative communities",
	}
)

type StrategyEntityDialog struct {
	SourceLabel      string `json:"source_label"`
	TargetLabel      string `json:"target_label"`
	DescriptionLabel string `json:"description_label"`
}

type StrategyAssociation struct {
	Source           string `json:"source"` // hash
	Target           string `json:"target"` // range
	Advocate         string `json:"advocate"`
	PlatformID       string `json:"platform_id"`
	SourceTargetType string `json:"source_target_type"`
	Description      string `json:"description"`
	CreatedBy        string `json:"created_by"`
	CreatedAt        string `json:"created_at"`
	LastUpdatedBy    string `json:"last_updated_by"`
	LastUpdatedAt    string `json:"last_updated_at"`
}
