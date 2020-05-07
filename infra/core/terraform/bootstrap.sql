create table if not exists users
(
    id           int,
    user_id      TEXT  NOT NULL,
    platform_id  TEXT  NOT NULL,
    first_name   TEXT NULL,
    last_name    TEXT NULL,
    display_name TEXT,
    timezone     TEXT,
    `created_at` TIMESTAMP,
    `updated_at` TIMESTAMP,
    `deleted_at` TIMESTAMP,
    primary key (id)
);
create table if not exists vision
(
    id           int,
    vision_id    TEXT      NOT NULL,
    platform_id  TEXT   NOT NULL,
    vision       TEXT NOT NULL,
    advocate     TEXT   NOT NULL,
    `created_at` TIMESTAMP,
    `updated_at` TIMESTAMP,
    `deleted_at` TIMESTAMP,
    primary key (id)
);
create table if not exists competencies
(
    id                     int,
    competency_id          TEXT      NOT NULL,
    competency_name        TEXT NOT NULL,
    competency_description TEXT NOT NULL,
    competency_type        TEXT   NOT NULL,
    platform_id            TEXT   NOT NULL,
    deactivated_on         TEXT   NOT NULL,
    `created_at`           TIMESTAMP,
    `updated_at`           TIMESTAMP,
    `deleted_at`           TIMESTAMP,
    primary key (id)
);
create table if not exists user_feedback
(
    id                int,
    feedback_id       TEXT      NOT NULL,
    source            TEXT   NOT NULL,
    target            TEXT   NOT NULL,
    competency_id     TEXT      NOT NULL,
    quarter           int,
    year              int,
    confidence_factor int,
    feedback          TEXT NOT NULL,
    platform_id       TEXT   NOT NULL,
    `created_at`      TIMESTAMP,
    `updated_at`      TIMESTAMP,
    `deleted_at`      TIMESTAMP,
    primary key (id)
);
create table if not exists user_objectives
(
    id                             int,
    objective_id                   TEXT      NOT NULL,
    user_id                        TEXT   NOT NULL,
    name                           TEXT NOT NULL,
    description                    TEXT NOT NULL,
    accountability_partner         TEXT   NOT NULL,
    accepted                       int,
    type                           TEXT  NOT NULL,
    strategy_alignment_entity_id   TEXT,
    strategy_alignment_entity_type TEXT,
    quarter                        int           NOT NULL,
    year                           int           NOT NULL,
    created_date                   TEXT   NOT NULL,
    expected_end_date              TEXT   NOT NULL,
    completed                      int           NOT NULL,
    comments                       TEXT,
    canceled                       int,
    platform_id                    TEXT   NOT NULL,
    `created_at`                   TIMESTAMP,
    `updated_at`                   TIMESTAMP,
    `deleted_at`                   TIMESTAMP
);
create table if not exists user_objective_progress
(
    id                        int,
    objective_id              TEXT    NOT NULL,
    user_id                   TEXT NOT NULL,
    created_on                TEXT NOT NULL,
    comments                  TEXT,
    closeout                  int,
    percent_time_lapsed       TEXT,
    status_color              TEXT,
    reviewed_by_partner       bool,
    partner_comments          TEXT,
    partner_reported_progress TEXT,
    platform_id               TEXT NOT NULL,
    `created_at`              TIMESTAMP,
    `updated_at`              TIMESTAMP,
    `deleted_at`              TIMESTAMP
);