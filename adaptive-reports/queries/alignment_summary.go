package queries

const AlignmentSummary = `
SELECT
        user_objective.id,
CASE
        WHEN objective.name is not null THEN objective.name
        WHEN initiative.name is not null THEN initiative.name
        WHEN competency.name is not null THEN competency.name
        ELSE 'no alignment'
END AS "FocusedOn",
CASE
        WHEN objective.id is not null THEN objective.id
        WHEN initiative.id is not null THEN initiative.id
        WHEN competency.id is not null THEN competency.id
        ELSE 'no alignment'
END AS "FocusedOnID",
CASE
        WHEN user_objective.strategy_alignment_entity_type = "strategy_objective" THEN "Objective"
        WHEN user_objective.strategy_alignment_entity_type = "strategy_initiative" THEN "Initiative"
        WHEN user_objective.strategy_alignment_entity_type = "competency" THEN "Competency"
END AS "IsA",
CASE
        WHEN objective.name is not null THEN "Our vision"
        WHEN initiative.name is not null THEN aligned_objective.name
        WHEN competency.name is not null THEN "Team Strength"
        ELSE 'no alignment'
END AS "Driving",
CASE
        WHEN initiative.id is not null THEN aligned_objective.id
        ELSE null
END AS "DrivingID",
CASE
        WHEN updated_recently.latest_update is not null AND DATE_ADD(CURDATE(), INTERVAL -7 DAY) <= updated_recently.latest_update THEN "Yes"
        ELSE "No"
END as "Updated",
CASE
        WHEN updated_recently.status is not null THEN updated_recently.status
        ELSE "No Status"
END as "AdvocateStatus",
CASE
        WHEN updated_recently.coach_status is not null THEN updated_recently.coach_status
        ELSE "No Status"
END as "CoachStatus",
team_member.display_name as "Advocate",
coach.display_name as "Coach",
user_objective.name "IDOName",
user_objective.created_date as "CreatedOn",
DATE(updated_recently.latest_update) as "UpdatedOn",
user_objective.expected_end_date as "CompleteBy"

FROM

user AS team_member,
user AS coach,
user_objective

LEFT JOIN objective ON
user_objective.strategy_alignment_entity_id = objective.id

LEFT JOIN initiative ON
user_objective.strategy_alignment_entity_id = initiative.id

LEFT JOIN competency ON
user_objective.strategy_alignment_entity_id = competency.id

LEFT JOIN objective AS aligned_objective ON
initiative.capability_objective_id = aligned_objective.id

LEFT JOIN (
        SELECT
                user_id,
                objective_id,
                MAX(created_date) as latest_update,
                updated_at,
                CASE
                        WHEN user_objective_progress.partner_reported_progress is null THEN 'No Status'
                        WHEN user_objective_progress.partner_reported_progress = 'Red' THEN 'Off Track'
                        WHEN user_objective_progress.partner_reported_progress = 'Yellow' THEN 'At Risk'
                        WHEN user_objective_progress.partner_reported_progress = 'Green' THEN 'On Track'
                END AS 'coach_status',
                CASE
                        WHEN user_objective_progress.status_color is null THEN 'No Status'
                        WHEN user_objective_progress.status_color = 'Red' THEN 'Off Track'
                        WHEN user_objective_progress.status_color = 'Yellow' THEN 'At Risk'
                        WHEN user_objective_progress.status_color = 'Green' THEN 'On Track'
                END AS 'status'
        FROM
                user_objective_progress
        GROUP BY objective_id
) AS updated_recently ON
user_objective.id = updated_recently.objective_id

where

user_objective.type = "individual" and
user_objective.platform_id = ? and
user_objective.user_id = team_member.id and
user_objective.accountability_partner = coach.id
`
