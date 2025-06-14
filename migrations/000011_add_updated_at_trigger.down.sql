DO $$
BEGIN
    DROP TRIGGER IF EXISTS trigger_set_updated_at ON email_verifications;
END;
$$;

DROP FUNCTION IF EXISTS set_updated_at();
