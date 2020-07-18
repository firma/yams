CREATE DATABASE mirdata CHARACTER SET utf8 COLLATE utf8_general_ci;

CREATE TABLE IF NOT EXISTS `map_info` (
    `id` INT UNSIGNED AUTO_INCREMENT,
    `file_name` varchar(100) default NULL,
    `title` varchar(100) default NULL,
    `mini_map` integer default NULL,
    `big_map` integer default NULL,
    `music` integer default NULL,
    `light` integer default NULL,
    `map_dark_light` integer default NULL,
    `mine_index` integer default NULL,
    `no_teleport` integer default NULL,
    `no_reconnect` integer default NULL,
    `no_random` integer default NULL,
    `no_escape` integer default NULL,
    `no_recall` integer default NULL,
    `no_drug` integer default NULL,
    `no_position` integer default NULL,
    `no_fight` integer default NULL,
    `no_throw_item` integer default NULL,
    `no_drop_player` integer default NULL,
    `no_drop_monster` integer default NULL,
    `no_names` integer default NULL,
    `no_mount` integer default NULL,
    `need_bridle` integer default NULL,
    `fight` integer default NULL,
    `fire` integer default NULL,
    `lightning` integer default NULL,
    `no_town_teleport` integer default NULL,
    `no_reincarnation` integer default NULL,
    `no_reconnect_map` varchar(100) default NULL,
    `fire_damage` integer default NULL,
    `lightning_damage` integer default NULL,
    PRIMARY KEY (`id`)
);

INSERT INTO map_info VALUES(1,'0','比奇省',101,135,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,'',0,0);

CREATE TABLE `npc_info` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `map_id` int(11) DEFAULT NULL,
  `file_name` varchar(200) DEFAULT NULL,
  `name` varchar(200) DEFAULT NULL,
  `image` int(11) DEFAULT NULL,
  `location_x` int(11) DEFAULT NULL,
  `location_y` int(11) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE = InnoDB AUTO_INCREMENT = 294 DEFAULT CHARSET = utf8;

INSERT INTO npc_info VALUES(1,1,'比奇省/边境村/边境传送员.txt','边境传送员',15,287,615);

CREATE TABLE `item_info` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(200) DEFAULT NULL,
  `type` int(11) DEFAULT NULL,
  `grade` int(11) DEFAULT NULL,
  `required_type` int(11) DEFAULT NULL,
  `required_class` int(11) DEFAULT NULL,
  `required_gender` int(11) DEFAULT NULL,
  `item_set` int(11) DEFAULT NULL,
  `shape` int(11) DEFAULT NULL,
  `weight` int(11) DEFAULT NULL,
  `light` int(11) DEFAULT NULL,
  `required_amount` int(11) DEFAULT NULL,
  `image` int(11) DEFAULT NULL,
  `durability` int(11) DEFAULT NULL,
  `stack_size` int(11) DEFAULT NULL,
  `price` int(11) DEFAULT NULL,
  `min_ac` int(11) DEFAULT NULL,
  `max_ac` int(11) DEFAULT NULL,
  `min_mac` int(11) DEFAULT NULL,
  `max_mac` int(11) DEFAULT NULL,
  `min_dc` int(11) DEFAULT NULL,
  `max_dc` int(11) DEFAULT NULL,
  `min_mc` int(11) DEFAULT NULL,
  `max_mc` int(11) DEFAULT NULL,
  `min_sc` int(11) DEFAULT NULL,
  `max_sc` int(11) DEFAULT NULL,
  `hp` int(11) DEFAULT NULL,
  `mp` int(11) DEFAULT NULL,
  `accuracy` int(11) DEFAULT NULL,
  `agility` int(11) DEFAULT NULL,
  `luck` int(11) DEFAULT NULL,
  `attack_speed` int(11) DEFAULT NULL,
  `start_item` int(11) DEFAULT NULL,
  `bag_weight` int(11) DEFAULT NULL,
  `hand_weight` int(11) DEFAULT NULL,
  `wear_weight` int(11) DEFAULT NULL,
  `effect` int(11) DEFAULT NULL,
  `strong` int(11) DEFAULT NULL,
  `magic_resist` int(11) DEFAULT NULL,
  `poison_resist` int(11) DEFAULT NULL,
  `health_recovery` int(11) DEFAULT NULL,
  `spell_recovery` int(11) DEFAULT NULL,
  `poison_recovery` int(11) DEFAULT NULL,
  `hp_rate` int(11) DEFAULT NULL,
  `mp_rate` int(11) DEFAULT NULL,
  `critical_rate` int(11) DEFAULT NULL,
  `critical_damage` int(11) DEFAULT NULL,
  `bools` int(11) DEFAULT NULL,
  `max_ac_rate` int(11) DEFAULT NULL,
  `max_mac_rate` int(11) DEFAULT NULL,
  `holy` int(11) DEFAULT NULL,
  `freezing` int(11) DEFAULT NULL,
  `poison_attack` int(11) DEFAULT NULL,
  `bind` int(11) DEFAULT NULL,
  `reflect` int(11) DEFAULT NULL,
  `hp_drain_rate` int(11) DEFAULT NULL,
  `unique_item` int(11) DEFAULT NULL,
  `random_stats_id` int(11) DEFAULT NULL,
  `can_fast_run` int(11) DEFAULT NULL,
  `can_awakening` int(11) DEFAULT NULL,
  `tool_tip` varchar(2000) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1347 DEFAULT CHARSET=utf8;

INSERT INTO item_info VALUES(1,'屠龙',1,2,0,7,3,0,29,92,0,40,57,33000,1,75000,0,0,0,0,5,40,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,2,0,0,0,0,0,0,0,0,0,1,0,1,'');
INSERT INTO item_info VALUES(2,'重盔甲(男)',2,1,0,7,1,0,3,23,0,22,62,25000,1,10000,4,7,2,3,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,2,0,1,'');
