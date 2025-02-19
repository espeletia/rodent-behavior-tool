-- +goose Up
-- +goose StatementBegin
ALTER TABLE video_analysis ALTER COLUMN owner_id DROP NOT NULL;
ALTER TABLE video_analysis ADD CONSTRAINT video_analysis_cage_or_user_chk CHECK (cage_id IS NOT NULL OR owner_id IS NOT NULL);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
