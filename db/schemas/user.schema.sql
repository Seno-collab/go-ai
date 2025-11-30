CREATE EXTENSION IF NOT EXISTS pgcrypto;  -- để dùng gen_random_uuid()
CREATE EXTENSION IF NOT EXISTS citext;    -- email không phân biệt hoa/thường

CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
  NEW.updated_at := NOW();
  RETURN NEW;
END; $$;

-- =========================
-- ROLES
-- =========================
CREATE TABLE IF NOT EXISTS role (
  id           INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  role_name    TEXT NOT NULL UNIQUE,
  created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER trg_role_updated_at
BEFORE UPDATE ON role
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

-- =========================
-- USERS (chú ý: "user" là từ khóa; nếu giữ tên này, luôn để trong dấu ")
-- =========================
CREATE TABLE IF NOT EXISTS "user" (
  id             UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  full_name      TEXT NOT NULL,
  email          CITEXT UNIQUE,
  password_hash  TEXT NOT NULL,
  role_id        INT REFERENCES role(id) ON UPDATE CASCADE ON DELETE SET NULL,
  is_active      BOOLEAN NOT NULL DEFAULT TRUE,
  created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TRIGGER trg_user_updated_at
BEFORE UPDATE ON "user"
FOR EACH ROW EXECUTE FUNCTION set_updated_at();
