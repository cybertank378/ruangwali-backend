-- =========================================================
-- TEACHING ASSIGNMENTS
-- =========================================================

CREATE TABLE teaching_assignments (
                                      id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

                                      teacher_id UUID NOT NULL,
                                      classroom_id UUID NOT NULL,
                                      subject_id UUID NOT NULL,
                                      semester_id UUID NOT NULL,

                                      status VARCHAR(30) NOT NULL DEFAULT 'ACTIVE',

                                      created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                      updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

                                      CONSTRAINT fk_teaching_assignments_teacher
                                          FOREIGN KEY (teacher_id)
                                              REFERENCES teachers(id)
                                              ON DELETE CASCADE,

                                      CONSTRAINT fk_teaching_assignments_classroom
                                          FOREIGN KEY (classroom_id)
                                              REFERENCES classrooms(id)
                                              ON DELETE CASCADE,

                                      CONSTRAINT fk_teaching_assignments_subject
                                          FOREIGN KEY (subject_id)
                                              REFERENCES subjects(id)
                                              ON DELETE CASCADE,

                                      CONSTRAINT fk_teaching_assignments_semester
                                          FOREIGN KEY (semester_id)
                                              REFERENCES semesters(id)
                                              ON DELETE CASCADE,

                                      CONSTRAINT uq_teaching_assignments
                                          UNIQUE (
                                                  teacher_id,
                                                  classroom_id,
                                                  subject_id,
                                                  semester_id
                                              ),

                                      CONSTRAINT chk_teaching_assignments_status
                                          CHECK (
                                              status IN (
                                                         'ACTIVE',
                                                         'INACTIVE'
                                                  )
                                              )
);

CREATE INDEX idx_teaching_assignments_teacher_id
    ON teaching_assignments (teacher_id);

CREATE INDEX idx_teaching_assignments_classroom_id
    ON teaching_assignments (classroom_id);

CREATE INDEX idx_teaching_assignments_subject_id
    ON teaching_assignments (subject_id);

CREATE INDEX idx_teaching_assignments_semester_id
    ON teaching_assignments (semester_id);

-- =========================================================
-- HOMEROOM ASSIGNMENTS
-- WALI_KELAS = DOMAIN ASSIGNMENT, BUKAN GLOBAL ROLE
-- =========================================================

CREATE TABLE homeroom_assignments (
                                      id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

                                      teacher_id UUID NOT NULL,
                                      classroom_id UUID NOT NULL,

                                      assigned_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                      ended_at TIMESTAMPTZ,

                                      created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                                      updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

                                      CONSTRAINT fk_homeroom_assignments_teacher
                                          FOREIGN KEY (teacher_id)
                                              REFERENCES teachers(id)
                                              ON DELETE CASCADE,

                                      CONSTRAINT fk_homeroom_assignments_classroom
                                          FOREIGN KEY (classroom_id)
                                              REFERENCES classrooms(id)
                                              ON DELETE CASCADE,

                                      CONSTRAINT chk_homeroom_assignments_date_range
                                          CHECK (
                                              ended_at IS NULL
                                                  OR ended_at >= assigned_at
                                              )
);

CREATE UNIQUE INDEX uq_homeroom_assignments_active_classroom
    ON homeroom_assignments (classroom_id)
    WHERE ended_at IS NULL;

CREATE UNIQUE INDEX uq_homeroom_assignments_active_teacher
    ON homeroom_assignments (teacher_id)
    WHERE ended_at IS NULL;

CREATE INDEX idx_homeroom_assignments_teacher_id
    ON homeroom_assignments (teacher_id);

CREATE INDEX idx_homeroom_assignments_classroom_id
    ON homeroom_assignments (classroom_id);