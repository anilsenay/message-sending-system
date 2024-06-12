CREATE TYPE message_status as enum ('unsent', 'processing', 'sent');

CREATE TABLE message (
  id BIGSERIAL CONSTRAINT message_pk PRIMARY KEY,
  content VARCHAR(1000) NOT NULL,
  recipient_phone_number VARCHAR(20) NOT NULL,
  status message_status NOT NULL DEFAULT 'unsent' :: message_status,
  created_at TIMESTAMP DEFAULT NOW() NOT NULL,
  sent_at TIMESTAMP
);