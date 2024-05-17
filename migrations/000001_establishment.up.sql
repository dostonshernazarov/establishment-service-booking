CREATE TABLE "location_table"(
    "location_id" UUID PRIMARY KEY NOT NULL,
    "establishment_id" UUID NOT NULL,
    "address" VARCHAR(255) DEFAULT '',
    "latitude" FLOAT DEFAULT 0,
    "longitude" FLOAT DEFAULT 0,
    "country" VARCHAR(255) DEFAULT '',
    "city" VARCHAR(255) DEFAULT '',
    "state_province" VARCHAR(255) DEFAULT '',
    "category"  VARCHAR(255) DEFAULT '',
    "created_at" TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP(0)
);

CREATE TABLE "room_table"(
    "room_id" UUID PRIMARY KEY NOT NULL,
    "hotel_id" UUID NOT NULL,
    "price" FLOAT DEFAULT 0,
    "description" TEXT DEFAULT '',
    "number_of_rooms" BIGINT DEFAULT 0,
    "holidays" VARCHAR(255) DEFAULT '',
    "free_days" VARCHAR(255) DEFAULT '',
    "discount" FLOAT DEFAULT 0,
    "created_at" TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP(0)
);

CREATE TABLE "favourite_table"(
    "favourite_id" UUID PRIMARY KEY NOT NULL,
    "establishment_id" UUID NOT NULL,
    "user_id" UUID NOT NULL,
    "created_at" TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP(0)
);

CREATE TABLE "review_table"(
    "review_id" UUID PRIMARY KEY NOT NULL,
    "establishment_id" UUID NOT NULL,
    "user_id" UUID NOT NULL,
    "rating" FLOAT DEFAULT 0,
    "comment" TEXT DEFAULT '',
    "created_at" TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP(0)
);

CREATE TABLE "image_table"(
    "image_id" UUID PRIMARY KEY NOT NULL,
    "establishment_id" UUID NOT NULL,
    "image_url" VARCHAR(255) DEFAULT '',
    "category"  VARCHAR(255) DEFAULT '',
    "created_at" TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP(0)
);

CREATE TABLE "restaurant_table"(
    "restaurant_id" UUID PRIMARY KEY NOT NULL,
    "owner_id" UUID NOT NULL,
    "restaurant_name" VARCHAR(255) DEFAULT '',
    "description" TEXT DEFAULT '',
    "rating" FLOAT DEFAULT 0,
    "opening_hours" VARCHAR(255) DEFAULT '',
    "contact_number" VARCHAR(255) DEFAULT '',
    "licence_url" VARCHAR(255) DEFAULT '',
    "website_url" VARCHAR(255) DEFAULT '',
    "created_at" TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP(0)
);

CREATE TABLE "attraction_table"(
    "attraction_id" UUID PRIMARY KEY NOT NULL,
    "owner_id" UUID NOT NULL,
    "attraction_name" VARCHAR(255) DEFAULT '',
    "description" VARCHAR(255) DEFAULT '',
    "rating" FLOAT DEFAULT 0,
    "contact_number" VARCHAR(255) DEFAULT '',
    "licence_url" VARCHAR(255) DEFAULT '',
    "website_url" VARCHAR(255) DEFAULT '',
    "created_at" TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP(0)
);

CREATE TABLE "hotel_table"(
    "hotel_id" UUID PRIMARY KEY NOT NULL,
    "owner_id" UUID NOT NULL,
    "hotel_name" VARCHAR(255) DEFAULT '',
    "description" TEXT DEFAULT '',
    "rating" FLOAT DEFAULT 0,
    "contact_number" VARCHAR(255) DEFAULT '',
    "licence_url" VARCHAR(255) DEFAULT '',
    "website_url" VARCHAR(255) DEFAULT '',
    "created_at" TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP(0) DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP(0)
);