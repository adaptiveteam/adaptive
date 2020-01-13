package queries

const StrategyStatus = `
SELECT
objective.id AS 'objective_id',
initiative_updates.id AS 'initiative_id',
objective.name AS 'Objective Name',
objective.type AS 'Objective Type',
DATE(objective.created_at) as 'Objective Created On',
DATE(objective_updates.updated_at) as 'Objective Updated On',
objective.expected_end_date as 'Objective End Date',
CONCAT(ROUND(DATEDIFF(objective.expected_end_date, CURDATE())/DATEDIFF(objective.expected_end_date, objective.created_at)*100,0),"%") as 'Objective Time Left',
objective_advocates.display_name as 'Objective Advocate',
CASE
        WHEN objective_updates.status_color is null THEN 'No Status'
        WHEN objective_updates.status_color = 'Red' THEN 'Off Track'
        WHEN objective_updates.status_color = 'Yellow' THEN 'At Risk'
        WHEN objective_updates.status_color = 'Green' THEN 'On Track'
END AS 'Objective Status',
CASE
        WHEN objective_updates.comments is null THEN 'No Updates'
        ELSE objective_updates.comments
END AS 'Objective Update',
objective.description as 'Objective Description',
CASE
        WHEN initiative_updates.name is null THEN 'No Initiatives'
        ELSE initiative_updates.name
END AS 'Initiative Name',
CASE
        WHEN initiative_advocates.display_name is null THEN 'No Initiatives'
        ELSE initiative_advocates.display_name
END AS 'Initiative Advocate',
CASE
        WHEN initiative_updates.status_color is null THEN 'No Status'
        WHEN initiative_updates.status_color = 'Red' THEN 'Off Track'
        WHEN initiative_updates.status_color = 'Yellow' THEN 'At Risk'
        WHEN initiative_updates.status_color = 'Green' THEN 'On Track'
END AS 'Initiative Status',
CASE
        WHEN initiative_updates.created_at is null THEN 'No Initiatives'
        ELSE DATE(initiative_updates.created_at)
END as 'Initiative Created On',
CASE
        WHEN initiative_updates.updated_at is null THEN 'No Initiatives'
        ELSE DATE(initiative_updates.last_updated)
END as 'Initiative Updated On',
CASE
        WHEN initiative_updates.expected_end_date is null THEN 'No Initiatives'
        ELSE DATE(initiative_updates.expected_end_date)
END as 'Initiative End Date',
CASE
        WHEN initiative_updates.expected_end_date is null THEN 'No Initiatives'
        ELSE CONCAT(ROUND(DATEDIFF(initiative_updates.expected_end_date, CURDATE())/DATEDIFF(initiative_updates.expected_end_date, initiative_updates.created_at)*100,0),"%")
        
END as 'Initiative Time Left',
CASE
        WHEN initiative_updates.comments is null THEN 'No Updates'
        ELSE initiative_updates.comments
END AS 'Initiative Update',
CASE
        WHEN initiative_updates.description is null THEN 'No Initiatives'
        ELSE initiative_updates.description
END AS 'Initiative Description'

FROM

user as objective_advocates,
objective

LEFT JOIN (
        SELECT

        user_objective_progress.created_date as updated_at,
        user_objective_progress.objective_id,
        user_objective_progress.status_color,
        user_objective_progress.comments

        FROM

        user_objective_progress,
        objective

        WHERE

        user_objective_progress.objective_id = objective.id AND
        user_objective_progress.platform_id = objective.platform_id

) AS objective_updates ON
objective_updates.objective_id = objective.id

LEFT JOIN (
        SELECT
                initiative.id,
                initiative.name,
                initiative.description,
                initiative.advocate,
                initiative.created_at,
                initiative.updated_at,
                initiative.expected_end_date,
                initiative.capability_objective_id as objective_id,
                latest_update.status_color,
                latest_update.comments,
                MAX(latest_update.created_date) as last_updated

        FROM

        initiative

        LEFT JOIN (
                SELECT
                        user_objective_progress.objective_id,
                        user_objective_progress.platform_id,
                        user_objective_progress.status_color,
                        user_objective_progress.comments,
                        user_objective_progress.partner_comments,
                        user_objective_progress.created_date
                FROM
                        user_objective_progress
        ) AS latest_update ON
        initiative.id = latest_update.objective_id AND
        initiative.platform_id = latest_update.platform_id
        GROUP BY
                initiative.name
        ORDER BY
                initiative.name
) AS initiative_updates ON
objective.id = initiative_updates.objective_id

LEFT JOIN user as initiative_advocates ON
initiative_advocates.id = initiative_updates.advocate

WHERE

objective.platform_id = ? AND
objective_advocates.id = objective.advocate

ORDER BY

objective.type,
objective.name
`
