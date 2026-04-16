-- +goose Up
-- Match outbox lease query: WHERE processing_until <= ? ... FOR UPDATE SKIP LOCKED LIMIT ?
ALTER TABLE unpublished_events
    DROP INDEX idx_unpublished_events_processing_until,
    ADD INDEX idx_unpublished_events_processing_until_id (processing_until, id);

-- +goose Down
ALTER TABLE unpublished_events
    DROP INDEX idx_unpublished_events_processing_until_id,
    ADD INDEX idx_unpublished_events_processing_until (processing_until);

