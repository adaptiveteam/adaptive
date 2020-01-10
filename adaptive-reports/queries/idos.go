package queries

const IDOs = `
SELECT

user_objective.completed,
user_objective.created_date as ido_created_at,
DATE(user_objective_progress.created_date) as updated_at,
advocate.display_name as advocate,
coach.display_name as coach,
user_objective.name as ido_name,
user_objective.description as ido_description,
CASE
        WHEN user_objective_progress.status_color = "" THEN "No Status"
        WHEN user_objective_progress.status_color = "Red" THEN "Off Track"
        WHEN user_objective_progress.status_color = "Yellow" THEN "At Risk"
        WHEN user_objective_progress.status_color = "Green" THEN "On Track"
        ELSE user_objective_progress.status_color
END as advocate_status,

CASE
        WHEN user_objective_progress.partner_reported_progress = "" THEN "No Status"
        WHEN user_objective_progress.partner_reported_progress = "Red" THEN "Off Track"
        WHEN user_objective_progress.partner_reported_progress = "Yellow" THEN "At Risk"
        WHEN user_objective_progress.partner_reported_progress = "Green" THEN "On Track"
        ELSE user_objective_progress.partner_reported_progress
END as coach_status,

CASE
        WHEN user_objective_progress.comments = "" THEN "No comments"
        ELSE user_objective_progress.comments
END as advocate_comments,
CASE
        WHEN user_objective_progress.partner_comments = "" THEN "No comments"
        ELSE user_objective_progress.partner_comments
END as coach_comments,

CASE
        WHEN objective.name is not null THEN objective.name
        WHEN initiative.name is not null THEN initiative.name
        WHEN competency.name is not null THEN competency.name
        ELSE 'no alignment'
END AS "focused_on_name",

CASE
        WHEN objective.description is not null THEN objective.description
        WHEN initiative.description is not null THEN initiative.description
        WHEN competency.description is not null THEN competency.description
        ELSE 'no alignment'
END AS "focused_on_description",

CASE
        WHEN user_objective.strategy_alignment_entity_type = "strategy_objective" THEN "Objective"
        WHEN user_objective.strategy_alignment_entity_type = "strategy_initiative" THEN "Initiative"
        WHEN user_objective.strategy_alignment_entity_type = "competency" THEN "Competency"
END AS "is_a",
CASE
        WHEN objective.name is not null THEN "Our vision"
        WHEN initiative.name is not null THEN aligned_objective.name
        WHEN competency.name is not null THEN "Team Strength"
        ELSE 'no alignment'
END AS "driving"


FROM

user as coach,
user as advocate,
user_objective

LEFT JOIN user_objective_progress ON
user_objective_progress.objective_id = user_objective.id

LEFT JOIN objective ON
user_objective.strategy_alignment_entity_id = objective.id

LEFT JOIN initiative ON
user_objective.strategy_alignment_entity_id = initiative.id

LEFT JOIN competency ON
user_objective.strategy_alignment_entity_id = competency.id

LEFT JOIN objective AS aligned_objective ON
initiative.capability_objective_id = aligned_objective.id

WHERE

user_objective.user_id = ? AND
user_objective.user_id = advocate.id aND
user_objective.accountability_partner = coach.id

ORDER BY

completed,
ido_name,
updated_at DESC
`
