-- =========================================================
-- ROLES
-- =========================================================

INSERT INTO roles (
    code,
    name,
    description,
    is_system,
    is_active
)
VALUES
    (
        'ADMIN',
        'Administrator',
        'Mengelola sistem, pengguna, akses, dan integrasi',
        TRUE,
        TRUE
    ),
    (
        'GURU',
        'Guru',
        'Mengakses data akademik sesuai penugasan mengajar',
        TRUE,
        TRUE
    ),
    (
        'WAKABID_KURIKULUM',
        'Wakabid Kurikulum',
        'Mengelola fondasi akademik, kurikulum, dan penugasan guru',
        TRUE,
        TRUE
    )
    ON CONFLICT (code) DO NOTHING;

-- =========================================================
-- PERMISSIONS
-- =========================================================

INSERT INTO permissions (
    code,
    resource,
    action,
    description
)
VALUES
    (
        'user.read',
        'user',
        'read',
        'Melihat data pengguna'
    ),
    (
        'user.create',
        'user',
        'create',
        'Membuat pengguna'
    ),
    (
        'user.update',
        'user',
        'update',
        'Memperbarui pengguna'
    ),
    (
        'user.delete',
        'user',
        'delete',
        'Menghapus pengguna'
    ),
    (
        'role.read',
        'role',
        'read',
        'Melihat role'
    ),
    (
        'role.manage',
        'role',
        'manage',
        'Mengelola role'
    ),
    (
        'permission.read',
        'permission',
        'read',
        'Melihat permission'
    ),
    (
        'teacher.read',
        'teacher',
        'read',
        'Melihat data guru'
    ),
    (
        'teacher.create',
        'teacher',
        'create',
        'Membuat data guru'
    ),
    (
        'teacher.update',
        'teacher',
        'update',
        'Memperbarui data guru'
    ),
    (
        'teacher.delete',
        'teacher',
        'delete',
        'Menghapus data guru'
    ),
    (
        'student.read',
        'student',
        'read',
        'Melihat data siswa'
    ),
    (
        'student.create',
        'student',
        'create',
        'Membuat data siswa'
    ),
    (
        'student.update',
        'student',
        'update',
        'Memperbarui data siswa'
    ),
    (
        'student.delete',
        'student',
        'delete',
        'Menghapus data siswa'
    ),
    (
        'academic_year.read',
        'academic_year',
        'read',
        'Melihat tahun ajaran'
    ),
    (
        'academic_year.manage',
        'academic_year',
        'manage',
        'Mengelola tahun ajaran'
    ),
    (
        'semester.read',
        'semester',
        'read',
        'Melihat semester'
    ),
    (
        'semester.manage',
        'semester',
        'manage',
        'Mengelola semester'
    ),
    (
        'subject.read',
        'subject',
        'read',
        'Melihat mata pelajaran'
    ),
    (
        'subject.manage',
        'subject',
        'manage',
        'Mengelola mata pelajaran'
    ),
    (
        'classroom.read',
        'classroom',
        'read',
        'Melihat kelas'
    ),
    (
        'classroom.manage',
        'classroom',
        'manage',
        'Mengelola kelas'
    ),
    (
        'teaching_assignment.read',
        'teaching_assignment',
        'read',
        'Melihat penugasan mengajar'
    ),
    (
        'teaching_assignment.manage',
        'teaching_assignment',
        'manage',
        'Mengelola penugasan mengajar'
    ),
    (
        'homeroom_assignment.read',
        'homeroom_assignment',
        'read',
        'Melihat penugasan wali kelas'
    ),
    (
        'homeroom_assignment.manage',
        'homeroom_assignment',
        'manage',
        'Mengelola penugasan wali kelas'
    ),
    (
        'integration.read',
        'integration',
        'read',
        'Melihat konfigurasi dan aktivitas integrasi'
    ),
    (
        'integration.manage',
        'integration',
        'manage',
        'Mengelola integrasi'
    ),
    (
        'integration.sync',
        'integration',
        'sync',
        'Menjalankan sinkronisasi integrasi'
    )
    ON CONFLICT (code) DO NOTHING;

-- =========================================================
-- ADMIN PERMISSIONS
-- =========================================================

INSERT INTO role_permissions (
    role_id,
    permission_id
)
SELECT
    r.id,
    p.id
FROM roles r
         CROSS JOIN permissions p
WHERE r.code = 'ADMIN'
    ON CONFLICT DO NOTHING;

-- =========================================================
-- GURU PERMISSIONS
-- Scope tetap divalidasi oleh domain policy.
-- =========================================================

INSERT INTO role_permissions (
    role_id,
    permission_id
)
SELECT
    r.id,
    p.id
FROM roles r
         JOIN permissions p
              ON p.code IN (
                            'teacher.read',
                            'student.read',
                            'academic_year.read',
                            'semester.read',
                            'subject.read',
                            'classroom.read',
                            'teaching_assignment.read',
                            'homeroom_assignment.read'
                  )
WHERE r.code = 'GURU'
    ON CONFLICT DO NOTHING;

-- =========================================================
-- WAKABID KURIKULUM PERMISSIONS
-- =========================================================

INSERT INTO role_permissions (
    role_id,
    permission_id
)
SELECT
    r.id,
    p.id
FROM roles r
         JOIN permissions p
              ON p.code IN (
                            'teacher.read',
                            'student.read',
                            'academic_year.read',
                            'academic_year.manage',
                            'semester.read',
                            'semester.manage',
                            'subject.read',
                            'subject.manage',
                            'classroom.read',
                            'classroom.manage',
                            'teaching_assignment.read',
                            'teaching_assignment.manage',
                            'homeroom_assignment.read',
                            'homeroom_assignment.manage'
                  )
WHERE r.code = 'WAKABID_KURIKULUM'
    ON CONFLICT DO NOTHING;