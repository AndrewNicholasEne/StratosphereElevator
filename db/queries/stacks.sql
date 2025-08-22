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
where ($1::bool is false)
   or (archived_at is null)
order by created_at desc
limit $2::int offset $3::int;

-- name: ArchiveStack :one
update stacks
set archived_at = now()
where id = sqlc.arg(id)
  and archived_at is null
returning id, name, slug, created_at, archived_at;