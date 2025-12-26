-- migrations/sql/20251226035601_init_complete_schema.sql
-- +goose Up
-- 初始化完整数据库结构，按依赖顺序创建表，并补充约束与索引

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- 幂等初始化：表存在则跳过创建，保留已有数据

-- 1. 用户基表（按领域划分）
CREATE TABLE IF NOT EXISTS authorization_users
(
    id                 CHAR(36)     NOT NULL PRIMARY KEY,
    email              VARCHAR(254) NOT NULL,
    number             VARCHAR(20)  NOT NULL,
    password_encrypted VARCHAR(255) NOT NULL,
    signed_up_at       DATETIME(3)  NULL,
    UNIQUE KEY idx_authorization_users_email (email),
    UNIQUE KEY idx_authorization_users_number (number)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS chat_users
(
    id CHAR(36) NOT NULL PRIMARY KEY
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS friendship_users
(
    id CHAR(36) NOT NULL PRIMARY KEY
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS profile_users
(
    id CHAR(36) NOT NULL PRIMARY KEY
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS roomship_users
(
    id CHAR(36) NOT NULL PRIMARY KEY
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

-- 2. 房间（按领域划分）
CREATE TABLE IF NOT EXISTS chat_rooms
(
    id CHAR(36) NOT NULL PRIMARY KEY
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS profile_rooms
(
    id CHAR(36) NOT NULL PRIMARY KEY
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS roomship_rooms
(
    id                 CHAR(36)         NOT NULL PRIMARY KEY,
    number             VARCHAR(20)      NOT NULL,
    owner_id           CHAR(36)         NOT NULL,
    password_encrypted VARCHAR(255)     NULL,
    max_member_count   BIGINT DEFAULT 0 NOT NULL,
    created_at         DATETIME(3)      NULL,
    UNIQUE KEY idx_roomship_rooms_number (number),
    KEY idx_roomship_rooms_owner_id (owner_id),
    CONSTRAINT fk_roomship_rooms_owner FOREIGN KEY (owner_id)
        REFERENCES roomship_users (id) ON DELETE RESTRICT ON UPDATE CASCADE,
    CHECK (`max_member_count` >= 0)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

-- 3. 个人资料（Profile）
CREATE TABLE IF NOT EXISTS profile_user_profiles
(
    id           CHAR(36)    NOT NULL PRIMARY KEY,
    name         LONGTEXT    NULL,
    gender       BIGINT      NULL,
    email        LONGTEXT    NULL,
    phone_number LONGTEXT    NULL,
    address      LONGTEXT    NULL,
    sign         LONGTEXT    NULL,
    signed_up_at DATETIME(3) NULL,
    CONSTRAINT fk_profile_user_profiles_user FOREIGN KEY (id)
        REFERENCES profile_users (id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS profile_room_profiles
(
    id           CHAR(36)    NOT NULL PRIMARY KEY,
    name         LONGTEXT    NULL,
    introduction LONGTEXT    NULL,
    created_at   DATETIME(3) NULL,
    CONSTRAINT fk_profile_room_profiles_room FOREIGN KEY (id)
        REFERENCES profile_rooms (id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

-- 4. 关系表
CREATE TABLE IF NOT EXISTS chat_friendships
(
    id       CHAR(36) NOT NULL PRIMARY KEY,
    user_id1 CHAR(36) NOT NULL,
    user_id2 CHAR(36) NOT NULL,
    UNIQUE KEY idx_friendship_pair (user_id1, user_id2),
    KEY idx_chat_friendships_user2 (user_id2),
    CONSTRAINT fk_chat_friendships_user1 FOREIGN KEY (user_id1)
        REFERENCES chat_users (id) ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT fk_chat_friendships_user2 FOREIGN KEY (user_id2)
        REFERENCES chat_users (id) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS friendship_friendships
(
    id         CHAR(36)    NOT NULL PRIMARY KEY,
    user_id1   CHAR(36)    NOT NULL,
    user_id2   CHAR(36)    NOT NULL,
    created_at DATETIME(3) NOT NULL,
    UNIQUE KEY idx_friendship_pair (user_id1, user_id2),
    KEY idx_friendship_friendships_created_at (created_at),
    CONSTRAINT fk_friendship_friendships_user1 FOREIGN KEY (user_id1)
        REFERENCES friendship_users (id) ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT fk_friendship_friendships_user2 FOREIGN KEY (user_id2)
        REFERENCES friendship_users (id) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS profile_roomships
(
    id      CHAR(36) NOT NULL PRIMARY KEY,
    role    CHAR(20) NOT NULL,
    room_id CHAR(36) NOT NULL,
    user_id CHAR(36) NOT NULL,
    UNIQUE KEY idx_roomship_pair (room_id, user_id),
    KEY idx_profile_roomships_role (role),
    CONSTRAINT fk_profile_roomships_room FOREIGN KEY (room_id)
        REFERENCES profile_rooms (id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_profile_roomships_user FOREIGN KEY (user_id)
        REFERENCES profile_users (id) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS roomship_roomships
(
    id         CHAR(36)    NOT NULL PRIMARY KEY,
    room_id    CHAR(36)    NOT NULL,
    user_id    CHAR(36)    NOT NULL,
    role       CHAR(20)    NOT NULL,
    created_at DATETIME(3) NOT NULL,
    UNIQUE KEY idx_roomship_pair (room_id, user_id),
    KEY idx_roomship_roomships_created_at (created_at),
    KEY idx_roomship_roomships_role (role),
    CONSTRAINT fk_roomship_roomships_room FOREIGN KEY (room_id)
        REFERENCES roomship_rooms (id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_roomship_roomships_user FOREIGN KEY (user_id)
        REFERENCES roomship_users (id) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS chat_roomships
(
    id      CHAR(36) NOT NULL PRIMARY KEY,
    room_id CHAR(36) NOT NULL,
    user_id CHAR(36) NOT NULL,
    UNIQUE KEY idx_roomship_pair (room_id, user_id),
    CONSTRAINT fk_chat_roomships_room FOREIGN KEY (room_id)
        REFERENCES chat_rooms (id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_chat_roomships_user FOREIGN KEY (user_id)
        REFERENCES chat_users (id) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

-- 5. 消息表
CREATE TABLE IF NOT EXISTS chat_private_messages
(
    id           CHAR(36)    NOT NULL PRIMARY KEY,
    content      TEXT        NOT NULL,
    recipient_id CHAR(36)    NOT NULL,
    sender_id    CHAR(36)    NOT NULL,
    state        VARCHAR(36) NOT NULL,
    sent_at      DATETIME(3) NULL,
    KEY idx_chat_private_messages_recipient_id (recipient_id),
    KEY idx_chat_private_messages_sender_id (sender_id),
    KEY idx_chat_private_messages_state (state),
    CONSTRAINT fk_chat_private_messages_recipient FOREIGN KEY (recipient_id)
        REFERENCES chat_users (id) ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT fk_chat_private_messages_sender FOREIGN KEY (sender_id)
        REFERENCES chat_users (id) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS chat_room_messages
(
    id        CHAR(36)    NOT NULL PRIMARY KEY,
    sender_id CHAR(36)    NOT NULL,
    room_id   CHAR(36)    NOT NULL,
    content   TEXT        NOT NULL,
    sent_at   DATETIME(3) NULL,
    KEY idx_chat_room_messages_sender_id (sender_id),
    KEY idx_chat_room_messages_room_id (room_id),
    CONSTRAINT fk_chat_room_messages_sender FOREIGN KEY (sender_id)
        REFERENCES chat_users (id) ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT fk_chat_room_messages_room FOREIGN KEY (room_id)
        REFERENCES chat_rooms (id) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS chat_room_message_states
(
    message_id CHAR(36)    NOT NULL,
    user_id    CHAR(36)    NOT NULL,
    state      VARCHAR(36) NOT NULL,
    UNIQUE KEY idx_room_msg_state (message_id, user_id),
    KEY idx_chat_room_message_states_message_id (message_id),
    KEY idx_chat_room_message_states_user_id (user_id),
    KEY idx_chat_room_message_states_state (state),
    CONSTRAINT fk_chat_room_message_states_message FOREIGN KEY (message_id)
        REFERENCES chat_room_messages (id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_chat_room_message_states_user FOREIGN KEY (user_id)
        REFERENCES chat_users (id) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

-- 6. 请求表
CREATE TABLE IF NOT EXISTS friendship_friend_requests
(
    id      CHAR(36)     NOT NULL PRIMARY KEY,
    `from`  CHAR(36)     NOT NULL,
    `to`    CHAR(36)     NOT NULL,
    content VARCHAR(500) NULL,
    state   VARCHAR(191) NOT NULL,
    sent_at DATETIME(3)  NOT NULL,
    KEY idx_friend_req_from (`from`),
    KEY idx_friend_req_to (`to`),
    KEY idx_friend_req_state (state),
    KEY idx_friend_req_sentat (sent_at),
    KEY idx_friend_req_from_to_state (`from`, `to`, state),
    KEY idx_friend_req_to_sentat (`to`, sent_at),
    CONSTRAINT fk_friendship_friend_requests_from FOREIGN KEY (`from`)
        REFERENCES friendship_users (id) ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT fk_friendship_friend_requests_to FOREIGN KEY (`to`)
        REFERENCES friendship_users (id) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS roomship_member_requests
(
    id           CHAR(36)     NOT NULL PRIMARY KEY,
    applicant_id CHAR(36)     NOT NULL,
    room_id      CHAR(36)     NOT NULL,
    content      VARCHAR(500) NULL,
    state        VARCHAR(191) NOT NULL,
    sent_at      DATETIME(3)  NOT NULL,
    operator_id  CHAR(36)     NULL,
    operated_at  DATETIME(3)  NULL,
    KEY idx_member_req_applicant (applicant_id),
    KEY idx_member_req_room_id (room_id),
    KEY idx_member_req_state (state),
    KEY idx_member_req_sentat (sent_at),
    KEY idx_member_req_applicant_room_state (applicant_id, room_id, state),
    KEY idx_member_req_to_sentat (room_id, sent_at),
    KEY idx_roomship_member_requests_operator_id (operator_id),
    CONSTRAINT fk_roomship_member_requests_applicant FOREIGN KEY (applicant_id)
        REFERENCES roomship_users (id) ON DELETE RESTRICT ON UPDATE CASCADE,
    CONSTRAINT fk_roomship_member_requests_room FOREIGN KEY (room_id)
        REFERENCES roomship_rooms (id) ON DELETE CASCADE ON UPDATE CASCADE,
    CONSTRAINT fk_roomship_member_requests_operator FOREIGN KEY (operator_id)
        REFERENCES roomship_users (id) ON DELETE RESTRICT ON UPDATE CASCADE
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

-- 7. 通知表
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

-- 8. 事件表
CREATE TABLE IF NOT EXISTS unpublished_events
(
    id               CHAR(36)             NOT NULL PRIMARY KEY,
    aggregate_id     CHAR(36)             NOT NULL,
    topic            VARCHAR(100)         NOT NULL,
    published        TINYINT(1) DEFAULT 0 NOT NULL,
    processing_until DATETIME(3)          NULL,
    payload          LONGBLOB             NULL,
    headers          LONGBLOB             NULL,
    created_at       DATETIME(3)          NULL,
    published_at     DATETIME(3)          NULL,
    KEY idx_unpublished_events_aggregate_id (aggregate_id),
    KEY idx_unpublished_events_topic (topic),
    KEY idx_unpublished_events_published (published),
    KEY idx_unpublished_events_processing_until (processing_until)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS dead_letters
(
    id               CHAR(36)             NOT NULL PRIMARY KEY,
    aggregate_id     CHAR(36)             NOT NULL,
    topic            VARCHAR(100)         NOT NULL,
    published        TINYINT(1) DEFAULT 0 NOT NULL,
    processing_until DATETIME(3)          NULL,
    payload          LONGBLOB             NULL,
    headers          LONGBLOB             NULL,
    created_at       DATETIME(3)          NULL,
    published_at     DATETIME(3)          NULL,
    reason           TEXT                 NOT NULL,
    failed_at        DATETIME(3)          NULL,
    KEY idx_dead_letters_aggregate_id (aggregate_id),
    KEY idx_dead_letters_topic (topic),
    KEY idx_dead_letters_published (published),
    KEY idx_dead_letters_processing_until (processing_until)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;

-- +goose Down
-- 按依赖顺序删除表
SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS chat_room_message_states;
DROP TABLE IF EXISTS chat_room_messages;
DROP TABLE IF EXISTS chat_private_messages;
DROP TABLE IF EXISTS roomship_member_requests;
DROP TABLE IF EXISTS friendship_friend_requests;
DROP TABLE IF EXISTS notification_system_messages;
DROP TABLE IF EXISTS chat_roomships;
DROP TABLE IF EXISTS profile_roomships;
DROP TABLE IF EXISTS roomship_roomships;
DROP TABLE IF EXISTS chat_friendships;
DROP TABLE IF EXISTS friendship_friendships;
DROP TABLE IF EXISTS dead_letters;
DROP TABLE IF EXISTS unpublished_events;
DROP TABLE IF EXISTS profile_room_profiles;
DROP TABLE IF EXISTS profile_user_profiles;
DROP TABLE IF EXISTS roomship_rooms;
DROP TABLE IF EXISTS chat_rooms;
DROP TABLE IF EXISTS profile_rooms;
DROP TABLE IF EXISTS authorization_users;
DROP TABLE IF EXISTS chat_users;
DROP TABLE IF EXISTS friendship_users;
DROP TABLE IF EXISTS profile_users;
DROP TABLE IF EXISTS roomship_users;

SET FOREIGN_KEY_CHECKS = 1;

