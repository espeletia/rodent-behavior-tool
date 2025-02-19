-- +goose Up
-- +goose StatementBegin
ALTER TABLE video_analysis ADD cage_id UUID;
ALTER TABLE video_analysis ADD CONSTRAINT video_analysis_cage_id_fk FOREIGN KEY (cage_id) REFERENCES rodent_cages (id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
