CREATE TABLE `schools` (
                           `id` int NOT NULL AUTO_INCREMENT,
                           `name` varchar(255) NOT NULL,
                           `brief_introduction` text,
                           `school_code` varchar(100) DEFAULT NULL,
                           `master_point` int DEFAULT NULL,
                           `phd_point` int DEFAULT NULL,
                           `research_project` int DEFAULT NULL,
                           `title_double_first_class` tinyint(1) DEFAULT NULL,
                           `title_985` tinyint(1) DEFAULT NULL,
                           `title_211` tinyint(1) DEFAULT NULL,
                           `title_college` tinyint(1) DEFAULT NULL,
                           `title_undergraduate` tinyint(1) DEFAULT NULL,
                           `region` varchar(255) DEFAULT NULL,
                           `website` varchar(255) DEFAULT NULL,
                           `recruitment_phone` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
                           `email` varchar(100) DEFAULT NULL,
                           `promotion_rate` varchar(50) DEFAULT NULL,
                           `abroad_rate` varchar(50) DEFAULT NULL,
                           `employment_rate` varchar(50) DEFAULT NULL,
                           `double_first_class_disciplines` text,
                           PRIMARY KEY (`id`),
                           UNIQUE KEY `name` (`name`),
                           UNIQUE KEY `school_code` (`school_code`)
);

CREATE TABLE `scores` (
                          `id` int NOT NULL AUTO_INCREMENT,
                          `school_id` int DEFAULT NULL,
                          `location` int DEFAULT NULL,
                          `year` int DEFAULT NULL,
                          `type_id` int DEFAULT NULL,
                          `tag` varchar(50) DEFAULT NULL,
                          `lowest` int DEFAULT NULL,
                          `lowest_rank` int DEFAULT NULL,
                          `sg_name` varchar(50) DEFAULT NULL,
                          `batch_name` varchar(50) DEFAULT NULL,
                          PRIMARY KEY (`id`),
                          UNIQUE KEY `scores_school_id_IDX` (`school_id`,`location`,`type_id`,`year`,`tag`,`sg_name`,`batch_name`) USING BTREE,
                          CONSTRAINT `scores_ibfk_1` FOREIGN KEY (`school_id`) REFERENCES `schools` (`id`)
);


CREATE TABLE major_scores (
                            id INT NOT NULL AUTO_INCREMENT,
                            special_id INT,
                            location VARCHAR(255),
                            year INT,
                            kelei VARCHAR(255),
                            batch VARCHAR(255),
                            recruitment_number INT,
                            lowest_score INT,
                            lowest_rank INT,
                            PRIMARY KEY (`id`)
);


CREATE TABLE majors (
                       id INT NOT NULL AUTO_INCREMENT,
                       name VARCHAR(255),
                       national_feature bool,  -- 这里需要确认布尔类型的处理方式
                       level VARCHAR(255),
                       discipline_category VARCHAR(255),
                       major_category VARCHAR(255),
                       limit_year VARCHAR(255),
                       school_id INT,
                       special_id VARCHAR(255),
                       PRIMARY KEY (`id`)
);

create table school_nums (
                        id INT NOT NULL AUTO_INCREMENT,
                        school_id int,
                        year int,
                        type_id varchar(50),
                        number int,
                        PRIMARY KEY (`id`),
                        CONSTRAINT `school_id_fk` FOREIGN KEY (`school_id`) REFERENCES `schools` (`id`)
);

create table users(
    id       INT          NOT NULL AUTO_INCREMENT,
    username varchar(50)  not null unique,
    email    varchar(100) not null unique,
    password varchar(100) not null,
    province VARCHAR(255),
    exam_type VARCHAR(255),
    school_type VARCHAR(255),
    physics BOOLEAN,
    history BOOLEAN,
    chemistry BOOLEAN,
    biology BOOLEAN,
    geography BOOLEAN,
    politics BOOLEAN,
    score INT,
    province_rank INT,
    holland VARCHAR(255),
    interests TEXT,
    PRIMARY KEY (`id`)
);