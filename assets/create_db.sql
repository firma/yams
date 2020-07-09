CREATE DATABASE yams CHARACTER SET utf8 COLLATE utf8_general_ci;

CREATE TABLE `account`(
    `id` INT UNSIGNED AUTO_INCREMENT,
    `username` VARCHAR(100) NOT NULL,
    `password` VARCHAR(100) NOT NULL,
    `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8;

CREATE TABLE IF NOT EXISTS "map_info" (
    id integer not null primary key autoincrement,
    file_name varchar(100) default NULL,
    title varchar(100) default NULL,
    mini_map integer default NULL,
    big_map integer default NULL,
    music integer default NULL,
    light integer default NULL,
    map_dark_light integer default NULL,
    mine_index integer default NULL,
    no_teleport integer default NULL,
    no_reconnect integer default NULL,
    no_random integer default NULL,
    no_escape integer default NULL,
    no_recall integer default NULL,
    no_drug integer default NULL,
    no_position integer default NULL,
    no_fight integer default NULL,
    no_throw_item integer default NULL,
    no_drop_player integer default NULL,
    no_drop_monster integer default NULL,
    no_names integer default NULL,
    no_mount integer default NULL,
    need_bridle integer default NULL,
    fight integer default NULL,
    fire integer default NULL,
    lightning integer default NULL,
    no_town_teleport integer default NULL,
    no_reincarnation integer default NULL,
    no_reconnect_map varchar(100) default NULL,
    fire_damage integer default NULL,
    lightning_damage integer default NULL
);