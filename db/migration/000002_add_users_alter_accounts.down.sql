
ALTER TABLE IF EXISTS ACCOUNTS DROP CONSTRAINT IF EXISTS "owner_currency_key";
ALTER TABLE IF EXISTS ACCOUNTS DROP CONSTRAINT IF EXISTS "accounts_owner_fkey";

drop TABLE IF EXISTS  USERS;
