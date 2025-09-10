CREATE TABLE subscriptions (
    id SERIAL PRIMARY KEY,
    service_name VARCHAR(255) NOT NULL,
    price BIGINT NOT NULL,
    user_id VARCHAR(255) NOT NULL,
    start_date VARCHAR(10) NOT NULL,
    end_date VARCHAR(10) NULL
)