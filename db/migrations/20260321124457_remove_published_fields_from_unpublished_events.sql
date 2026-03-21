-- +goose Up
-- 首先需要删除索引，因为字段被索引引用
DROP INDEX idx_unpublished_events_published ON unpublished_events;

-- 删除 published 和 published_at 字段
ALTER TABLE unpublished_events
    DROP COLUMN published,
    DROP COLUMN published_at;

-- +goose Down
-- 重新添加字段（按照原来的结构）
ALTER TABLE unpublished_events
    ADD COLUMN published    TINYINT(1) DEFAULT 0 NOT NULL,
    ADD COLUMN published_at DATETIME(3)          NULL;

-- 重新创建索引
CREATE INDEX idx_unpublished_events_published ON unpublished_events (published);