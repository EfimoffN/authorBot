CREATE TABLE prj_user (
    userid CHARACTER VARYING(32) NOT NULL UNIQUE,
    nameuser CHARACTER VARYING(300) NOT NULL,
    chatid CHARACTER VARYING(32) NOT NULL UNIQUE,

    CONSTRAINT pk_prj_user PRIMARY KEY (userid)
)
WITH(
	OIDS=FALSE
);

ALTER TABLE prj_user OWNER TO postgres;
COMMENT ON TABLE prj_user IS 'Таблица пользователей';
COMMENT ON TABLE prj_user.userid IS 'Идентификатор пользователя';
COMMENT ON TABLE prj_user.nameuser IS 'Имя пользовтаеля в tlg';
COMMENT ON TABLE prj_user.chatid IS 'Идентификатор бота tlg';

CREATE TABLE prj_link(
    linkid CHARACTER VARYING(32) NOT NULL UNIQUE,
    link CHARACTER VARYING(300) NOT NULL,
    
    CONSTRAINT pk_prj_link PRIMARY KEY (linkid)
)
WITH(
	OIDS=FALSE
);

ALTER TABLE prj_link OWNER TO postgres;
COMMENT ON TABLE prj_link IS 'Таблица  ссылок';
COMMENT ON TABLE prj_link.linkid IS 'Идентификатор ссылки';
COMMENT ON TABLE prj_link.link IS 'Ссылка на отслеживаемую страницу';

CREATE TABLE ref_link_user(
    refid CHARACTER VARYING(32) NOT NULL UNIQUE,
    userid CHARACTER VARYING(32) NOT NULL,
    linkid CHARACTER VARYING(32) NOT NULL,
    
    CONSTRAINT pk_ref_link_user PRIMARY KEY (refid) REFERENCES prj_user (userid) ON DELETE CASCADE
)
WITH(
	OIDS=FALSE
);

ALTER TABLE ref_link_user OWNER TO postgres;
COMMENT ON TABLE ref_link_user IS 'Таблица связи ссылка - пользоваетль';
COMMENT ON TABLE ref_link_user.refid IS 'Идентификатор связи ссылка - рользователь';
COMMENT ON TABLE ref_link_user.userid IS 'Идентификатор пользователя';
COMMENT ON TABLE ref_link_user.linkid IS 'Идентификатор ссылки на страницу';

