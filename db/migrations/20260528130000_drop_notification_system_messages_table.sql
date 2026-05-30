-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS notification_system_messages;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS notification_system_messages
(
    id           CHAR(36)     NOT NULL PRIMARY KEY,
    recipient_id CHAR(36)     NOT NULL,
    state        VARCHAR(36)  NOT NULL,
    content      VARCHAR(255) NOT NULL,
    sent_at      DATETIME(3)  NULL,
    KEY idx_notification_system_messages_recipient_id (recipient_id),
    KEY idx_notification_system_messages_state (state),
    CONSTRAINT fk_notification_system_messages_recipient FOREIGN KEY (recipient_id)
        REFERENCES chat_users (id) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;
-- +goose StatementEnd

