CREATE TABLE "users" (
    "username" varchar PRIMARY KEY,
    "role" varchar NOT NULL DEFAULT 'depositor',
    "hashed_password" varchar NOT NULL,
    "full_name" varchar NOT NULL,
    "email" varchar UNIQUE NOT NULL,
    "is_email_verified" bool NOT NULL DEFAULT false,
    "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01',
    "created_at" timestamptz NOT NULL DEFAULT (now())
);
CREATE TABLE "verify_emails" (
    "id" bigserial PRIMARY KEY,
    "username" varchar NOT NULL,
    "email" varchar NOT NULL,
    "secret_code" varchar NOT NULL,
    "is_used" bool NOT NULL DEFAULT false,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "expired_at" timestamptz NOT NULL DEFAULT (now() + interval '15 minutes')
);
CREATE TABLE "accounts" (
    "id" bigserial PRIMARY KEY,
    "owner" varchar NOT NULL,
    "balance" bigint NOT NULL,
    "currency" varchar NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);
CREATE TABLE "entries" (
    "id" bigserial PRIMARY KEY,
    "account_id" bigint NOT NULL,
    "amount" bigint NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);
CREATE TABLE "transfers" (
    "id" bigserial PRIMARY KEY,
    "from_account_id" bigint NOT NULL,
    "to_account_id" bigint NOT NULL,
    "amount" bigint NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);
