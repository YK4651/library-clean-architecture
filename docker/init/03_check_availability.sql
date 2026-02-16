-- Check book availability (Single Source of Truth!)
SELECT
    b.id,
    b.title,
    CASE
        WHEN l.id IS NULL THEN 'AVAILABLE'
        ELSE 'BORROWED'
    END AS availability
FROM books b
LEFT JOIN loans l ON b.id = l.book_id AND l.returned_at IS NULL;
