CREATE TABLE FOLLOWERS (
    user_id UUID NOT NULL,
    follower_id UUID NOT NULL,
    PRIMARY KEY(user_id, follower_id),
    FOREIGN KEY(user_id)
        REFERENCES USERS(id)
        ON DELETE CASCADE,
    FOREIGN KEY(follower_id)
        REFERENCES USERS(id)
        ON DELETE CASCADE
);
