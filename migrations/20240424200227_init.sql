-- +goose Up
-- +goose StatementBegin
    CREATE TABLE IF NOT EXISTS public.gauges (
        id varchar NOT NULL,
        value double precision NOT NULL
    );
	ALTER TABLE public.gauges DROP CONSTRAINT IF EXISTS gauges_un;
    ALTER TABLE public.gauges ADD CONSTRAINT gauges_un UNIQUE (id);

	CREATE TABLE IF NOT EXISTS public.counters (id varchar NOT NULL,delta bigint NOT NULL);
	ALTER TABLE public.counters DROP CONSTRAINT IF EXISTS counters_un;
	ALTER TABLE public.counters ADD CONSTRAINT counters_un UNIQUE (id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
    DROP TABLE IF EXISTS public.gauges;
    DROP TABLE IF EXISTS public.counters;
-- +goose StatementEnd