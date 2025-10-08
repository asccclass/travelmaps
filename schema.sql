CREATE TABLE `travel` (
  `travelID` int(11) NOT NULL,
  `title` varchar(200) NOT NULL,
  `startDate` date NOT NULL,
  `endDate` date NOT NULL,
  `summary` text NOT NULL,
  `usrNo` int(11) NOT NULL DEFAULT '0'
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `trips` (
  `id` int(11) NOT NULL,
  `name` varchar(200) NOT NULL,
  `lat` varchar(20) NOT NULL,
  `lng` varchar(20) NOT NULL,
  `description` text NOT NULL,
  `dayz` int(11) NOT NULL DEFAULT '1',
  `logTime` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `photos` (
  `photoID` int(11) NOT NULL,
  `id` int(11) NOT NULL COMMENT '參考 trips.id',
  `thumbnail` text NOT NULL,
  `full` text NOT NULL,
  `caption` text NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8;