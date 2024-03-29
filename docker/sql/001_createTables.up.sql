CREATE TABLE prj_user (
    userid CHARACTER VARYING(32) NOT NULL UNIQUE,
    nameuser CHARACTER VARYING(300) NOT NULL,
    chatid CHARACTER VARYING(32) NOT NULL UNIQUE,

    CONSTRAINT pk_prj_user PRIMARY KEY (userid)
);

CREATE TABLE prj_link(
    linkid CHARACTER VARYING(32) NOT NULL UNIQUE,
    link CHARACTER VARYING(300) NOT NULL,
    
    CONSTRAINT pk_prj_link PRIMARY KEY (linkid)
);

CREATE TABLE ref_link_user(
    refid CHARACTER VARYING(32) NOT NULL UNIQUE,
    userid CHARACTER VARYING(32) NOT NULL,
    linkid CHARACTER VARYING(32) NOT NULL,
    
    CONSTRAINT pk_ref_link_user PRIMARY KEY (refid),
    FOREIGN KEY (userid) REFERENCES prj_user (userid) ON DELETE CASCADE
);
