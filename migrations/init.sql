CREATE TABLE countries (
                           iso2_code CHAR(2) PRIMARY KEY,
                           name VARCHAR(255) NOT NULL
);

CREATE TABLE swift_codes (
                             swift_code VARCHAR(11) PRIMARY KEY,
                             address TEXT NOT NULL,
                             bank_name TEXT NOT NULL,
                             country_iso2 CHAR(2) NOT NULL REFERENCES countries(iso2_code),
                             is_headquarter BOOLEAN NOT NULL
);

CREATE INDEX idx_swift_country ON swift_codes (country_iso2);
