CREATE TABLE `products` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `sku` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
  `price` int(11) NOT NULL,
  `qty` int(11) NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

insert  into `products`(`id`,`sku`,`name`,`price`,`qty`) values 
(1,'120P90','Google Home',49,10),
(2,'43N23P','MacBook Pro',5399,5),
(3,'A304SD','Alexa Speaker',109,10),
(4,'234234','Raspberry Pi B',30,2);