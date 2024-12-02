-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS video_analysis (
    id UUID NOT NULL,
    media_id UUID NOT NULL,
    owner_id UUID NOT NULL,
    description TEXT,
    name TEXT NOT NULL,
    analysed_video UUID,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT video_pkey PRIMARY KEY (id),
    CONSTRAINT video_media_id_fk FOREIGN KEY (media_id) REFERENCES media (id) ON DELETE CASCADE,
    CONSTRAINT video_owner_id_fk FOREIGN KEY (owner_id) REFERENCES users (id) ON DELETE CASCADE,
    CONSTRAINT video_analysed_video_id_fk FOREIGN KEY (analysed_video) REFERENCES media (id) ON DELETE SET NULL
);

CREATE TRIGGER update_video_analysis_updated_at
BEFORE UPDATE ON video_analysis
FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
