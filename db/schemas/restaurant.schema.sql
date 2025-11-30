-- =========================
-- RESTAURANTS
-- =========================
CREATE TABLE IF NOT EXISTS restaurant (
  id             INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  name           TEXT NOT NULL,
  description    TEXT,
  address        TEXT,
  category       TEXT,
  city           TEXT,
  district       TEXT,
  logo_url       TEXT,
  banner_url     TEXT,
  phone_number   TEXT,
  website_url    TEXT,
  email          CITEXT,
  created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  user_id        UUID
);

CREATE TRIGGER trg_restaurant_updated_at
BEFORE UPDATE ON restaurant
FOR EACH ROW EXECUTE FUNCTION set_updated_at();


CREATE TABLE restaurant_hours (
  restaurant_id INT REFERENCES restaurant(id) ON DELETE CASCADE,
  day_of_week   INT NOT NULL CHECK (day_of_week BETWEEN 0 AND 6), -- 0=Sun
  open_time     TIME,
  close_time    TIME,
  is_closed     BOOLEAN NOT NULL DEFAULT FALSE,
  PRIMARY KEY (restaurant_id, day_of_week)
);
