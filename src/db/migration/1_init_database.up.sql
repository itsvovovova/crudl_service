CREATE TABLE subscriptions (
    id SERIAL PRIMARY KEY,
    service_name VARCHAR(255) NOT NULL,
    price BIGINT NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NULL
);

CREATE INDEX subscriptions_user_id ON subscriptions (user_id);
CREATE INDEX subscriptions_service_name ON subscriptions (service_name);
CREATE INDEX subscriptions_price ON subscriptions (price);