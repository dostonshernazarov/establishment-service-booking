CREATE TABLE location_table(
    location_id UUID NOT NULL,
    address VARCHAR(255) DEFAULT '',
    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,
    country VARCHAR(255) DEFAULT '',
    city VARCHAR(255) DEFAULT '',
    state_province VARCHAR(255) DEFAULT '',
    created_at TIMESTAMP(0) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP(0) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP(0) WITH TIME ZONE,
    PRIMARY KEY(location_id)
);

CREATE TABLE room_table(
    room_id UUID NOT NULL,
    hotel_id UUID NOT NULL,
    price DOUBLE PRECISION DEFAULT 0,
    description TEXT DEFAULT '',
    number_of_rooms BIGINT DEFAULT 0,
    holidays VARCHAR(255) DEFAULT '',
    free_days VARCHAR(255) DEFAULT '',
    discount DOUBLE PRECISION DEFAULT 0,
    created_at TIMESTAMP(0) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP(0) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP(0) WITH TIME ZONE,
    PRIMARY KEY(room_id)
);

CREATE TABLE favourite_table(
    favourite_id UUID NOT NULL,
    establishment_id UUID NOT NULL,
    user_id UUID NOT NULL,
    PRIMARY KEY(favourite_id)
);

CREATE TABLE review_table(
    review_id UUID NOT NULL,
    establishment_id UUID NOT NULL,
    user_id UUID NOT NULL,
    rating DOUBLE PRECISION DEFAULT 0,
    comment TEXT DEFAULT '',
    PRIMARY KEY(review_id)
);

CREATE TABLE image_table(
    image_id UUID NOT NULL,
    establishment_id UUID NOT NULL,
    image_url VARCHAR(255) NOT NULL,
    created_at TIMESTAMP(0) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP(0) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP(0) WITH TIME ZONE,
    PRIMARY KEY(image_id)
);

CREATE TABLE restaurant_table(
    restaurant_id UUID NOT NULL,
    owner_id UUID NOT NULL,
    restaurant_name VARCHAR(255) DEFAULT '',
    description TEXT DEFAULT '',
    rating DOUBLE PRECISION DEFAULT 0,
    opening_hours VARCHAR(255) DEFAULT '',
    contact_number VARCHAR(255) DEFAULT '',
    licence_url VARCHAR(255) DEFAULT '',
    website_url VARCHAR(255) DEFAULT '',
    created_at TIMESTAMP(0) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP(0) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP(0) WITH TIME ZONE,
    location_id UUID NOT NULL,
    PRIMARY KEY(restaurant_id),
    CONSTRAINT restaurant_table_location_id_foreign FOREIGN KEY(location_id) REFERENCES location_table(location_id)
);

CREATE TABLE attraction_table(
    attraction_id UUID NOT NULL,
    owner_id UUID NOT NULL,
    attraction_name VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    rating DOUBLE PRECISION DEFAULT 0,
    contact_number VARCHAR(255) DEFAULT '',
    licence_url VARCHAR(255) DEFAULT '',
    website_url VARCHAR(255) DEFAULT '',
    created_at TIMESTAMP(0) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP(0) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP(0) WITH TIME ZONE,
    location_id UUID NOT NULL,
    PRIMARY KEY(attraction_id),
    CONSTRAINT attraction_table_location_id_foreign FOREIGN KEY(location_id) REFERENCES location_table(location_id)
);

CREATE TABLE hotel_table(
    hotel_id UUID NOT NULL,
    owner_id BIGINT NOT NULL,
    hotel_name VARCHAR(255) DEFAULT '',
    description TEXT DEFAULT '',
    rating DOUBLE PRECISION DEFAULT 0,
    contact_number VARCHAR(255) DEFAULT '',
    licence_url VARCHAR(255) DEFAULT '',
    website_url VARCHAR(255) DEFAULT '',
    created_at TIMESTAMP(0) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP(0) WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP(0) WITH TIME ZONE,
    location_id UUID NOT NULL,
    PRIMARY KEY(hotel_id),
    CONSTRAINT hotel_table_location_id_foreign FOREIGN KEY(location_id) REFERENCES location_table(location_id)
);

ALTER TABLE image_table ADD CONSTRAINT image_table_establishment_id_foreign FOREIGN KEY(establishment_id) REFERENCES restaurant_table(restaurant_id);

