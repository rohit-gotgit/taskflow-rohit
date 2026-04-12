-- Insert test user
INSERT INTO users (id, name, email, password)
VALUES (
    uuid_generate_v4(),
    'Test User',
    'test@example.com',
    '$2a$12$CwTycUXWue0Thq9StjUM0uJ8Zt0pniS3pSkeCZMt2rt7NmBGG99G6'
);

-- Insert project
INSERT INTO projects (id, name, description, owner_id)
VALUES (
    uuid_generate_v4(),
    'Demo Project',
    'This is a sample project',
    (SELECT id FROM users WHERE email = 'test@example.com')
);

-- Insert tasks
INSERT INTO tasks (title, status, priority, project_id)
VALUES
('Task 1', 'todo', 'low', (SELECT id FROM projects LIMIT 1)),
('Task 2', 'in_progress', 'medium', (SELECT id FROM projects LIMIT 1)),
('Task 3', 'done', 'high', (SELECT id FROM projects LIMIT 1));