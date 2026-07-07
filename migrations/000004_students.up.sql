-- =========================================================
-- STUDENTS
-- =========================================================

CREATE TABLE students (
                          id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

                          nis VARCHAR(50),
                          nisn VARCHAR(50),

                          full_name VARCHAR(255) NOT NULL,

                          gender VARCHAR(20) NOT NULL,

                          birth_place VARCHAR(150),
                          birth_date DATE,

                          religion VARCHAR(100),
                          address TEXT,

                          status VARCHAR(30) NOT NULL DEFAULT 'ACTIVE',

                          created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                          updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

                          CONSTRAINT chk_students_full_name_not_blank
                              CHECK (BTRIM(full_name) <> ''),

                          CONSTRAINT chk_students_gender
                              CHECK (
                                  gender IN (
                                             'MALE',
                                             'FEMALE'
                                      )
                                  ),

                          CONSTRAINT chk_students_status
                              CHECK (
                                  status IN (
                                             'ACTIVE',
                                             'INACTIVE',
                                             'GRADUATED',
                                             'TRANSFERRED'
                                      )
                                  )
);

CREATE UNIQUE INDEX uq_students_nis
    ON students (nis)
    WHERE nis IS NOT NULL;

CREATE UNIQUE INDEX uq_students_nisn
    ON students (nisn)
    WHERE nisn IS NOT NULL;

CREATE INDEX idx_students_full_name
    ON students (full_name);

CREATE INDEX idx_students_status
    ON students (status);

-- =========================================================
-- CLASSROOM ENROLLMENTS
-- =========================================================

CREATE TABLE classroom_enrollments (
                                       id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

                                       classroom_id UUID NOT NULL,
                                       student_id UUID NOT NULL,

                                       enrolled_at DATE NOT NULL DEFAULT CURRENT_DATE,
                                       left_at DATE,

                                       status VARCHAR(30) NOT NULL DEFAULT 'ACTIVE',

                                       created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                       updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

                                       CONSTRAINT fk_classroom_enrollments_classroom
                                           FOREIGN KEY (classroom_id)
                                               REFERENCES classrooms(id)
                                               ON DELETE CASCADE,

                                       CONSTRAINT fk_classroom_enrollments_student
                                           FOREIGN KEY (student_id)
                                               REFERENCES students(id)
                                               ON DELETE CASCADE,

                                       CONSTRAINT uq_classroom_enrollments
                                           UNIQUE (
                                                   classroom_id,
                                                   student_id
                                               ),

                                       CONSTRAINT chk_classroom_enrollments_status
                                           CHECK (
                                               status IN (
                                                          'ACTIVE',
                                                          'INACTIVE'
                                                   )
                                               ),

                                       CONSTRAINT chk_classroom_enrollments_date_range
                                           CHECK (
                                               left_at IS NULL
                                                   OR left_at >= enrolled_at
                                               )
);

CREATE INDEX idx_classroom_enrollments_classroom_id
    ON classroom_enrollments (classroom_id);

CREATE INDEX idx_classroom_enrollments_student_id
    ON classroom_enrollments (student_id);

CREATE UNIQUE INDEX uq_classroom_enrollments_active_student
    ON classroom_enrollments (student_id)
    WHERE status = 'ACTIVE';