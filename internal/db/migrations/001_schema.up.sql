-- Food exclusion tags (14 EU allergens + vegetarian + vegan)
CREATE TABLE food_tags (
    id           SERIAL PRIMARY KEY,
    slug         TEXT UNIQUE NOT NULL,
    display_name TEXT NOT NULL
);

-- Users: single-row identity for both guests and registered accounts
CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email         TEXT UNIQUE,
    password_hash TEXT,
    guest_token   TEXT UNIQUE NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Per-user calorie preference
CREATE TABLE user_preferences (
    user_id        UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    calorie_target INT NOT NULL
);

-- Food tags a user has excluded
CREATE TABLE user_exclusions (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    tag_id  INT  NOT NULL REFERENCES food_tags(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, tag_id)
);

-- Master ingredient list with canonical unit
CREATE TABLE ingredients (
    id             SERIAL PRIMARY KEY,
    name           TEXT UNIQUE NOT NULL,
    canonical_unit TEXT NOT NULL
);

-- Allergen tags per ingredient
CREATE TABLE ingredient_tags (
    ingredient_id INT NOT NULL REFERENCES ingredients(id) ON DELETE CASCADE,
    tag_id        INT NOT NULL REFERENCES food_tags(id) ON DELETE CASCADE,
    PRIMARY KEY (ingredient_id, tag_id)
);

-- Recipe library
CREATE TABLE recipes (
    id                   SERIAL PRIMARY KEY,
    name                 TEXT NOT NULL,
    calories_per_serving INT  NOT NULL,
    description          TEXT NOT NULL DEFAULT '',
    cooking_instruction  TEXT NOT NULL DEFAULT '',
    cook_time_minutes    INT  NOT NULL DEFAULT 0
);

-- Per-recipe ingredient quantities (stored in canonical_unit)
CREATE TABLE recipe_ingredients (
    recipe_id     INT     NOT NULL REFERENCES recipes(id) ON DELETE CASCADE,
    ingredient_id INT     NOT NULL REFERENCES ingredients(id) ON DELETE RESTRICT,
    quantity      NUMERIC NOT NULL,
    PRIMARY KEY (recipe_id, ingredient_id)
);

-- A user's weekly meal plan
CREATE TABLE meal_plans (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- One meal per calendar date in a plan
CREATE TABLE plan_meals (
    id        UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plan_id   UUID NOT NULL REFERENCES meal_plans(id) ON DELETE CASCADE,
    meal_date DATE NOT NULL,
    recipe_id INT  NOT NULL REFERENCES recipes(id) ON DELETE RESTRICT,
    UNIQUE (plan_id, meal_date)
);

-- Pre-computed grocery list materialized at plan generation time
CREATE TABLE grocery_items (
    id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plan_id        UUID    NOT NULL REFERENCES meal_plans(id) ON DELETE CASCADE,
    ingredient_id  INT     NOT NULL REFERENCES ingredients(id) ON DELETE RESTRICT,
    total_quantity NUMERIC NOT NULL,
    unit           TEXT    NOT NULL,
    UNIQUE (plan_id, ingredient_id)
);

-- Access-pattern indexes
CREATE INDEX ON plan_meals(plan_id);
CREATE INDEX ON plan_meals(meal_date);
CREATE INDEX ON grocery_items(plan_id);
CREATE INDEX ON recipe_ingredients(recipe_id);
CREATE INDEX ON ingredient_tags(ingredient_id);
