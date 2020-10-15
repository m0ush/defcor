CREATE INDEX ON prices (secid, date DESC);
ALTER TABLE stocks ADD COLUMN date_inactive date;