CREATE TABLE "sessions" (
    "id" uuid PRIMARY KEY,
    "username" varchar NOT NULL,
    "refresh_token" varchar NOT NULL,
    "user_agent" varchar NOT NULL,
    "client_ip" varchar NOT NULL,
    "is_blocked" boolean NOT NULL DEFAULT false,
    "expires_at" timestamptz NOT NULL,
    "created_at" timestamptz NOT NULL DEFAULT (now())
);
CREATE TABLE "charges"(
    "id": "ch_3MmlLrLkdIwHu7ix0snN0B15",
    "object": "charge",
    "amount": 1099,
    "amount_captured": 1099,
    "amount_refunded": 0,
    "application": null,
    "application_fee": null,
    "application_fee_amount": null,
    "balance_transaction": "txn_3MmlLrLkdIwHu7ix0uke3Ezy",
    "billing_details": varchar NOT NULL,
    "calculated_statement_descriptor": "Stripe",
    "captured": true,
    "created": 1679090539,
    "currency": "usd",
    "customer": null,
    "description": null,
    "disputed": false,
    "failure_balance_transaction": null,
    "failure_code": null,
    "failure_message": null,
    "fraud_details": { },
    "invoice": null,
    "livemode": false,
    "metadata": { },
    "on_behalf_of": null,
    "outcome": { "network_status": "approved_by_network",
    "reason": null,
    "risk_level": "normal",
    "risk_score": 32,
    "seller_message": "Payment complete.",
    "type": "authorized" },
    "paid": true,
    "payment_intent": null,
    "payment_method": "card_1MmlLrLkdIwHu7ixIJwEWSNR",
    "payment_method_details": "checks": varchar NOT NULL,
    "receipt_email": null,
    "receipt_number": null,
    "receipt_url": "https://pay.stripe.com/receipts/payment/CAcaFwoVYWNjdF8xTTJKVGtMa2RJd0h1N2l4KOvG06AGMgZfBXyr1aw6LBa9vaaSRWU96d8qBwz9z2J_CObiV_H2-e8RezSK_sw0KISesp4czsOUlVKY",
    "refunded": false,
    "review": null,
    "shipping": null,
    "source_transfer": null,
    "statement_descriptor": null,
    "statement_descriptor_suffix": null,
    "status": "succeeded",
    "transfer_data": null,
    "transfer_group": null
);
CREATE TABLE "billing_details"(
    "address" varchar NOT NULL "email": null,
    "name": null,
    "phone": null
);
CREATE TABLE "address"(
    "city": null,
    "country": null,
    "line1": null,
    "line2": null,
    "postal_code": null,
    "state": null
);
CREATE TABLE "payment_method"("card");
CREATE TABLE "card"(
    "id": "card_1MvoiELkdIwHu7ixOeFGbN9D",
    "object": "card",
    "address_city": null,
    "address_country": null,
    "address_line1": null,
    "address_line1_check": null,
    "address_line2": null,
    "address_state": null,
    "address_zip": null,
    "address_zip_check": null,
    "brand": "Visa",
    "country": "US",
    "customer": "cus_NhD8HD2bY8dP3V",
    "cvc_check": null,
    "dynamic_last4": null,
    "exp_month": 4,
    "exp_year": 2024,
    "fingerprint": "mToisGZ01V71BCos",
    "funding": "credit",
    "last4": "4242",
    "metadata": { },
    "name": null,
    "tokenization_method": null,
    "wallet": null
);
CREATE TABLE "payout"(
    "id": "po_1OaFDbEcg9tTZuTgNYmX0PKB",
    "object": "payout",
    "amount": 1100,
    "arrival_date": 1680652800,
    "automatic": false,
    "balance_transaction": "txn_1OaFDcEcg9tTZuTgYMR25tSe",
    "created": 1680648691,
    "currency": "usd",
    "description": null,
    "destination": "ba_1MtIhL2eZvKYlo2CAElKwKu2",
    "failure_balance_transaction": null,
    "failure_code": null,
    "failure_message": null,
    "livemode": false,
    "metadata": { },
    "method": "standard",
    "original_payout": null,
    "reconciliation_status": "not_applicable",
    "reversed_by": null,
    "source_type": "card",
    "statement_descriptor": null,
    "status": "pending",
    "type": "bank_account"
);
CREATE TABLE "refund"(
    { "id": "re_1Nispe2eZvKYlo2Cd31jOCgZ",
    "object": "refund",
    "amount": 1000,
    "balance_transaction": "txn_1Nispe2eZvKYlo2CYezqFhEx",
    "charge": "ch_1NirD82eZvKYlo2CIvbtLWuY",
    "created": 1692942318,
    "currency": "usd",
    "destination_details": { "card": { "reference": "123456789012",
    "reference_status": "available",
    "reference_type": "acquirer_reference_number",
    "type": "refund" },
    "type": "card" },
    "metadata": { },
    "payment_intent": "pi_1GszsK2eZvKYlo2CfhZyoZLp",
    "reason": null,
    "receipt_number": null,
    "source_transfer_reversal": null,
    "status": "succeeded",
    "transfer_reversal": null }
);
CREATE TABLE "payment_method"(
    "id": "pm_1Q0PsIJvEtkwdCNYMSaVuRz6",
    "object": "payment_method",
    "allow_redisplay": "unspecified",
    "billing_details": { "address": { "city": null,
    "country": null,
    "line1": null,
    "line2": null,
    "postal_code": null,
    "state": null },
    "email": null,
    "name": "John Doe",
    "phone": null },
    "created": 1726673582,
    "customer": null,
    "livemode": false,
    "metadata": { },
    "type": "us_bank_account",
    "us_bank_account": { "account_holder_type": "individual",
    "account_type": "checking",
    "bank_name": "STRIPE TEST BANK",
    "financial_connections_account": null,
    "fingerprint": "LstWJFsCK7P349Bg",
    "last4": "6789",
    "networks": { "preferred": "ach",
    "supported": [
                "ach"
            ] },
    "routing_number": "110000000",
    "status_details": { } }
);
CREATE TABLE "bank_account"(
    "id": "ba_1MvoIJ2eZvKYlo2CO9f0MabO",
    "object": "bank_account",
    "account_holder_name": "Jane Austen",
    "account_holder_type": "company",
    "account_type": null,
    "bank_name": "STRIPE TEST BANK",
    "country": "US",
    "currency": "usd",
    "customer": "cus_9s6XI9OFIdpjIg",
    "fingerprint": "1JWtPxqbdX5Gamtc",
    "last4": "6789",
    "metadata": { },
    "routing_number": "110000000",
    "status": "new"
);
CREATE INDEX ON "accounts" ("owner");
CREATE UNIQUE INDEX ON "accounts" ("owner", "currency");
CREATE INDEX ON "entries" ("account_id");
CREATE INDEX ON "transfers" ("from_account_id");
CREATE INDEX ON "transfers" ("to_account_id");
CREATE INDEX ON "transfers" ("from_account_id", "to_account_id");
COMMENT ON COLUMN "entries"."amount" IS 'can be negative or positive';
COMMENT ON COLUMN "transfers"."amount" IS 'must be positive';
ALTER TABLE "verify_emails"
ADD FOREIGN KEY ("username") REFERENCES "users" ("username");
ALTER TABLE "accounts"
ADD FOREIGN KEY ("owner") REFERENCES "users" ("username");
ALTER TABLE "entries"
ADD FOREIGN KEY ("account_id") REFERENCES "accounts" ("id");
ALTER TABLE "transfers"
ADD FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id");
ALTER TABLE "transfers"
ADD FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id");
ALTER TABLE "sessions"
ADD FOREIGN KEY ("username") REFERENCES "users" ("username");