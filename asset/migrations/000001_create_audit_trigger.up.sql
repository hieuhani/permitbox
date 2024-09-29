create or replace function auto_timestamps()
    returns trigger as $$
begin
    new.updated_at = now();
    return new;
end;
$$ language 'plpgsql';

create or replace function create_updated_at_trigger(table_name text) returns void as $$
begin
    execute 'create trigger ' || table_name || '_updated_at before update on ' || table_name || ' for each row execute procedure auto_timestamps()';
end;
$$ language plpgsql;