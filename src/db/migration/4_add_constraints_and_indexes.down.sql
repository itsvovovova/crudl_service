DROP INDEX IF EXISTS subscriptions_user_id_id;

ALTER TABLE subscriptions
    DROP CONSTRAINT IF EXISTS price_positive,
    DROP COLUMN IF EXISTS updated_at,
    DROP COLUMN IF EXISTS created_at;
