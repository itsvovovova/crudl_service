ALTER TABLE subscriptions
    ADD COLUMN created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    ADD COLUMN updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    ADD CONSTRAINT price_positive CHECK (price > 0);

CREATE INDEX subscriptions_user_id_id ON subscriptions (user_id, id);
