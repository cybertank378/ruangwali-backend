-- =========================================================
-- TEACHERS
-- =========================================================

CREATE TABLE teachers (
                          id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

                          user_id UUID UNIQUE,

                          employee_number VARCHAR(100),
                          nip VARCHAR(50),

                          full_name VARCHAR(255) NOT NULL,

                          gender VARCHAR(20),
                          phone VARCHAR(50),

                          status VARCHAR(30) NOT NULL DEFAULT 'ACTIVE',

                          created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                          updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

                          CONSTRAINT fk_teachers_user
                              FOREIGN KEY (user_id)
                                  REFERENCES users(id)
                                  ON DELETE SET NULL,

                          CONSTRAINT chk_teachers_full_name_not_blank
                              CHECK (BTRIM(full_name) <> ''),

                          CONSTRAINT chk_teachers_gender
                              CHECK (
                                  gender IS NULL
                                      OR gender IN (
                                                    'MALE',
                                                    'FEMALE'
                                      )
                                  ),

                          CONSTRAINT chk_teachers_status
                              CHECK (
                                  status IN (
                                             'ACTIVE',
                                             'INACTIVE'
                                      )
                                  )
);

CREATE UNIQUE INDEX uq_teachers_employee_number
    ON teachers (employee_number)
    WHERE employee_number IS NOT NULL;

CREATE UNIQUE INDEX uq_teachers_nip
    ON teachers (nip)
    WHERE nip IS NOT NULL;

CREATE INDEX idx_teachers_full_name
    ON teachers (full_name);

CREATE INDEX idx_teachers_status
    ON teachers (status);

-- =========================================================
-- ACADEMIC YEARS
-- =========================================================

CREATE TABLE academic_years (
                                id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

                                name VARCHAR(50) NOT NULL,

                                start_date DATE NOT NULL,
                                end_date DATE NOT NULL,

                                is_active BOOLEAN NOT NULL DEFAULT FALSE,

                                created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

                                CONSTRAINT uq_academic_years_name
                                    UNIQUE (name),

                                CONSTRAINT chk_academic_years_name_not_blank
                                    CHECK (BTRIM(name) <> ''),

                                CONSTRAINT chk_academic_years_date_range
                                    CHECK (end_date > start_date)
);

CREATE UNIQUE INDEX uq_academic_years_single_active
    ON academic_years ((TRUE))
    WHERE is_active = TRUE;

-- =========================================================
-- SEMESTERS
-- =========================================================

CREATE TABLE semesters (
                           id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

                           academic_year_id UUID NOT NULL,

                           code VARCHAR(30) NOT NULL,
                           name VARCHAR(100) NOT NULL,

                           start_date DATE NOT NULL,
                           end_date DATE NOT NULL,

                           is_active BOOLEAN NOT NULL DEFAULT FALSE,

                           created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                           updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

                           CONSTRAINT fk_semesters_academic_year
                               FOREIGN KEY (academic_year_id)
                                   REFERENCES academic_years(id)
                                   ON DELETE CASCADE,

                           CONSTRAINT uq_semesters_academic_year_code
                               UNIQUE (
                                       academic_year_id,
                                       code
                                   ),

                           CONSTRAINT chk_semesters_code_not_blank
                               CHECK (BTRIM(code) <> ''),

                           CONSTRAINT chk_semesters_name_not_blank
                               CHECK (BTRIM(name) <> ''),

                           CONSTRAINT chk_semesters_date_range
                               CHECK (end_date > start_date)
);

CREATE INDEX idx_semesters_academic_year_id
    ON semesters (academic_year_id);

CREATE UNIQUE INDEX uq_semesters_single_active
    ON semesters ((TRUE))
    WHERE is_active = TRUE;

-- =========================================================
-- SUBJECTS
-- =========================================================

CREATE TABLE subjects (
                          id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

                          code VARCHAR(50) NOT NULL,
                          name VARCHAR(150) NOT NULL,

                          description TEXT,

                          is_active BOOLEAN NOT NULL DEFAULT TRUE,

                          created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                          updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

                          CONSTRAINT uq_subjects_code
                              UNIQUE (code),

                          CONSTRAINT chk_subjects_code_not_blank
                              CHECK (BTRIM(code) <> ''),

                          CONSTRAINT chk_subjects_name_not_blank
                              CHECK (BTRIM(name) <> '')
);

CREATE INDEX idx_subjects_name
    ON subjects (name);

CREATE INDEX idx_subjects_is_active
    ON subjects (is_active);

-- =========================================================
-- CLASSROOMS
-- =========================================================

CREATE TABLE classrooms (
                            id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

                            academic_year_id UUID NOT NULL,

                            code VARCHAR(50) NOT NULL,
                            name VARCHAR(150) NOT NULL,

                            grade_level VARCHAR(50),

                            status VARCHAR(30) NOT NULL DEFAULT 'ACTIVE',

                            created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                            updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

                            CONSTRAINT fk_classrooms_academic_year
                                FOREIGN KEY (academic_year_id)
                                    REFERENCES academic_years(id)
                                    ON DELETE CASCADE,

                            CONSTRAINT uq_classrooms_academic_year_code
                                UNIQUE (
                                        academic_year_id,
                                        code
                                    ),

                            CONSTRAINT chk_classrooms_code_not_blank
                                CHECK (BTRIM(code) <> ''),

                            CONSTRAINT chk_classrooms_name_not_blank
                                CHECK (BTRIM(name) <> ''),

                            CONSTRAINT chk_classrooms_status
                                CHECK (
                                    status IN (
                                               'ACTIVE',
                                               'INACTIVE'
                                        )
                                    )
);

CREATE INDEX idx_classrooms_academic_year_id
    ON classrooms (academic_year_id);

CREATE INDEX idx_classrooms_status
    ON classrooms (status);