-- name: GetUser :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users 
WHERE id = $1;

-- name: GetPlanByName :one
SELECT * FROM plans 
WHERE name = $1;

-- name: RegisterUser :exec
INSERT INTO users (name, email, password_hash, email_verified_at)
VALUES ($1, $2, $3, $4);

-- name: CreatePlan :exec
INSERT INTO plans (name, price_monthly, max_websites, check_interval_seconds, has_performance_reports, has_seo_audits, has_public_status_page )
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: CreateSubscription :exec
INSERT INTO subscriptions (user_id ,plan_id, status, stripe_subscription_id, current_period_ends_at )
VALUES ($1 , $2, $3, $4, $5);
