CREATE INDEX IF NOT EXISTS idx_subs_id 
ON subscriptions USING HASH (id);

CREATE INDEX IF NOT EXISTS idx_subs_user_id 
ON subscriptions USING HASH (user_id);

CREATE INDEX IF NOT EXISTS idx_subs_service_name 
ON subscriptions USING HASH (service_name);

CREATE INDEX IF NOT EXISTS idx_subs_dates 
ON subscriptions (start_date, end_date);