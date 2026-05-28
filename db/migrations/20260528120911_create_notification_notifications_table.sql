-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS notification_notifications
(
	id           CHAR(36)    NOT NULL PRIMARY KEY,
	recipient_id CHAR(36)    NOT NULL,
	state        VARCHAR(36) NOT NULL DEFAULT 'undelivered',
	raw_payload  LONGBLOB    NULL,
	KEY idx_notification_notifications_recipient_id (recipient_id),
	KEY idx_notification_notifications_state (state),
	KEY idx_notification_notifications_recipient_state (recipient_id, state)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS notification_notifications;
-- +goose StatementEnd
