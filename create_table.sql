CREATE Table "user"(
    "id" varchar(255) PRIMARY KEY,
    "email" varchar(255)
    "type" varchar(255)
);

CREATE TYPE "status" as ENUM(
    'New', 
    'Inprogress', 
    'Closed'
);

CREATE Table "issue"(
    "id" varchar(255) PRIMARY KEY,
    "title" varchar(255),
    "description" varchar(255),
    "reported_by" varchar(255) REFERENCES "user"("id"),
    "resolved_by" varchar(255) REFERENCES "user"("id"),
    "status" Status DEFAULT 'New',
    "resolved_at" timestamp,
    "created_by" varchar(255) REFERENCES "user"("id"),
    "created_at" timestamp,
    "updated_at" timestamp,
    "belongs_to" varchar(255) REFERENCES "user"("id")
);