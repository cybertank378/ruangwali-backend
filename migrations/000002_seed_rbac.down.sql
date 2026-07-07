-- =========================================================
-- DEVELOPMENT USER ROLES
-- =========================================================

DELETE FROM user_roles
WHERE user_id IN (
    SELECT id
    FROM users
    WHERE LOWER(BTRIM(email)) IN (
                                  'admin@ruangwali.local',
                                  'guru@ruangwali.local',
                                  'walikelas@ruangwali.local',
                                  'kurikulum@ruangwali.local'
        )
);

-- =========================================================
-- DEVELOPMENT USERS
-- =========================================================

DELETE FROM users
WHERE LOWER(BTRIM(email)) IN (
                              'admin@ruangwali.local',
                              'guru@ruangwali.local',
                              'walikelas@ruangwali.local',
                              'kurikulum@ruangwali.local'
    );

-- =========================================================
-- ROLE PERMISSIONS
-- =========================================================

DELETE FROM role_permissions
WHERE role_id IN (
    SELECT id
    FROM roles
    WHERE code IN (
                   'ADMIN',
                   'GURU',
                   'WALI_KELAS',
                   'WAKABID_KURIKULUM'
        )
);

-- =========================================================
-- PERMISSIONS
-- =========================================================

DELETE FROM permissions
WHERE code IN (
               'user.read',
               'user.create',
               'user.update',
               'user.delete',

               'role.read',
               'role.manage',

               'permission.read',

               'teacher.read',
               'teacher.create',
               'teacher.update',
               'teacher.delete',

               'student.read',
               'student.create',
               'student.update',
               'student.delete',

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
               'homeroom_assignment.manage',

               'integration.read',
               'integration.manage',
               'integration.sync'
    );

-- =========================================================
-- ROLES
-- =========================================================

DELETE FROM roles
WHERE code IN (
               'ADMIN',
               'GURU',
               'WALI_KELAS',
               'WAKABID_KURIKULUM'
    );