ALTER TABLE image_table ADD CONSTRAINT image_table_establishment_id_foreign2 FOREIGN KEY(establishment_id) REFERENCES attraction_table(attraction_id);

ALTER TABLE review_table ADD CONSTRAINT review_table_establishment_id_foreign FOREIGN KEY(establishment_id) REFERENCES hotel_table(hotel_id);

ALTER TABLE image_table ADD CONSTRAINT image_table_establishment_id_foreign3 FOREIGN KEY(establishment_id) REFERENCES hotel_table(hotel_id);

ALTER TABLE favourite_table ADD CONSTRAINT favourite_table_establishment_id_foreign FOREIGN KEY(establishment_id) REFERENCES hotel_table(hotel_id);

ALTER TABLE favourite_table ADD CONSTRAINT favourite_table_establishment_id_foreign2 FOREIGN KEY(establishment_id) REFERENCES attraction_table(attraction_id);

ALTER TABLE favourite_table ADD CONSTRAINT favourite_table_establishment_id_foreign3 FOREIGN KEY(establishment_id) REFERENCES restaurant_table(restaurant_id);

ALTER TABLE room_table ADD CONSTRAINT room_table_room_id_foreign FOREIGN KEY(room_id) REFERENCES hotel_table(hotel_id);

ALTER TABLE review_table ADD CONSTRAINT review_table_establishment_id_foreign2 FOREIGN KEY(establishment_id) REFERENCES restaurant_table(restaurant_id);

ALTER TABLE review_table ADD CONSTRAINT review_table_establishment_id_foreign3 FOREIGN KEY(establishment_id) REFERENCES attraction_table(attraction_id);

ALTER TABLE attraction_table ADD CONSTRAINT attraction_table_location_id_foreign FOREIGN KEY(location_id) REFERENCES location_table(location_id);














































ALTER TABLE "review_table" DROP CONSTRAINT IF EXISTS "review_table_establishment_id_foreign";
ALTER TABLE "review_table" DROP CONSTRAINT IF EXISTS "review_table_establishment_id_foreign";

ALTER TABLE "attraction_table" DROP CONSTRAINT IF EXISTS "attraction_table_location_id_foreign";
ALTER TABLE "review_table" DROP CONSTRAINT IF EXISTS "review_table_establishment_id_foreign";

ALTER TABLE "room_table" DROP CONSTRAINT IF EXISTS "room_table_room_id_foreign";
ALTER TABLE "favourite_table" DROP CONSTRAINT IF EXISTS "favourite_table_establishment_id_foreign";

ALTER TABLE "favourite_table" DROP CONSTRAINT IF EXISTS "favourite_table_establishment_id_foreign";
ALTER TABLE "restaurant_table" DROP CONSTRAINT IF EXISTS "restaurant_table_location_id_foreign";

ALTER TABLE "favourite_table" DROP CONSTRAINT IF EXISTS "favourite_table_establishment_id_foreign";
ALTER TABLE "image_table" DROP CONSTRAINT IF EXISTS "image_table_establishment_id_foreign";

ALTER TABLE "image_table" DROP CONSTRAINT IF EXISTS "image_table_establishment_id_foreign";
ALTER TABLE "review_table" DROP CONSTRAINT IF EXISTS "review_table_establishment_id_foreign";

ALTER TABLE "hotel_table" DROP CONSTRAINT IF EXISTS "hotel_table_location_id_foreign";
ALTER TABLE "image_table" DROP CONSTRAINT IF EXISTS "image_table_establishment_id_foreign";

ALTER TABLE "review_table" DROP CONSTRAINT IF EXISTS "review_table_establishment_id_foreign";
ALTER TABLE "image_table" DROP CONSTRAINT IF EXISTS "image_table_establishment_id_foreign";

ALTER TABLE "hotel_table" DROP CONSTRAINT IF EXISTS "hotel_table_location_id_foreign";
ALTER TABLE "attraction_table" DROP CONSTRAINT IF EXISTS "attraction_table_location_id_foreign";

DROP TABLE IF EXISTS "attraction_table";
DROP TABLE IF EXISTS "hotel_table";
DROP TABLE IF EXISTS "restaurant_table";
DROP TABLE IF EXISTS "favourite_table";
DROP TABLE IF EXISTS "review_table";
DROP TABLE IF EXISTS "image_table";
DROP TABLE IF EXISTS "room_table";
DROP TABLE IF EXISTS "location_table";