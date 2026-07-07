INSERT INTO permissions (code, resource, action, description) VALUES
('dashboard.read','dashboard','read','Melihat dashboard'),
('student.read','student','read','Melihat siswa'),
('student.create','student','create','Menambah siswa'),
('student.update','student','update','Mengubah siswa'),
('student.delete','student','delete','Menghapus siswa'),
('student.import','student','import','Impor siswa'),
('student.export','student','export','Ekspor siswa'),
('guardian.read','guardian','read','Melihat wali'),
('guardian.create','guardian','create','Menambah wali'),
('guardian.update','guardian','update','Mengubah wali'),
('guardian.delete','guardian','delete','Menghapus wali'),
('classroom.read','classroom','read','Melihat kelas'),
('classroom.update','classroom','update','Mengubah kelas'),
('classroom.organization.read','classroom.organization','read','Melihat struktur kelas'),
('classroom.organization.manage','classroom.organization','manage','Mengelola struktur kelas'),
('classroom.schedule.read','classroom.schedule','read','Melihat jadwal'),
('classroom.schedule.manage','classroom.schedule','manage','Mengelola jadwal'),
('classroom.duty.read','classroom.duty','read','Melihat piket'),
('classroom.duty.manage','classroom.duty','manage','Mengelola piket'),
('classroom.seating.read','classroom.seating','read','Melihat denah'),
('classroom.seating.manage','classroom.seating','manage','Mengelola denah'),
('classroom.agreement.read','classroom.agreement','read','Melihat kesepakatan'),
('classroom.agreement.manage','classroom.agreement','manage','Mengelola kesepakatan'),
('attendance.read','attendance','read','Melihat absensi'),
('attendance.record','attendance','record','Mencatat absensi'),
('attendance.update','attendance','update','Mengubah absensi'),
('attendance.delete','attendance','delete','Menghapus absensi'),
('attendance.export','attendance','export','Ekspor absensi'),
('assessment.read','assessment','read','Melihat nilai'),
('assessment.record','assessment','record','Mencatat nilai'),
('assessment.update','assessment','update','Mengubah nilai'),
('assessment.delete','assessment','delete','Menghapus nilai'),
('assessment.export','assessment','export','Ekspor nilai'),
('achievement.read','achievement','read','Melihat prestasi'),
('achievement.create','achievement','create','Menambah prestasi'),
('achievement.update','achievement','update','Mengubah prestasi'),
('achievement.delete','achievement','delete','Menghapus prestasi'),
('violation.read','violation','read','Melihat pelanggaran'),
('violation.create','violation','create','Menambah pelanggaran'),
('violation.update','violation','update','Mengubah pelanggaran'),
('violation.delete','violation','delete','Menghapus pelanggaran'),
('journal.read','journal','read','Melihat jurnal'),
('journal.create','journal','create','Menambah jurnal'),
('journal.update','journal','update','Mengubah jurnal'),
('journal.delete','journal','delete','Menghapus jurnal'),
('school.read','school','read','Melihat sekolah'),
('school.update','school','update','Mengubah sekolah'),
('user.read','user','read','Melihat user'),
('user.create','user','create','Menambah user'),
('user.update','user','update','Mengubah user'),
('user.deactivate','user','deactivate','Menonaktifkan user'),
('role.read','role','read','Melihat role'),
('role.manage','role','manage','Mengelola role'),
('role.assign','role','assign','Menetapkan role'),
('audit.read','audit','read','Melihat audit'),
('backup.create','backup','create','Membuat backup'),
('backup.restore','backup','restore','Memulihkan backup'),
('system.reset','system','reset','Reset data sistem')
ON CONFLICT (code) DO NOTHING;

INSERT INTO roles (code, name, description, is_system)
VALUES
('SUPER_ADMIN','Super Admin','Operator platform RuangWali',TRUE),
('SCHOOL_ADMIN','School Admin','Administrator tenant sekolah',TRUE),
('HOMEROOM_TEACHER','Homeroom Teacher','Wali kelas',TRUE),
('SUBJECT_TEACHER','Subject Teacher','Guru mata pelajaran',TRUE)
ON CONFLICT DO NOTHING;

-- SUPER_ADMIN: semua permission
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r CROSS JOIN permissions p
WHERE r.code = 'SUPER_ADMIN'
ON CONFLICT DO NOTHING;

-- SCHOOL_ADMIN: semua tenant capability kecuali system.reset
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r CROSS JOIN permissions p
WHERE r.code = 'SCHOOL_ADMIN'
  AND p.code <> 'system.reset'
ON CONFLICT DO NOTHING;

-- HOMEROOM_TEACHER: administrasi kelas
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.code = ANY(ARRAY[
'dashboard.read',
'student.read','student.create','student.update','student.import','student.export',
'guardian.read','guardian.create','guardian.update',
'classroom.read','classroom.organization.read','classroom.organization.manage',
'classroom.schedule.read','classroom.schedule.manage',
'classroom.duty.read','classroom.duty.manage',
'classroom.seating.read','classroom.seating.manage',
'classroom.agreement.read','classroom.agreement.manage',
'attendance.read','attendance.record','attendance.update','attendance.export',
'assessment.read','assessment.record','assessment.update','assessment.export',
'achievement.read','achievement.create','achievement.update',
'violation.read','violation.create','violation.update',
'journal.read','journal.create','journal.update',
'school.read'
])
WHERE r.code = 'HOMEROOM_TEACHER'
ON CONFLICT DO NOTHING;

-- SUBJECT_TEACHER: capability dasar; fine-grained scope tetap diverifikasi policy
INSERT INTO role_permissions (role_id, permission_id)
SELECT r.id, p.id
FROM roles r
JOIN permissions p ON p.code = ANY(ARRAY[
'dashboard.read',
'student.read',
'classroom.read',
'classroom.schedule.read',
'classroom.duty.read',
'classroom.seating.read',
'classroom.agreement.read',
'attendance.read',
'assessment.read','assessment.record','assessment.update',
'achievement.read',
'violation.read',
'school.read'
])
WHERE r.code = 'SUBJECT_TEACHER'
ON CONFLICT DO NOTHING;
