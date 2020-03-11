create table if not exists users
(
    id           int,
    user_id      varchar(10)  NOT NULL,
    platform_id  varchar(10)  NOT NULL,
    first_name   varchar(100) NULL,
    last_name    varchar(100) NULL,
    display_name varchar(100),
    timezone     varchar(20),
    `created_at` TIMESTAMP,
    `updated_at` TIMESTAMP,
    `deleted_at` TIMESTAMP,
    primary key (id)
);
create table if not exists vision
(
    id           int,
    vision_id    char(36)      NOT NULL,
    platform_id  varchar(10)   NOT NULL,
    vision       varchar(3000) NOT NULL,
    advocate     varchar(10)   NOT NULL,
    `created_at` TIMESTAMP,
    `updated_at` TIMESTAMP,
    `deleted_at` TIMESTAMP,
    primary key (id)
);
create table if not exists competencies
(
    id                     int,
    competency_id          char(36)      NOT NULL,
    competency_name        varchar(3000) NOT NULL,
    competency_description varchar(3000) NOT NULL,
    competency_type        varchar(20)   NOT NULL,
    platform_id            varchar(10)   NOT NULL,
    deactivated_on         varchar(30)   NOT NULL,
    `created_at`           TIMESTAMP,
    `updated_at`           TIMESTAMP,
    `deleted_at`           TIMESTAMP,
    primary key (id)
);
create table if not exists user_feedback
(
    id                int,
    feedback_id       char(36)      NOT NULL,
    source            varchar(10)   NOT NULL,
    target            varchar(10)   NOT NULL,
    competency_id     char(36)      NOT NULL,
    quarter           int,
    year              int,
    confidence_factor int,
    feedback          varchar(3000) NOT NULL,
    platform_id       varchar(10)   NOT NULL,
    `created_at`      TIMESTAMP,
    `updated_at`      TIMESTAMP,
    `deleted_at`      TIMESTAMP,
    primary key (id)
);
create table if not exists user_objectives
(
    id                             int,
    objective_id                   char(36)      NOT NULL,
    user_id                        varchar(10)   NOT NULL,
    name                           varchar(3000) NOT NULL,
    description                    varchar(3000) NOT NULL,
    accountability_partner         varchar(10)   NOT NULL,
    accepted                       int,
    type                           varchar(100)  NOT NULL,
    strategy_alignment_entity_id   char(36),
    strategy_alignment_entity_type varchar(100),
    quarter                        int           NOT NULL,
    year                           int           NOT NULL,
    created_date                   varchar(20)   NOT NULL,
    expected_end_date              varchar(20)   NOT NULL,
    completed                      int           NOT NULL,
    comments                       varchar(3000),
    canceled                       int,
    platform_id                    varchar(10)   NOT NULL,
    `created_at`                   TIMESTAMP,
    `updated_at`                   TIMESTAMP,
    `deleted_at`                   TIMESTAMP
);
create table if not exists user_objective_progress
(
    id                        int,
    objective_id              char(36)    NOT NULL,
    user_id                   varchar(10) NOT NULL,
    created_on                varchar(20) NOT NULL,
    comments                  varchar(3000),
    closeout                  int,
    percent_time_lapsed       varchar(10),
    status_color              varchar(20),
    reviewed_by_partner       bool,
    partner_comments          varchar(3000),
    partner_reported_progress varchar(3000),
    platform_id               varchar(10) NOT NULL,
    `created_at`              TIMESTAMP,
    `updated_at`              TIMESTAMP,
    `deleted_at`              TIMESTAMP
);