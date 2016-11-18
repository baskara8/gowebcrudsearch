/*
Navicat MySQL Data Transfer

Source Server         : MYSQL-MYLOCAL
Source Server Version : 50505
Source Host           : 127.0.0.1:3306
Source Database       : latihan

Target Server Type    : MYSQL
Target Server Version : 50505
File Encoding         : 65001

Date: 2016-11-18 14:26:38
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for siswa
-- ----------------------------
DROP TABLE IF EXISTS `siswa`;
CREATE TABLE `siswa` (
  `id` varchar(25) NOT NULL,
  `nama` varchar(26) NOT NULL,
  `email` varchar(32) NOT NULL,
  `password` varchar(50) DEFAULT NULL,
  `jekel` char(1) DEFAULT NULL,
  `foto` varchar(150) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=latin1;

-- ----------------------------
-- Records of siswa
-- ----------------------------
INSERT INTO `siswa` VALUES ('bima', 'Bima', 'bima@gmail.com', '0a5a5f6907194bba404d67fed4dec84f6f82e474', 'L', '1464550022_France-Flag.png');
INSERT INTO `siswa` VALUES ('dodiamanda', 'Dodi Amanda', 'dodi@gmail.com', '0a5a5f6907194bba404d67fed4dec84f6f82e474', 'L', null);
INSERT INTO `siswa` VALUES ('jani', 'Jani', 'jani@gmail.com', '646c3fad45809c2a958cce7b5df636ef8de93e7d', 'P', '1464550034_United-States-Flag.png');
INSERT INTO `siswa` VALUES ('rina', 'Rina Nose', 'rina@gmail.com', '70e21878d268fa8f82817f9278f8bae0fb108950', 'P', '1464550055_Indonesia-Flag.png');
INSERT INTO `siswa` VALUES ('riva', 'Riva Ananta Baskara', 'brockbask@gmail.com', '65547a478d66f0cf7b33c6bcf3654d214b3db295', 'L', null);
