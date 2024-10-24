CREATE TABLE IF NOT EXISTS warnings (
    warning_id SERIAL PRIMARY KEY, 
    ruleId text NOT NULL,
    uri text NOT NULL,
    startLine int NOT NULL CHECK (xseverity >= 0),
    xseverity int NOT NULL CHECK (xseverity >= 0 AND xseverity <= 2)
);
