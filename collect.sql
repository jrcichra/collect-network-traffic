create database collect_network_traffic;

create table collect_network_traffic.throughput(
	`id` bigint primary key auto_increment,
	`date` datetime not null default current_timestamp(),
	`interface` varchar(20) not null,
	`bytes` bigint not null,
	`src_name` varchar(1024) not null,
	`dst_name` varchar(1025) not null,
	`hostname` varchar(200) not null,
	`proto` varchar(10) not null,
	`src_port` int not null,
	`dst_port` int not null,
	`interval` int not null
);

CREATE INDEX throughput_date_IDX USING BTREE ON collect_network_traffic.throughput (`date`);
CREATE INDEX throughput_hostname_IDX USING BTREE ON collect_network_traffic.throughput (hostname);
CREATE INDEX throughput_interface_IDX USING BTREE ON collect_network_traffic.throughput (interface);


create user 'collect'@'%' identified by 'test';

grant all on collect_network_traffic.throughput to 'collect'@'%';

CREATE EVENT collect_network_traffic.auto_delete
ON SCHEDULE AT CURRENT_TIMESTAMP + INTERVAL 1 DAY 
ON COMPLETION PRESERVE
DO 
DELETE LOW_PRIORITY FROM collect_network_traffic.throughput WHERE datetime < DATE_SUB(NOW(), INTERVAL 7 DAY);

create user 'collect'@'%' identified by 'test';

grant all on collect_network_traffic.throughput to 'collect'@'%';