-- +goose Up

-- Create Course table
create table course
(
    id          bigserial
        primary key,
    name        varchar(255) not null,
    description text,
    created_at  timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at  timestamp with time zone default CURRENT_TIMESTAMP,
    deleted_at  timestamp with time zone
);

create index idx_course_deleted_at
    on course (deleted_at);

create table chapter
(
    id          bigserial
        primary key,
    name        varchar(255) not null,
    description text,
    "order"     bigint       not null,
    course_id   bigint       not null
        constraint fk_course_chapters
            references course
            on delete cascade,
    created_at  timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at  timestamp with time zone default CURRENT_TIMESTAMP,
    deleted_at  timestamp with time zone
);


create index idx_chapter_deleted_at
    on chapter (deleted_at);

create table lesson
(
    id          bigserial
        primary key,
    name        varchar(255) not null,
    description text,
    content     text,
    "order"     bigint       not null,
    chapter_id  bigint       not null
        constraint fk_chapter_lessons
            references chapter
            on delete cascade,
    created_at  timestamp with time zone default CURRENT_TIMESTAMP,
    updated_at  timestamp with time zone default CURRENT_TIMESTAMP,
    deleted_at  timestamp with time zone
);


create index idx_lesson_deleted_at
    on lesson (deleted_at);

-- +goose Down
DROP TABLE IF EXISTS lesson;
DROP TABLE IF EXISTS chapter;
DROP TABLE IF EXISTS course;
