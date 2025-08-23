-- name: CreateStack :one
insert into stacks (id, name, slug)
values (sqlc.arg(id), sqlc.arg(name), sqlc.arg(slug))
returning id, name, slug, created_at, archived_at;

-- name: GetStackBySlug :one
select id, name, slug, created_at, archived_at
from stacks
where slug = sqlc.arg(slug);

-- name: ListStacks :many
select id, name, slug, created_at, archived_at
from stacks
where ($1::bool is true)
   or (archived_at is null)
order by created_at desc
limit $2::int offset $3::int;

-- name: ArchiveStack :execrows
UPDATE stacks
SET archived_at = now()
WHERE id = sqlc.arg(id)
  AND archived_at IS NULL;

-- name: StackArchivedStatus :one
SELECT EXISTS (SELECT 1
               FROM stacks
               WHERE id = sqlc.arg(id)
                 AND archived_at IS NOT NULL) AS already_archived;