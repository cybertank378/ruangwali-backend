CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- =========================================================
-- USERS
-- =========================================================

CREATE TABLE users (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

                       email VARCHAR(255) NOT NULL,
                       password_hash TEXT NOT NULL,

                       status VARCHAR(30) NOT NULL DEFAULT 'ACTIVE',

                       last_login_at TIMESTAMPTZ,

                       created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                       updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

                       CONSTRAINT chk_users_email_not_blank
                           CHECK (BTRIM(email) <> ''),

                       CONSTRAINT chk_users_password_hash_not_blank
                           CHECK (BTRIM(password_hash) <> ''),

                       CONSTRAINT chk_users_status
                           CHECK (
                               status IN (
                                          'ACTIVE',
                                          'INACTIVE',
                                          'SUSPENDED'
                                   )
                               )
);

CREATE UNIQUE INDEX uq_users_email_normalized
    ON users (LOWER(BTRIM(email)));

CREATE INDEX idx_users_status
    ON users (status);

-- =========================================================
-- ROLES
-- =========================================================

CREATE TABLE roles (
                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

                       code VARCHAR(100) NOT NULL,
                       name VARCHAR(150) NOT NULL,
                       description TEXT,

                       is_system BOOLEAN NOT NULL DEFAULT FALSE,
                       is_active BOOLEAN NOT NULL DEFAULT TRUE,

                       created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                       updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

                       CONSTRAINT uq_roles_code
                           UNIQUE (code),

                       CONSTRAINT chk_roles_code_not_blank
                           CHECK (BTRIM(code) <> ''),

                       CONSTRAINT chk_roles_name_not_blank
                           CHECK (BTRIM(name) <> '')
);

CREATE INDEX idx_roles_is_active
    ON roles (is_active);

-- =========================================================
-- PERMISSIONS
-- =========================================================

CREATE TABLE permissions (
                             id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

                             code VARCHAR(150) NOT NULL,
                             resource VARCHAR(100) NOT NULL,
                             action VARCHAR(100) NOT NULL,

                             description TEXT,

                             created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                             updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

                             CONSTRAINT uq_permissions_code
                                 UNIQUE (code),

                             CONSTRAINT uq_permissions_resource_action
                                 UNIQUE (
                                         resource,
                                         action
                                     ),

                             CONSTRAINT chk_permissions_code_not_blank
                                 CHECK (BTRIM(code) <> ''),

                             CONSTRAINT chk_permissions_resource_not_blank
                                 CHECK (BTRIM(resource) <> ''),

                             CONSTRAINT chk_permissions_action_not_blank
                                 CHECK (BTRIM(action) <> '')
);

CREATE INDEX idx_permissions_resource
    ON permissions (resource);

-- =========================================================
-- USER ROLES
-- =========================================================

CREATE TABLE user_roles (
                            user_id UUID NOT NULL,
                            role_id UUID NOT NULL,

                            assigned_by UUID,
                            assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

                            PRIMARY KEY (
                                         user_id,
                                         role_id
                                ),

                            CONSTRAINT fk_user_roles_user
                                FOREIGN KEY (user_id)
                                    REFERENCES users(id)
                                    ON DELETE CASCADE,

                            CONSTRAINT fk_user_roles_role
                                FOREIGN KEY (role_id)
                                    REFERENCES roles(id)
                                    ON DELETE CASCADE,

                            CONSTRAINT fk_user_roles_assigned_by
                                FOREIGN KEY (assigned_by)
                                    REFERENCES users(id)
                                    ON DELETE SET NULL
);

CREATE INDEX idx_user_roles_role_id
    ON user_roles (role_id);

CREATE INDEX idx_user_roles_assigned_by
    ON user_roles (assigned_by);

-- =========================================================
-- ROLE PERMISSIONS
-- =========================================================

CREATE TABLE role_permissions (
                                  role_id UUID NOT NULL,
                                  permission_id UUID NOT NULL,

                                  granted_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

                                  PRIMARY KEY (
                                               role_id,
                                               permission_id
                                      ),

                                  CONSTRAINT fk_role_permissions_role
                                      FOREIGN KEY (role_id)
                                          REFERENCES roles(id)
                                          ON DELETE CASCADE,

                                  CONSTRAINT fk_role_permissions_permission
                                      FOREIGN KEY (permission_id)
                                          REFERENCES permissions(id)
                                          ON DELETE CASCADE
);

CREATE INDEX idx_role_permissions_permission_id
    ON role_permissions (permission_id);