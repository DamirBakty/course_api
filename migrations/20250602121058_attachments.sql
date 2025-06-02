-- +goose Up

-- Create Attachment table
create table attachment
(
    id         bigserial
        primary key,
    name       varchar(255) not null,
    url        varchar(255) not null,
    lesson_id  bigint       not null
        constraint fk_lesson_attachments
            references lesson
            on delete cascade,
    created_at timestamp with time zone default CURRENT_TIMESTAMP,
    deleted_at timestamp with time zone
);

create index idx_attachment_deleted_at
    on attachment (deleted_at);

-- +goose Down
DROP TABLE IF EXISTS attachment;