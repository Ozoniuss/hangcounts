CREATE TABLE individuals (
    id          SERIAL PRIMARY KEY,
    name        TEXT,
    email       TEXT UNIQUE NOT NULL,
    username    TEXT UNIQUE NOT NULL,

    created_at  TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMPTZ
);

-- used to query individual data when they log in
CREATE INDEX idx_individuals_username ON individuals(username);
CREATE INDEX idx_individuals_email ON individuals(email);

CREATE TABLE hangouts (
    id               BIGSERIAL PRIMARY KEY,
    public_id        UUID DEFAULT gen_random_uuid() UNIQUE NOT NULL,
    location         TEXT NOT NULL,
    description      TEXT,
    duration_minutes INT NOT NULL,
    date             TIMESTAMPTZ NOT NULL,
    created_by       INT NOT NULL,

    created_at  TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    deleted_at  TIMESTAMPTZ,

    CONSTRAINT fk_hangout_creator FOREIGN KEY (created_by) REFERENCES individuals (id) ON DELETE SET NULL
);

CREATE TABLE hangout_individuals (
    id            BIGSERIAL PRIMARY KEY,
    hangout_id    BIGINT NOT NULL,
    individual_id INT NOT NULL,

    created_at  TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP, -- this mirrors the one from hangouts but is used here for optimisations
    deleted_at  TIMESTAMPTZ,


    CONSTRAINT fk_hangout FOREIGN KEY (hangout_id) REFERENCES hangouts (id) ON DELETE CASCADE, -- if a hangout is permanently deleted, these make no sense
    CONSTRAINT fk_individual FOREIGN KEY (individual_id) REFERENCES individuals (id) -- for now users will only be soft deleted so restrict their deletion if they have hangouts
);

-- support showing the hangouts for a user on his page, and showing all the users of a hangout
CREATE INDEX idx_hangout_individuals_hangout ON hangout_individuals(hangout_id);
CREATE INDEX idx_hangout_individuals_individual ON hangout_individuals(individual_id);

-- Optimize fetching a user's hangouts in descending order to display them on their page
CREATE INDEX idx_hangout_individuals_created_at ON hangout_individuals(created_at DESC NULLS LAST);
