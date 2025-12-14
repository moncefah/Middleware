PRAGMA foreign_keys = ON;

-- =========================
-- AGENDAS
-- =========================

INSERT INTO agendas (id, name, uca_id) VALUES
('11111111-1111-1111-1111-111111111111', 'Mathématiques', 'UCA-MATH'),
('22222222-2222-2222-2222-222222222222', 'Informatique', 'UCA-INFO'),
('33333333-3333-3333-3333-333333333333', 'Physique', 'UCA-PHYS');

-- =========================
-- ALERTS
-- =========================

INSERT INTO alerts (id, agenda_id, email) VALUES
('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', '11111111-1111-1111-1111-111111111111', 'alice@example.com'),
('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', '11111111-1111-1111-1111-111111111111', 'bob@example.com'),
('cccccccc-cccc-cccc-cccc-cccccccccccc', '22222222-2222-2222-2222-222222222222', 'charlie@example.com');

