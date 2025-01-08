-- +goose Up
-- +goose StatementBegin
CREATE SEQUENCE cage_messages_id_seq
START WITH 30000
INCREMENT BY 1
NO MINVALUE
NO MAXVALUE
CACHE 1;

CREATE TABLE IF NOT EXISTS cage_messages (
    id INT NOT NULL DEFAULT nextval('cage_messages_id_seq'),
    cage_id UUID NOT NULL,
    revision INT NOT NULL,
    water INT NOT NULL,
    food INT NOT NULL,
    light INT NOT NULL,
    temp INT NOT NULL,
    humidity INT NOT NULL,
    video_url TEXT,
    video_id UUID,
    time_sent TIMESTAMP NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT cage_messages_pkey PRIMARY KEY (id),
    CONSTRAINT rodent_cages_cage_id_fk FOREIGN KEY (cage_id) REFERENCES rodent_cages (id) ON DELETE CASCADE,
    CONSTRAINT rodent_cages_video_id_fk FOREIGN KEY (video_id) REFERENCES video_analysis (id) ON DELETE CASCADE
);

CREATE TRIGGER update_cage_messages_updated_at
BEFORE UPDATE ON cage_messages
FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

ALTER TABLE rodent_cages ADD secret_token TEXT NOT NULL; 

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
