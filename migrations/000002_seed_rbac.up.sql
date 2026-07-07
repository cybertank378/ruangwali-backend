

INSERT INTO roles (
    code,
    name,
    description,
    is_system,
    is_active
)
VALUES-- =========================================================
-- ROLES
-- =========================================================
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
        'WALI_KELAS',
        'Wali Kelas',
        'Mengelola dan memantau data siswa pada kelas yang menjadi tanggung jawabnya',
        TRUE,
        TRUE
    ),
    (
        'WAKABID_KURIKULUM',
        'Wakabid Kurikulum',
        'Mengelola fondasi akademik, kurikulum, kelas, dan penugasan guru',
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
                            'teaching_assignment.read'
                  )
WHERE r.code = 'GURU'
    ON CONFLICT DO NOTHING;

-- =========================================================
-- WALI KELAS PERMISSIONS
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
                            'student.update',
                            'academic_year.read',
                            'semester.read',
                            'subject.read',
                            'classroom.read',
                            'teaching_assignment.read',
                            'homeroom_assignment.read'
                  )
WHERE r.code = 'WALI_KELAS'
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

-- =========================================================
-- DEVELOPMENT USERS
--
-- Password awal seluruh akun:
-- RuangWali@2026
--
-- Argon2id:
-- memory      = 65536 KiB
-- iterations  = 3
-- parallelism = 2
-- salt length = 16 bytes
-- key length  = 32 bytes
-- =========================================================

INSERT INTO users (
    email,
    password_hash,
    status
)
SELECT
    seed.email,
    seed.password_hash,
    seed.status
FROM (
         VALUES
             (
                 'admin@ruangwali.local',
                 '$argon2id$v=19$m=65536,t=3,p=2$8Tfyq7eJT6sR3wsAmf7e1Q$upMO7fyvFzKTIljhj5+R2Zshm3QnPBV8sot6IzZn7IA',
                 'ACTIVE'
             ),
             (
                 'guru@ruangwali.local',
                 '$argon2id$v=19$m=65536,t=3,p=2$Eojx/kZVRHhCViZ6yImAjA$jXp4fs6TqaK/WNir0UYdDQSDtJQRtdy2z9RmIfAbhFQ',
                 'ACTIVE'
             ),
             (
                 'walikelas@ruangwali.local',
                 '$argon2id$v=19$m=65536,t=3,p=2$j7vhlfcFNRirXrPAi/9UZw$FaiYKQ0zGsXhhFmh0cBp72C/SdvJZFmTfBcN8JNrP/c',
                 'ACTIVE'
             ),
             (
                 'kurikulum@ruangwali.local',
                 '$argon2id$v=19$m=65536,t=3,p=2$Fk9HDF7Q2QIPtfhvbVFt8A$PFiADm1maDNgHj0Z93eqG/YA8QZJ/EfqwDpLGcAV5kM',
                 'ACTIVE'
             )
     ) AS seed (
                email,
                password_hash,
                status
    )
WHERE NOT EXISTS (
    SELECT 1
    FROM users u
    WHERE LOWER(BTRIM(u.email)) =
          LOWER(BTRIM(seed.email))
);

-- =========================================================
-- ADMIN USER ROLE
-- =========================================================

INSERT INTO user_roles (
    user_id,
    role_id,
    assigned_by
)
SELECT
    u.id,
    r.id,
    NULL
FROM users u
         JOIN roles r
              ON r.code = 'ADMIN'
WHERE LOWER(BTRIM(u.email)) =
      'admin@ruangwali.local'
    ON CONFLICT DO NOTHING;

-- =========================================================
-- GURU USER ROLE
-- =========================================================

INSERT INTO user_roles (
    user_id,
    role_id,
    assigned_by
)
SELECT
    u.id,
    r.id,
    admin_user.id
FROM users u
         JOIN roles r
              ON r.code = 'GURU'
         LEFT JOIN users admin_user
                   ON LOWER(BTRIM(admin_user.email)) =
                      'admin@ruangwali.local'
WHERE LOWER(BTRIM(u.email)) =
      'guru@ruangwali.local'
    ON CONFLICT DO NOTHING;

-- =========================================================
-- WALI KELAS USER ROLE
-- =========================================================

INSERT INTO user_roles (
    user_id,
    role_id,
    assigned_by
)
SELECT
    u.id,
    r.id,
    admin_user.id
FROM users u
         JOIN roles r
              ON r.code = 'WALI_KELAS'
         LEFT JOIN users admin_user
                   ON LOWER(BTRIM(admin_user.email)) =
                      'admin@ruangwali.local'
WHERE LOWER(BTRIM(u.email)) =
      'walikelas@ruangwali.local'
    ON CONFLICT DO NOTHING;

-- =========================================================
-- WAKABID KURIKULUM USER ROLE
-- =========================================================

INSERT INTO user_roles (
    user_id,
    role_id,
    assigned_by
)
SELECT
    u.id,
    r.id,
    admin_user.id
FROM users u
         JOIN roles r
              ON r.code = 'WAKABID_KURIKULUM'
         LEFT JOIN users admin_user
                   ON LOWER(BTRIM(admin_user.email)) =
                      'admin@ruangwali.local'
WHERE LOWER(BTRIM(u.email)) =
      'kurikulum@ruangwali.local'
    ON CONFLICT DO NOTHING;