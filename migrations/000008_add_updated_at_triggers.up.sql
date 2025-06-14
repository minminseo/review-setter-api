-- 全てのテーブルのupdated_atカラムに対して、updated_atを自動更新するトリガーの追加
CREATE OR REPLACE FUNCTION set_updated_at()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at := CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 各テーブルにトリガーを追加
DO $$
BEGIN
    -- users
    DROP TRIGGER IF EXISTS trigger_set_updated_at ON users;
    CREATE TRIGGER trigger_set_updated_at
        BEFORE UPDATE ON users
        FOR EACH ROW EXECUTE FUNCTION set_updated_at();

    -- categories
    DROP TRIGGER IF EXISTS trigger_set_updated_at ON categories;
    CREATE TRIGGER trigger_set_updated_at
        BEFORE UPDATE ON categories
        FOR EACH ROW EXECUTE FUNCTION set_updated_at();

    -- review_patterns
    DROP TRIGGER IF EXISTS trigger_set_updated_at ON review_patterns;
    CREATE TRIGGER trigger_set_updated_at
        BEFORE UPDATE ON review_patterns
        FOR EACH ROW EXECUTE FUNCTION set_updated_at();

    -- review_boxes
    DROP TRIGGER IF EXISTS trigger_set_updated_at ON review_boxes;
    CREATE TRIGGER trigger_set_updated_at
        BEFORE UPDATE ON review_boxes
        FOR EACH ROW EXECUTE FUNCTION set_updated_at();

    -- review_items
    DROP TRIGGER IF EXISTS trigger_set_updated_at ON review_items;
    CREATE TRIGGER trigger_set_updated_at
        BEFORE UPDATE ON review_items
        FOR EACH ROW EXECUTE FUNCTION set_updated_at();

    -- pattern_steps
    DROP TRIGGER IF EXISTS trigger_set_updated_at ON pattern_steps;
    CREATE TRIGGER trigger_set_updated_at
        BEFORE UPDATE ON pattern_steps
        FOR EACH ROW EXECUTE FUNCTION set_updated_at();

    -- review_dates
    DROP TRIGGER IF EXISTS trigger_set_updated_at ON review_dates;
    CREATE TRIGGER trigger_set_updated_at
        BEFORE UPDATE ON review_dates
        FOR EACH ROW EXECUTE FUNCTION set_updated_at();

    -- email_verifications
    DROP TRIGGER IF EXISTS trigger_set_updated_at ON email_verifications;
    CREATE TRIGGER trigger_set_updated_at
        BEFORE UPDATE ON email_verifications
        FOR EACH ROW EXECUTE FUNCTION set_updated_at();
END;
$$;