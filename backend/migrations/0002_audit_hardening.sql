-- The API role may append and read audit events, but cannot alter history.
DO $$ BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'ssp_audit_owner') THEN
        CREATE ROLE ssp_audit_owner NOLOGIN;
    END IF;
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'ssp_app') THEN
        CREATE ROLE ssp_app NOLOGIN;
    END IF;
END $$;

GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO ssp_app;
GRANT USAGE, SELECT, UPDATE ON ALL SEQUENCES IN SCHEMA public TO ssp_app;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT ALL ON TABLES TO ssp_app;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT USAGE, SELECT, UPDATE ON SEQUENCES TO ssp_app;

-- Managed Postgres providers require explicit membership before ownership transfer.
GRANT ssp_audit_owner TO CURRENT_USER;
GRANT CREATE ON SCHEMA public TO ssp_audit_owner;
ALTER TABLE audit_events OWNER TO ssp_audit_owner;
REVOKE UPDATE, DELETE ON audit_events FROM PUBLIC, ssp_app;
GRANT SELECT, INSERT ON audit_events TO ssp_app;

CREATE OR REPLACE FUNCTION prune_audit_events(cutoff timestamptz)
RETURNS bigint
LANGUAGE plpgsql
SECURITY DEFINER
SET search_path = pg_catalog, public
AS $$
DECLARE removed bigint;
BEGIN
    DELETE FROM public.audit_events WHERE occurred_at < cutoff;
    GET DIAGNOSTICS removed = ROW_COUNT;
    RETURN removed;
END;
$$;
ALTER FUNCTION prune_audit_events(timestamptz) OWNER TO ssp_audit_owner;
REVOKE CREATE ON SCHEMA public FROM ssp_audit_owner;
REVOKE ssp_audit_owner FROM CURRENT_USER;
REVOKE ALL ON FUNCTION prune_audit_events(timestamptz) FROM PUBLIC;
GRANT EXECUTE ON FUNCTION prune_audit_events(timestamptz) TO ssp_app;
