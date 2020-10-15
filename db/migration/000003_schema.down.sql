ALTER TABLE stocks ADD COLUMN active boolean;
ALTER TABLE splits ADD COLUMN ratio numeric(10, 6), ADD COLUMN description varchar(120);
ALTER TABLE dividends ADD COLUMN description varchar(120);
