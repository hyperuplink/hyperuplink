-- Applied after the users exist, since every post here points at one of them
-- by username.

BEGIN;


INSERT INTO categories (id, name, slug, position) VALUES
  ('a0000000-0000-4000-8000-000000000001', 'General', 'general', 0),
  ('a0000000-0000-4000-8000-000000000002', 'Support', 'support', 1);

INSERT INTO forums (id, name, slug, category_id, position, description) VALUES
  ('b0000000-0000-4000-8000-000000000001', 'Announcements', 'announcements',
   'a0000000-0000-4000-8000-000000000001', 0,
   'Board news, release notes and downtime notices.'),
  ('b0000000-0000-4000-8000-000000000002', 'Introductions', 'introductions',
   'a0000000-0000-4000-8000-000000000001', 1,
   'New here? Tell everybody who you are.'),
  ('b0000000-0000-4000-8000-000000000003', 'Help', 'help',
   'a0000000-0000-4000-8000-000000000002', 0,
   'Questions about using the board, answered by whoever knows.');

INSERT INTO topics
  (id, short_id, name, slug, forum_id, author_id, kind, pinned, text, html)
VALUES
  ('c0000000-0000-4000-8000-000000000001', 'wlcm0001',
   'Welcome to the board', 'welcome-to-the-board',
   'b0000000-0000-4000-8000-000000000001',
   (SELECT id FROM users WHERE username = 'sysop'),
   'regular', true,
   'This board runs on Hyperuplink. Post in Markdown, keep the tone civil, and read the house rules before you start a thread.',
   '<p>This board runs on Hyperuplink. Post in Markdown, keep the tone civil, and read the house rules before you start a thread.</p>'),

  ('c0000000-0000-4000-8000-000000000002', 'rules001',
   'House rules', 'house-rules',
   'b0000000-0000-4000-8000-000000000001',
   (SELECT id FROM users WHERE username = 'sysop'),
   'regular', false,
   'Stay on topic, attack arguments and never people, and use the report button instead of starting a fight in the thread.',
   '<p>Stay on topic, attack arguments and never people, and use the report button instead of starting a fight in the thread.</p>'),

  ('c0000000-0000-4000-8000-000000000003', 'hello001',
   'Say hello here', 'say-hello-here',
   'b0000000-0000-4000-8000-000000000002',
   (SELECT id FROM users WHERE username = 'vera'),
   'regular', false,
   'I have been reading for a while and finally signed up. I run a small board of my own and came here for the themes.',
   '<p>I have been reading for a while and finally signed up. I run a small board of my own and came here for the themes.</p>'),

  ('c0000000-0000-4000-8000-000000000004', 'poll0001',
   'Which theme should be the default?', 'which-theme-should-be-the-default',
   'b0000000-0000-4000-8000-000000000003',
   (SELECT id FROM users WHERE username = 'vera'),
   'poll', false,
   'The board ships with more themes than anybody needs, so pick the one that should greet a first-time visitor.',
   '<p>The board ships with more themes than anybody needs, so pick the one that should greet a first-time visitor.</p>'),

  ('c0000000-0000-4000-8000-000000000005', 'attach01',
   'How do I add a screenshot to a post?', 'how-do-i-add-a-screenshot-to-a-post',
   'b0000000-0000-4000-8000-000000000003',
   (SELECT id FROM users WHERE username = 'juno'),
   'regular', false,
   'I wrote the post, but I cannot work out where the attachment goes. Is there a size limit I am running into?',
   '<p>I wrote the post, but I cannot work out where the attachment goes. Is there a size limit I am running into?</p>');

UPDATE topics
SET poll_options = ARRAY[
  'macOS 9',
  'Windows 9x',
  'CDE',
  'Hobbit'
]::varchar(78)[]
WHERE id = 'c0000000-0000-4000-8000-000000000004';

INSERT INTO replies (id, short_id, topic_id, author_id, text, html) VALUES
  ('d0000000-0000-4000-8000-000000000001', 'rep00001',
   'c0000000-0000-4000-8000-000000000001',
   (SELECT id FROM users WHERE username = 'vera'),
   'Good to be here. The macOS 9 theme with a light colorscheme is doing a lot of work for me already.',
   '<p>Good to be here. The macOS 9 theme with a light colorscheme is doing a lot of work for me already.</p>'),

  ('d0000000-0000-4000-8000-000000000002', 'rep00002',
   'c0000000-0000-4000-8000-000000000001',
   (SELECT id FROM users WHERE username = 'bitrot'),
   'Same, though I keep switching back to Windows 3.x whenever I feel nostalgic, which is more often than I would like to admit.',
   '<p>Same, though I keep switching back to Windows 3.x whenever I feel nostalgic, which is more often than I would like to admit.</p>'),

  ('d0000000-0000-4000-8000-000000000003', 'rep00003',
   'c0000000-0000-4000-8000-000000000003',
   (SELECT id FROM users WHERE username = 'juno'),
   'Welcome! What does your own board run on?',
   '<p>Welcome! What does your own board run on?</p>'),

  ('d0000000-0000-4000-8000-000000000004', 'rep00004',
   'c0000000-0000-4000-8000-000000000005',
   (SELECT id FROM users WHERE username = 'sysop'),
   'The attachment field sits underneath the editor, and the administrator sets both the maximum size and the formats that are accepted.',
   '<p>The attachment field sits underneath the editor, and the administrator sets both the maximum size and the formats that are accepted.</p>');

INSERT INTO postevents (type, author_id, target, topic_id, selection) VALUES
  ('pollvote', (SELECT id FROM users WHERE username = 'vera'),
   'topic', 'c0000000-0000-4000-8000-000000000004', 0),
  ('pollvote', (SELECT id FROM users WHERE username = 'juno'),
   'topic', 'c0000000-0000-4000-8000-000000000004', 0),
  ('pollvote', (SELECT id FROM users WHERE username = 'bitrot'),
   'topic', 'c0000000-0000-4000-8000-000000000004', 2);

INSERT INTO postevents (type, author_id, target, topic_id, selection) VALUES
  ('report', (SELECT id FROM users WHERE username = 'bitrot'),
   'topic', 'c0000000-0000-4000-8000-000000000005', 0);

INSERT INTO groups (id, name) VALUES
  ('members', 'Members'),
  ('veterans', 'Veterans');

INSERT INTO permissions (group_id, category_id, bits) VALUES
  ('veterans', 'a0000000-0000-4000-8000-000000000002', B'110');

COMMIT;

REFRESH MATERIALIZED VIEW vforums;
REFRESH MATERIALIZED VIEW vtopics;
REFRESH MATERIALIZED VIEW vreplies;
