CREATE TABLE IF NOT EXISTS invenory_movements (
  "movement_id" BigSerial PRIMARY KEY,
  "product_id" int,
  "change" int NOT NULL,
  "reason" text,
  "moved_by" int,
  "timestamp" timestamp DEFAULT (now())
);


ALTER TABLE "invenory_movements" ADD FOREIGN KEY ("product_id") REFERENCES "products" ("product_id");
ALTER TABLE "invenory_movements" ADD FOREIGN KEY ("moved_by") REFERENCES "users" ("user_id");