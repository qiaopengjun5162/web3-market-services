DO
$$
BEGIN
        IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'uint256') THEN
CREATE DOMAIN UINT256 AS NUMERIC
    CHECK (VALUE >= 0 AND VALUE < POWER(CAST(2 AS NUMERIC), CAST(256 AS NUMERIC)) AND SCALE(VALUE) = 0);
ELSE
ALTER DOMAIN UINT256 DROP CONSTRAINT uint256_check;
ALTER DOMAIN UINT256 ADD
    CHECK (VALUE >= 0 AND VALUE < POWER(CAST(2 AS NUMERIC), CAST(256 AS NUMERIC)) AND SCALE(VALUE) = 0);
END IF;
END
$$;

CREATE TABLE IF NOT EXISTS market_price(
                                           guid            VARCHAR PRIMARY KEY,
                                           asset_name      VARCHAR UNIQUE NOT NULL,
                                           price_usdt      VARCHAR        NOT NULL,
                                           volume          VARCHAR        NOT NULL,
                                           rate            VARCHAR        NOT NULL,
                                           timestamp       INTEGER        NOT NULL CHECK (timestamp > 0)
    );
CREATE INDEX IF NOT EXISTS market_asset_name ON market_price (asset_name);


CREATE TABLE IF NOT EXISTS official_coin_rate(
                                                 guid            VARCHAR PRIMARY KEY,
                                                 asset_name      VARCHAR UNIQUE NOT NULL,
                                                 base_asset      VARCHAR        NOT NULL,
                                                 price           VARCHAR        NOT NULL,
                                                 timestamp       INTEGER        NOT NULL CHECK (timestamp > 0)
    );
CREATE INDEX IF NOT EXISTS official_coin_rate_asset_name ON official_coin_rate(asset_name);
