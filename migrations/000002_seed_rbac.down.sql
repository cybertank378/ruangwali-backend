DELETE FROM role_permissions
WHERE role_id IN (
    SELECT id
    FROM roles
    WHERE code IN (
                   'ADMIN',
                   'GURU',
                   'WAKABID_KURIKULUM'
        )
);

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

DELETE FROM roles
WHERE code IN (
               'ADMIN',
               'GURU',
               'WAKABID_KURIKULUM'
    );