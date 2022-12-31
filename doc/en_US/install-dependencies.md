# Install Dependencies

* [CentOS 7](#centos-7)
* [Ubuntu 18](#ubuntu-18)

CentOS 7
========

> Install MariaDB

```shell
# Addition `MariaDB` Repo.

sudo cat > /etc/yum.repos.d/MariaDB.repo <<EOF
[mariadb]
name = MariaDB
baseurl = http://yum.mariadb.org/10.2/centos7-amd64
gpgkey=https://yum.mariadb.org/RPM-GPG-KEY-MariaDB
gpgcheck=1
EOF


# Install `MariaDB` Server and Client.

sudo yum -y install MariaDB-server MariaDB-client


# Start `MariaDB` Server.

sudo systemctl start mariadb


# Initialize `MariaDB` and set root password.

sudo mysql_secure_installation
```


Ubuntu 18
==========

> Install MariaDB

```shell
# Key is imported and the repository added.

sudo apt-get -y install software-properties-common
sudo apt-key adv --fetch-keys 'https://mariadb.org/mariadb_release_signing_key.asc'
sudo add-apt-repository 'deb [arch=amd64,arm64,ppc64el] http://mirror.hosting90.cz/mariadb/repo/10.2/ubuntu bionic main'
sudo apt update


# Install `MariaDB` and set root password (After installation, set the root password according to the system prompt).

sudo apt-get -y install mariadb-server
```
