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

create user 'collect'@'%' identified by 'test';

grant all on collect_network_traffic.throughput to 'collect'@'%';