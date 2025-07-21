-- name: SignInUser :one 
SELECT * FROM users
WHERE email = $1 AND password_hash = $2;

-- name: GetUser :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users 
WHERE id = $1;

-- name: GetPlanByName :one
SELECT * FROM plans 
WHERE name = $1;

-- name: GetSubscriptionByStripeSbId :one
SELECT * FROM subscriptions 
WHERE stripe_subscription_id = $1;

-- name: RegisterUser :exec
INSERT INTO users (name, email, password_hash, email_verified_at)
VALUES ($1, $2, $3, $4);

-- name: CreatePlan :exec
INSERT INTO plans (name, price_monthly, max_websites, check_interval_seconds, has_performance_reports, has_seo_audits, has_public_status_page )
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: CreateSubscription :exec
INSERT INTO subscriptions (user_id ,plan_id, status, stripe_subscription_id, current_period_ends_at )
VALUES ($1 , $2, $3, $4, $5);

-- name: GetAlertIncidents :many
(
    SELECT 
        w.id,
        w.name,
        w.url,
        'Uptime' AS alert_type, 
        i.cause,
        i.started_at,
        NULL AS lcp_ms, 
        NULL AS performance_checked_at
    FROM 
        websites w
    JOIN 
        incidents i ON w.id = i.website_id
    WHERE 
        w.is_active = TRUE 
        AND i.ended_at IS NULL 
)
UNION ALL
(
    WITH LatestPerformanceReport AS (
        SELECT
            website_id,
            lcp_ms,
            checked_at,
            ROW_NUMBER() OVER(PARTITION BY website_id ORDER BY checked_at DESC) as rn
        FROM 
            performance_reports
    )
    SELECT 
        w.id,
        w.name,
        w.url,
        'Performance' AS alert_type,
        'LCP acima de 2500ms' AS cause, 
        NULL AS started_at,
        pr.lcp_ms,
        pr.checked_at AS performance_checked_at
    FROM 
        websites w
    JOIN 
        LatestPerformanceReport pr ON w.id = pr.website_id
    WHERE 
        w.is_active = TRUE
        AND pr.rn = 1 
        AND pr.lcp_ms > 2500 
        AND w.id NOT IN (SELECT website_id FROM incidents WHERE ended_at IS NULL)
)
ORDER BY 
    COALESCE(started_at, performance_checked_at) DESC;
