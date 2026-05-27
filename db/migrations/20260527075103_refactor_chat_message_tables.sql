-- +goose Up
SET FOREIGN_KEY_CHECKS = 0;

-- 将房间消息的旧 message_state 表迁移为 recipient 关系表
RENAME TABLE chat_room_message_states TO chat_room_message_recipients;

ALTER TABLE chat_room_message_recipients
	DROP FOREIGN KEY fk_chat_room_message_states_message,
	DROP FOREIGN KEY fk_chat_room_message_states_user,
	DROP INDEX idx_chat_room_message_states_message_id,
	DROP INDEX idx_chat_room_message_states_user_id,
	DROP INDEX idx_chat_room_message_states_state,
	DROP INDEX idx_room_msg_state,
	DROP COLUMN state,
	ADD KEY idx_chat_room_message_recipients_message_id (message_id),
	ADD KEY idx_chat_room_message_recipients_user_id (user_id),
	ADD UNIQUE KEY idx_room_msg (message_id, user_id),
	ADD CONSTRAINT fk_chat_room_message_recipients_message FOREIGN KEY (message_id)
		REFERENCES chat_room_messages (id) ON DELETE CASCADE ON UPDATE CASCADE,
	ADD CONSTRAINT fk_chat_room_message_recipients_user FOREIGN KEY (user_id)
		REFERENCES chat_users (id) ON DELETE RESTRICT ON UPDATE CASCADE;

-- 重构私信表结构，去掉 state 字段，和应用层模型保持一致
ALTER TABLE chat_private_messages
	DROP INDEX idx_chat_private_messages_state,
	DROP COLUMN state;

SET FOREIGN_KEY_CHECKS = 1;

-- +goose Down
SET FOREIGN_KEY_CHECKS = 0;

-- 回滚私信表结构
ALTER TABLE chat_private_messages
	ADD COLUMN state VARCHAR(36) NOT NULL DEFAULT 'undelivered' AFTER sender_id,
	ADD INDEX idx_chat_private_messages_state (state);

-- 回滚房间消息 recipient 表到旧的 message_state 表
ALTER TABLE chat_room_message_recipients
	DROP FOREIGN KEY fk_chat_room_message_recipients_message,
	DROP FOREIGN KEY fk_chat_room_message_recipients_user,
	DROP INDEX idx_chat_room_message_recipients_message_id,
	DROP INDEX idx_chat_room_message_recipients_user_id,
	DROP INDEX idx_room_msg,
	ADD COLUMN state VARCHAR(36) NOT NULL DEFAULT 'undelivered' AFTER user_id,
	ADD INDEX idx_chat_room_message_states_message_id (message_id),
	ADD INDEX idx_chat_room_message_states_user_id (user_id),
	ADD INDEX idx_chat_room_message_states_state (state),
	ADD UNIQUE KEY idx_room_msg_state (message_id, user_id),
	ADD CONSTRAINT fk_chat_room_message_states_message FOREIGN KEY (message_id)
		REFERENCES chat_room_messages (id) ON DELETE CASCADE ON UPDATE CASCADE,
	ADD CONSTRAINT fk_chat_room_message_states_user FOREIGN KEY (user_id)
		REFERENCES chat_users (id) ON DELETE RESTRICT ON UPDATE CASCADE;

RENAME TABLE chat_room_message_recipients TO chat_room_message_states;

SET FOREIGN_KEY_CHECKS = 1;
