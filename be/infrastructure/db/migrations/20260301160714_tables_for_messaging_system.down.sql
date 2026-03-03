-- ============================================
-- ROLLBACK: Chat Application Schema
-- Drops all tables and indexes safely
-- ============================================


-- --------------------------------------------
-- 1️⃣ Drop indexes (safe even if not exist)
-- --------------------------------------------

DROP INDEX IF EXISTS idx_read_messages_user_id;
DROP INDEX IF EXISTS idx_conversation_participants_user_id;
DROP INDEX IF EXISTS idx_messages_sent_at;
DROP INDEX IF EXISTS idx_messages_conversation_id;

-- --------------------------------------------
-- 2️⃣ Drop dependent tables first
-- --------------------------------------------

DROP TABLE IF EXISTS read_messages;
DROP TABLE IF EXISTS conversation_participants;

-- --------------------------------------------
-- 3️⃣ Remove FK constraint from conversations
-- (because it references messages)
-- --------------------------------------------

ALTER TABLE IF EXISTS conversations
DROP CONSTRAINT IF EXISTS fk_conversations_last_message;

-- --------------------------------------------
-- 4️⃣ Drop messages table
-- --------------------------------------------

DROP TABLE IF EXISTS messages;

-- --------------------------------------------
-- 5️⃣ Drop conversations table
-- --------------------------------------------

DROP TABLE IF EXISTS conversations;
