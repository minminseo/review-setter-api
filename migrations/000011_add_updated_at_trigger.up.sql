-- updated_atカラムに対して、updated_atを自動更新するトリガーの追加
CREATE OR REPLACE FUNCTION set_updated_at()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at := CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- トリガーを追加
DO $$
BEGIN
    -- email_verifications
    DROP TRIGGER IF EXISTS trigger_set_updated_at ON email_verifications;
    CREATE TRIGGER trigger_set_updated_at
        BEFORE UPDATE ON email_verifications
        FOR EACH ROW EXECUTE FUNCTION set_updated_at();
END;
$$;