CREATE TABLE subscriptions (
    service_name VARCHAR(255) NOT NULL,
    price NOT NULL,
    user_id UUID NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NULL
)