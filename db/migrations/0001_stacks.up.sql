create table if not exists stacks (
                                      id uuid primary key,
                                      name text not null,
                                      slug text not null unique,
                                      created_at timestamptz not null default now(),
    archived_at timestamptz
    );