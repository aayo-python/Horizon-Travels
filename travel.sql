CREATE TABLE "hotels" (
  "id" varchar NOT NULL,
  "title" varchar NOT NULL,
  "primaryInfo" varchar NOT NULL,
  "secondaryInfo" varchar NOT NULL,
  "bubbleRating" varchar NOT NULL,
  "isSponsored" varchar NOT NULL,
  "accentedLabel" bool NOT NULL DEFAULT false,
  "provider" varchar NOT NULL,
  "priceForDisplay" varchar NOT NULL,
  "strikethroughPrice" varchar NOT NULL,
  "priceDetails" varchar NOT NULL,
  "priceSummary" varchar NOT NULL,
  "cardPhotos" varchar NOT NULL
);

CREATE TABLE "bubbleratings" (
  "id" BIGSERIAL NOT NULL,
  "count" varchar NOT NULL,
  "rating" INT NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "cardphotos" (
  "id" BIGSERIAL NOT NULL,
  "__typename" varchar NOT NULL,
  "maxHeight" INT NOT NULL,
  "maxWidth" INT NOT NULL,
  "urlTemplate" varchar NOT NULL
);

CREATE TABLE "restaurants" (
  "heroImgUrl" varchar NOT NULL,
  "heroImgRawHeight" INT NOT NULL,
  "heroImgRawWidth" INT NOT NULL,
  "squareImgUrl" varchar NOT NULL,
  "squareImgRawLength" INT NOT NULL,
  "locationId" INT NOT NULL,
  "name" varchar NOT NULL,
  "averageRating" varchar NOT NULL,
  "userReviewCount" varchar NOT NULL,
  "currentOpenStatusCategory" varchar NOT NULL,
  "currentOpenStatusText" varchar NOT NULL,
  "establishmentTypeAndCuisineTags" varchar NOT NULL,
  "priceTag" varchar NOT NULL,
  "offers" varchar NOT NULL,
  "hasMenu" boolean NOT NULL,
  "menuUrl" varchar NOT NULL,
  "isDifferentGeo" boolean NOT NULL,
  "parentGeoName" varchar NOT NULL,
  "distanceTo" varchar NOT NULL,
  "reviewSnippets" varchar NOT NULL,
  "isLocalChefItem" boolean NOT NULL,
  "isPremium" BOOLEAN NOT NULL,
  "isStoryboardPublished" boolean NOT NULL
);

CREATE TABLE "establishmenttypeandcuisinetags" (
  "id" varchar NOT NULL,
  "type" varchar NOT NULL
);

CREATE TABLE "reviewsnippets" (
  "id" varchar NOT NULL,
  "reviewText" varchar NOT NULL,
  "reviewUrl" varchar NOT NULL
);

CREATE TABLE "offers" (
  "id" BIGSERIAL NOT NULL,
  "slot1Offer" varchar NOT NULL,
  "slot2Offer" varchar NOT NULL
);

CREATE TABLE "vacationrentals" (
  "geoId" INT NOT NULL,
  "locationId" INT NOT NULL,
  "localizedName" VARCHAR NOT NULL,
  "localizedAdditionalNames" VARCHAR NOT NULL,
  "locationV2" VARCHAR NOT NULL,
  "placeType" VARCHAR NOT NULL,
  "latitude" varchar NOT NULL,
  "longitude" varchar NOT NULL,
  "isGeo" boolean NOT NULL,
  "thumbnail" VARCHAR NOT NULL
);

CREATE TABLE "thumbnails" (
  "id" varchar NOT NULL,
  "maxWidth" INT NOT NULL,
  "maxHeight" INT NOT NULL,
  "urlTemplate" VARCHAR NOT NULL
);

CREATE TABLE "cars" (
  "geoId" Int NOT NULL,
  "locationId" Int NOT NULL,
  "localizedName" Varchar NOT NULL,
  "localizedAdditionalNames" Varchar NOT NULL,
  "locationV2" Varchar NOT NULL,
  "placeType" Varchar NOT NULL,
  "latitude" Int NOT NULL,
  "longitude" Int NOT NULL,
  "isGeo" boolean NOT NULL,
  "thumbnail" varchar NOT NULL
);

ALTER TABLE "hotels" ADD FOREIGN KEY ("bubbleRating") REFERENCES "bubbleratings" ("id");

ALTER TABLE "hotels" ADD FOREIGN KEY ("cardPhotos") REFERENCES "cardphotos" ("id");

ALTER TABLE "restaurants" ADD FOREIGN KEY ("establishmentTypeAndCuisineTags") REFERENCES "establishmenttypeandcuisinetags" ("id");

ALTER TABLE "restaurants" ADD FOREIGN KEY ("reviewSnippets") REFERENCES "reviewsnippets" ("id");

ALTER TABLE "restaurants" ADD FOREIGN KEY ("offers") REFERENCES "offers" ("id");

ALTER TABLE "vacationrentals" ADD FOREIGN KEY ("thumbnail") REFERENCES "thumbnails" ("id");
