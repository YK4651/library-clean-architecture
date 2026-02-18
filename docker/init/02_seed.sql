-- docker/init/02_seed.sql
-- Sample Data for Testing
-- Note: ID uses ULID (26 characters)

-- Users (status: 1=ACTIVE, 2=SUSPENDED, 3=INACTIVE)
INSERT INTO users (id, name, email, status, max_loans) VALUES
('01HQXK5V8N3YZGJB4QWERT001', '田中太郎', 'tanaka@example.com', 1, 5),
('01HQXK5V8N3YZGJB4QWERT002', '佐藤花子', 'sato@example.com', 1, 5),
('01HQXK5V8N3YZGJB4QWERT003', '鈴木一郎', 'suzuki@example.com', 2, 5),
('01HQXK5V8N3YZGJB4QWERT004', '高橋美咲', 'takahashi@example.com', 1, 3);

-- Books
INSERT INTO books (id, title, author, isbn) VALUES
('01HQXK6A2M5YZGJB4QWERT001', 'Clean Code', 'Robert C. Martin', '9780132350884'),
('01HQXK6A2M5YZGJB4QWERT002', 'Clean Architecture', 'Robert C. Martin', '9780134494166'),
('01HQXK6A2M5YZGJB4QWERT003', 'Domain-Driven Design', 'Eric Evans', '9780321125217'),
('01HQXK6A2M5YZGJB4QWERT004', 'Refactoring', 'Martin Fowler', '9780134757599'),
('01HQXK6A2M5YZGJB4QWERT005', 'The Pragmatic Programmer', 'David Thomas', '9780135957059');

-- Loans (some active, some returned, some overdue)
INSERT INTO loans (id, user_id, book_id, borrowed_at, due_date, returned_at, late_fee) VALUES
-- Active loan (book-001 is borrowed)
('01HQXK7B3P6YZGJB4QWERT001', '01HQXK5V8N3YZGJB4QWERT001', '01HQXK6A2M5YZGJB4QWERT001', NOW(), DATE_ADD(NOW(), INTERVAL 14 DAY), NULL, 0),
-- Returned loan
('01HQXK7B3P6YZGJB4QWERT002', '01HQXK5V8N3YZGJB4QWERT001', '01HQXK6A2M5YZGJB4QWERT002', DATE_SUB(NOW(), INTERVAL 30 DAY), DATE_SUB(NOW(), INTERVAL 16 DAY), DATE_SUB(NOW(), INTERVAL 20 DAY), 0),
-- Overdue loan (book-003 is overdue)
('01HQXK7B3P6YZGJB4QWERT003', '01HQXK5V8N3YZGJB4QWERT002', '01HQXK6A2M5YZGJB4QWERT003', DATE_SUB(NOW(), INTERVAL 20 DAY), DATE_SUB(NOW(), INTERVAL 6 DAY), NULL, 0);