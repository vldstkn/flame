CREATE TABLE swipes(
    user_id1 BIGINT NOT NULL,
    user_id2 BIGINT NOT NULL,
    user_is_liked1 bool,
    user_is_liked2 bool,
    PRIMARY KEY (user_id1, user_id2)
);
CREATE INDEX idx_swipes_user_id1 ON swipes(user_id1);
CREATE INDEX idx_swipes_user_id2 ON swipes(user_id2);